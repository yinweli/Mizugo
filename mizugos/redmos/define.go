package redmos

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// redmos 提供「雙層式資料庫」實作:
//   - 主要資料庫(Major): 以 Redis 為基礎(快取／高速查詢)
//   - 次要資料庫(Minor): 以 Mongo 為基礎(持久化／複雜查詢)
//   - 混合資料庫(Mixed): 以 Major+Minor 組合, 支援以「行為(Behavior)」序列化的雙段式提交
//
// # 組成與線程安全
//
// Major/Minor/Mixed 非執行緒安全, 若要跨 goroutine 共用, 請透過管理器 Redmomgr 做統一建立與存取(其內部有鎖保護), 建立順序為:
//   1) AddMajor / AddMinor 建立資料庫
//   2) AddMixed 綁定成 Mixed
//   3) 以 GetMajor / GetMinor / GetMixed 取得連線物件操作
//   4) 結束時呼叫 Finalize 釋放所有已註冊資料庫
//
// # URI 格式
//
// - RedisURI(前綴固定 redisdb://):
//   redisdb://[user:pass@]host1:port1[,hostN:portN]/[?options]
//   支援常見 go-redis 連線選項(如 dbid, poolSize, readTimeout, masterName 提供哨兵等)
// - MongoURI(標準 mongodb:// 語法):
//   mongodb://[user:pass@]host1[:port1][,hostN[:portN]]/[?options]
//   支援連線/逾時/連線池/副本集等常見選項
//
// # 快速上手
//
// 1) 以 Redmomgr 建立 Major/Minor/Mixed
//
//    mgr := redmos.NewRedmomgr()
//
//    // 建立 Redis 連線
//    majorURI := redmos.RedisURI("redisdb://127.0.0.1:6379/?dbid=0&clientName=redmos")
//
//    if _, err := mgr.AddMajor("cache", majorURI); err != nil {
//        panic(err)
//    } // if
//
//    // 建立 Mongo 連線
//    minorURI := redmos.MongoURI("mongodb://127.0.0.1:27017")
//
//    if _, err := mgr.AddMinor("store", minorURI, "demo"); err != nil {
//        panic(err)
//    } // if
//
//    // 綁定成 Mixed
//    mixed, err := mgr.AddMixed("demo", "cache", "store")
//
//    if err != nil {
//        panic(err)
//    } // if
//
// 2) 直接操作 Major(Redis)
//
//    major := mgr.GetMajor("cache")
//    submit := major.Submit()
//    submit.Set(ctx, "k", "v", 0)
//
//    if _, err := submit.Exec(ctx); err != nil {
//        panic(err)
//    } // if
//
// 3) 直接操作 Minor(Mongo)
//
//    minor := mgr.GetMinor("store")
//    submit := minor.Submit()
//
//    // 收集批次操作(以表名為 key)
//    submit.Operate("users", mongo.NewInsertOneModel().SetDocument(bson.M{"_KEY_":"u:1", "name":"alice"}))
//
//    if err := submit.Exec(ctx); err != nil {
//        panic(err)
//    } // if
//
// 4) 操作 Mixed (提交單一行為)
//
//    ctx := context.Background()
//    submit := mixed.Submit(ctx)
//
//    // 新增一個行為
//    submit.Add(&Set[Profile]{
//        MajorEnable: true,          // 寫入 Redis
//        MinorEnable: true,          // 寫入 Mongo
//        Meta:        UserMeta{},    // 提供 Major/Minor 索引資訊(需實作 Metaer)
//        Key:         "1001",        // 業務主鍵
//        Data:        &Profile{Name: "Alice", Lv: 12},
//    })
//
//    if err := submit.Exec(); err != nil {
//        panic(err)
//    } // if
//
// 5) 操作 Mixed (提交多個行為)
//
//    ctx := context.Background()
//    submit := mixed.Submit(ctx)
//
//    // 可混合不同型別的行為, 都會在同一次 Exec 中完成
//    submit.Add(&Set[Profile]{MajorEnable: true, MinorEnable: true, Meta: UserMeta{}, Key: "1002", Data: &Profile{Name: "Bob",   Lv: 8}})
//    submit.Add(&Set[Profile]{MajorEnable: true, MinorEnable: false,Meta: UserMeta{}, Key: "1003", Data: &Profile{Name: "Carol", Lv: 5}})  // 僅寫 Redis
//    submit.Add(&Set[Profile]{MajorEnable: false,MinorEnable: true, Meta: UserMeta{}, Key: "1004", Data: &Profile{Name: "Duke",  Lv: 20}}) // 僅寫 Mongo
//
//    if err := submit.Exec(); err != nil {
//        panic(err)
//    } // if
//
// # 內建行為
//
//   - Lock / Unlock: 鎖定 / 解鎖
//   - LockIf / UnlockIf: 視旗標決定是否鎖定 / 解鎖
//
// # 建立行為
//
// 這裡示範從零建立一個「寫入 Major(Redis) 與 Minor(Mongo)」的行為(非使用現有 Set[T]), 讓你了解行為的標準寫法與提交流程
//
// 行為生命週期, 由 Mixed.Submit.Exec() 驅動:
//   Initialize(ctx, major, minor) → Prepare() → (Major pipeline Exec) → Complete() → (Minor bulk Exec)
//
// 設計要點:
//   - 內嵌 Behave: 省去自己保存 ctx/major/minor 的程式碼(Submit.Add 會自動呼叫 Initialize 幫你注入)
//   - Prepare:
//       * 做參數檢查(Meta/Key/Data...)
//       * 若要對 Redis 下指令, 直接在這裡把指令「排進 Major 的 pipeline」並把返回的 Cmd 暫存起來
//       * 若要「跳過一次存檔」的業務條件, 這裡可判斷(例如 Data 實作 Saver 且 Save=false)
//   - Complete:
//       * 依前一步 Major 的結果(Cmd.Result())決定要不要對 Mongo 收集 BulkWrite
//       * 在這裡只做「收集」; 真正送出由 Submit.Exec() 在最後呼叫 minor.Exec(ctx)
//
// 使用前置:
//   - T 必須是 struct(非指標); 然而行為參數 Data 要是 *T(以便序列化/inline)
//   - 需要一個實作 Metaer 的 meta(提供 MajorKey/MinorKey/MinorTable)
//   - 若要寫入 Redis/Mongo, 就分別開 MajorEnable/MinorEnable
//   - 欄位需有適當的 bson tag; Minor 會以 MinorData[T]{K, D} 包裝, K 對應 MongoKey="_KEY_"
//
// 程式碼:
//   // 這個類別是示範, 名稱用 DemoSet 避免與你現有 Set[T] 衝突; 你可以複製後改名為自己的 Set 行為
//   type DemoSet[T any] struct {
//       redmos.Behave                       // 內嵌, 共用工具: Ctx / Major / Minor
//       MajorEnable bool                    // 是否寫 Major(Redis)
//       MinorEnable bool                    // 是否寫 Minor(Mongo)
//       Meta        redmos.Metaer           // 提供 Major/Minor 的鍵與表資訊
//       Key         string                  // 業務主鍵(用於組合 Major/Minor 的實際鍵)
//       Data        *T                      // 要寫入的資料(T 必須是 struct)
//       keepTTL     bool                    // 是否保留既有 TTL(示範可調參)
//       cmd         *redis.StatusCmd        // 暫存 Redis 指令回傳(供 Complete 判斷)
//   }
//
//   // Initialize 由 Submit.Add 自動呼叫(通常不需要覆寫, 這裡示範顯式轉呼叫 Behave.Initialize)
//   func (this *DemoSet[T]) Initialize(ctx context.Context, major redmos.MajorSubmit, minor *redmos.MinorSubmit) {
//       this.Behave.Initialize(ctx, major, minor)
//       // 額外初始(示範可選): 預設保留 TTL
//       this.keepTTL = true
//   }
//
//   // Prepare: 參數檢查 + 排 Major pipeline 指令
//   func (this *DemoSet[T]) Prepare() error {
//       // 1) 基本參數檢查
//       if this.Meta == nil {
//           return fmt.Errorf("prepare: meta nil")
//       } // if
//
//       if this.Key == "" {
//           return fmt.Errorf("prepare: key empty")
//       } // if
//
//       if this.Data == nil {
//           return fmt.Errorf("prepare: data nil")
//       } // if
//
//       // 2) 可選: 若資料實作 Saver 且 Save=false, 直接跳過(Major/Minor 都不做)
//       if saver, ok := any(this.Data).(redmos.Saver); ok && !saver.GetSave() {
//           return nil
//       } // if
//
//       // 3) 排 Major 指令(在 Prepare 階段排進 pipeline; Submit.Exec 會幫你執行)
//       if this.MajorEnable {
//           // Major 實際鍵
//           rkey := this.Meta.MajorKey(this.Key)
//
//           // 序列化資料(*T → JSON)
//           raw, err := json.Marshal(this.Data)
//
//           if err != nil {
//               return fmt.Errorf("prepare marshal: %w", err)
//           } // if
//
//           // 決定 TTL 策略(示範: 保留既有 TTL; 你也可以改為指定逾時)
//           ttl := redis.KeepTTL
//
//           if this.keepTTL == false {
//               // 例如: 指定 24h(示範)
//               ttl = 24 * time.Hour
//           } // if
//
//           // 把 SET 指令排進 pipeline(回傳 *StatusCmd 暫存, 供 Complete 判斷)
//           this.cmd = this.Major().Set(this.Ctx(), rkey, raw, ttl)
//       } // if
//
//       // 4) 檢查 Minor 前置(若要寫 Minor, 需要有效表名)
//       if this.MinorEnable && this.Meta.MinorTable() == "" {
//           return fmt.Errorf("prepare: table empty")
//       } // if
//
//       return nil
//   }
//
//   // Complete: 依 Major 的結果決定要不要收集 Minor 的 upsert
//   func (this *DemoSet[T]) Complete() error {
//       // 1) 可選: 資料旗標判斷(與 Prepare 對應, 如果 Prepare 已 return, 理論上不會到這裡; 保險起見再判一次)
//       if saver, ok := any(this.Data).(redmos.Saver); ok && saver.GetSave() == false {
//           return nil
//       } // if
//
//       // 2) 先確認 Major 執行結果(若前面有排 Major 指令)
//       if this.MajorEnable && this.cmd != nil {
//           s, err := this.cmd.Result()
//
//           if err != nil {
//               return fmt.Errorf("complete major: %w", err)
//           } // if
//
//           if s != redmos.RedisOk {
//               return fmt.Errorf("complete major: unexpected redis reply: %v", s)
//           } // if
//       } // if
//
//       // 3) 收集 Minor upsert(真正送出由 Submit.Exec 在最後呼叫 minor.Exec(ctx))
//       if this.MinorEnable {
//           key := this.Meta.MinorKey(this.Key)
//           table := this.Meta.MinorTable()
//
//           // 用 MinorData[T] 打包, K=_KEY_, D=資料(inline)
//           model := mongo.NewReplaceOneModel().
//               SetUpsert(true).
//               SetFilter(bson.M{redmos.MongoKey: key}).
//               SetReplacement(&redmos.MinorData[T]{K: key, D: this.Data})
//
//           this.Minor().Operate(table, model)
//       } // if
//
//       return nil
//   }

const (
	// Timeout 為與資料庫互動的預設逾時時間
	Timeout = 30 * time.Second

	// RedisNil 表示「查無資料」時的空字串占位符
	RedisNil = ""

	// RedisOk 表示 Redis 指令成功時的標準回覆字串
	RedisOk = "OK"

	// RedisTTL 表示 Redis 的 SET 命令是否維持原逾期時間
	RedisTTL = redis.KeepTTL

	// MongoKey 表示 Mongo 資料庫的索引欄位名稱
	MongoKey = "_KEY_"
)

// Metaer 元資料介面, 提供主要/次要資料庫操作所需的資訊
type Metaer interface {
	// MajorKey 取得主要資料庫的索引值
	MajorKey(key any) string

	// MinorKey 取得次要資料庫的索引值
	MinorKey(key any) string

	// MinorTable 取得次要資料庫的表名稱
	MinorTable() string
}

// Saver 儲存判斷介面, 用於描述當前物件是否需要儲存
type Saver interface {
	// GetSave 取得儲存旗標
	GetSave() bool
}
