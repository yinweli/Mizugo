package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Incr 遞增行為, 以索引值到主要/次要資料庫中遞增數值, 使用上有以下幾點須注意
//   - 執行前設定好 MinorEnable; 由於遞增行為只會在主要資料庫中執行, 因此次要資料庫僅用於備份
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Incr 指示遞增數值
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Data 來取得資料
//   - 由於遞增行為是以int64來運作, 因此使用時可能需要轉換
type Incr struct {
	Behave                    // 行為物件
	MinorEnable bool          // 啟用次要資料庫
	Meta        Metaer        // 元資料
	Key         string        // 索引值
	Incr        int64         // 遞增數值
	Data        *IncrData     // 資料物件
	cmd         *redis.IntCmd // 命令結果
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
	data, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("incr complete: %w: %v", err, this.Key)
	} // if

	if this.Data == nil {
		this.Data = new(IncrData)
	} // if

	this.Data.Data = data

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{MongoKey: key}).
			SetReplacement(&MinorData[IncrData]{
				K: key,
				D: this.Data,
			}))
	} // if

	return nil
}
