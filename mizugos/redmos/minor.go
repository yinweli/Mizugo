package redmos

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// newMinor 建立次要資料庫, 並依據傳入的 MongoURI 立刻進行連線, dbName 不可為空
func newMinor(uri MongoURI, dbName string) (minor *Minor, err error) {
	if dbName == "" {
		return nil, fmt.Errorf("newMinor: dbName empty")
	} // if

	client, err := uri.Connect(context.Background())

	if err != nil {
		return nil, fmt.Errorf("newMinor: %w", err)
	} // if

	return &Minor{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Minor 次要資料庫(Mongo)
//
// 以 Mongo 為基礎的資料庫封裝, 提供三種使用模式:
//   - Submit: 取得 MinorSubmit, 適合批次/管線化命令
//   - Client: 取得原生 mongo.Client, 適合需要原生 API 的場景
//   - Database: 取得原生 mongo.Database, 適合需要原生 API 的場景
//
// Minor 不是執行緒安全, 若跨 goroutine 共用, 請由上層管理器(Redmomgr)負責同步保護
type Minor struct {
	client   *mongo.Client   // 客戶端物件
	database *mongo.Database // 資料庫物件
}

// Submit 取得 Pipeline 執行器
func (this *Minor) Submit() *MinorSubmit {
	if this.client != nil && this.database != nil {
		return &MinorSubmit{
			database: this.database,
			operate:  map[string][]mongo.WriteModel{},
		}
	} // if

	return nil
}

// Client 取得 Mongo 客戶端
func (this *Minor) Client() *mongo.Client {
	return this.client
}

// Database 取得 Mongo 資料庫
func (this *Minor) Database() *mongo.Database {
	return this.database
}

// SwitchDB 切換 Mongo DB, dbName 不可為空
func (this *Minor) SwitchDB(dbName string) error {
	if this.client == nil {
		return fmt.Errorf("minor switch: client nil")
	} // if

	if dbName == "" {
		return fmt.Errorf("minor switch: dbName empty")
	} // if

	this.database = this.client.Database(dbName)
	return nil
}

// DropDB 清除資料庫
func (this *Minor) DropDB() {
	if this.database != nil {
		_ = this.database.Drop(context.Background())
	} // if
}

// stop 關閉並釋放資料庫
func (this *Minor) stop() {
	if this.client != nil {
		_ = this.client.Disconnect(context.Background())
		this.client = nil
		this.database = nil
	} // if
}

// MinorSubmit Pipeline 執行器
//
// 用法:
//   - 先以 Operate 收集以表格名稱為索引的批次操作
//   - 最後執行 Exec 將批次操作送出
//   - 清除批次操作列表
type MinorSubmit struct {
	database *mongo.Database               // 資料庫物件
	operate  map[string][]mongo.WriteModel // 批次操作列表
}

// Collection 取得表格物件
func (this *MinorSubmit) Collection(table string) *mongo.Collection {
	return this.database.Collection(table)
}

// Operate 新增批次操作
func (this *MinorSubmit) Operate(table string, operate mongo.WriteModel) *MinorSubmit {
	this.operate[table] = append(this.operate[table], operate)
	return this
}

// Exec 執行批次操作, 會清空已收集的批次操作
func (this *MinorSubmit) Exec(ctx context.Context) error {
	for table, itor := range this.operate {
		if _, err := this.database.Collection(table).BulkWrite(ctx, itor); err != nil {
			return fmt.Errorf("minorSubmit exec: %w", err)
		} // if
	} // for

	this.operate = map[string][]mongo.WriteModel{}
	return nil
}

// MinorData 泛型資料殼，用於在次要資料庫存取時維持固定索引欄位
type MinorData[T any] struct {
	K string `bson:"_KEY_"`   // 索引欄位, bson名稱必須與 MongoKey 一致
	D *T     `bson:",inline"` // 資料欄位, 此欄位利用inline達成內嵌效果(會將 T 的欄位展平成與 K 同層)
}

// MinorIndex 建立專門用於次要資料庫的索引
//   - 以 Metaer.MinorTable 為表格名稱
//   - 以 MongoKey 為索引欄位, 排序為遞增, 設定 Unique=true
func MinorIndex(meta Metaer) *Index {
	return &Index{
		Name:   fmt.Sprintf("%v_minor_index", meta.MinorTable()),
		Table:  meta.MinorTable(),
		Field:  MongoKey,
		Order:  1,
		Unique: true,
	}
}
