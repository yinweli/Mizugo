package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Incr 遞增行為
//
// 以索引鍵(Key)在主要資料庫中對數值進行原子性遞增, 並可選擇將最新結果備份至次要資料庫;
// 若成功完成後, 會將結果記錄於 Data, 最後可選擇以 Done 回呼帶出結果
//
// 事前準備:
//   - (可選)設定 MinorEnable: true 表示完成時會將最新結果寫入次要資料庫
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MajorKey；若啟用備份，還需提供 MinorKey 與 MinorTable)
//   - 設定 Key: 不可為空字串
//   - 設定 Incr: 遞增(或遞減)的步進值, 允許負值
//   - (可選)設定 Done: 完成時的回呼函式, 參數為遞增後的整數結果
//
// 注意:
//   - 本行為的「遞增」僅使用主要資料庫; 次要資料庫僅作為結果備份
//   - 遞增操作為原子性; 對應 Redis 指令為 INCRBY
//   - 若鍵不存在, 會以 0 為起始值再遞增
//   - 由於以 int64 運作, 超出 int64 範圍會導致錯誤
type Incr struct {
	Behave                       // 行為物件
	MinorEnable bool             // 啟用次要資料庫
	Meta        Metaer           // 元資料
	Key         string           // 索引值
	Incr        int64            // 遞增數值
	Data        int64            // 資料物件
	Done        func(data int64) // 完成回呼
	cmd         *redis.IntCmd    // 命令結果
}

// IncrData 遞增資料
type IncrData struct {
	Data int64 `bson:"data"` // 遞增結果
}

// Prepare 前置處理
func (this *Incr) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("incr prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("incr prepare: key empty")
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("incr prepare: table empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	this.cmd = this.Major().IncrBy(this.Ctx(), key, this.Incr)
	return nil
}

// Complete 完成處理
func (this *Incr) Complete() error {
	incr, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("incr complete: %w: %v", err, this.Key)
	} // if

	this.Data = incr

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{MongoKey: key}).
			SetReplacement(&MinorData[IncrData]{
				K: key,
				D: &IncrData{Data: this.Data},
			}))
	} // if

	if this.Done != nil {
		this.Done(this.Data)
	} // if

	return nil
}
