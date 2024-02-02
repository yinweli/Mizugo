package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Incr 遞增行為, 以索引字串到主要/次要資料庫中遞增數值, 使用上有以下幾點須注意
//   - 資料類型必須是int/int32/int64其中之一, 並填寫到泛型類型T中
//   - 執行前設定好 MinorEnable. 請注意! 遞增行為必定會在主資料庫中執行, 因此無法禁用主資料庫
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Incr 這是本次遞增值
//   - 執行後可用 Data 來取得遞增後的數值
//   - 由於redis遞增是以int64來運作, 因此如果泛型類型T是int/int32, 需要考慮是否會溢出
type Incr[T int | int32 | int64] struct {
	Behave                    // 行為物件
	MinorEnable bool          // 啟用次要資料庫
	Meta        Metaer        // 元資料
	Key         string        // 索引值
	Incr        T             // 遞增數值
	Data        T             // 資料物件
	cmd         *redis.IntCmd // 命令結果
}

// Prepare 前置處理
func (this *Incr[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("incr prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("incr prepare: key empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	this.cmd = this.Major().IncrBy(this.Ctx(), key, int64(this.Incr))

	if this.MinorEnable {
		if this.Meta.MinorTable() == "" {
			return fmt.Errorf("incr prepare: table empty")
		} // if

		if this.Meta.MinorField() == "" {
			return fmt.Errorf("incr prepare: field empty")
		} // if
	} // if

	return nil
}

// Complete 完成處理
func (this *Incr[T]) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("incr complete: meta nil")
	} // if

	data, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("incr complete: %w: %v", err, this.Key)
	} // if

	if this.MinorEnable {
		table := this.Meta.MinorTable()
		field := this.Meta.MinorField()
		key := this.Meta.MinorKey(this.Key)
		filter := bson.D{{Key: field, Value: key}}
		replace := bson.M{field: key, "value": data} // 在mongo中, 遞增值固定儲存在value欄位
		this.Minor().Operate(table, mongo.NewReplaceOneModel().SetUpsert(true).SetFilter(filter).SetReplacement(replace))
	} // if

	this.Data = T(data)
	return nil
}
