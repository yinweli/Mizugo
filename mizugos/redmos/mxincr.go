package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Incr 遞增行為, 以索引字串到主要/次要資料庫中遞增數值, 使用上有以下幾點須注意
//   - 此行為結構需與泛型共同運作, 填入的泛型類型T必須是int/int32/int64其中之一
//   - 執行前必須設定好 Table, Field, Key 並且不能為空字串
//   - 執行前必須設定好 Incr, 這是本次遞增值
//   - 執行後可用 Data 來取得遞增後的數值
//   - 在內部執行過程中, 索引欄位與索引字串會被轉為小寫
//   - 由於redis遞增都是以int64來運作, 因此如果泛型類型T是 int/int32, 需要考慮是否會溢出
type Incr[T int | int32 | int64] struct {
	Behave
	Table string        // 表格名稱
	Field string        // 索引名稱
	Key   string        // 索引值
	Incr  T             // 遞增數值
	Data  T             // 資料物件
	incr  *redis.IntCmd // 遞增命令結果
}

// Prepare 前置處理
func (this *Incr[T]) Prepare() error {
	if this.Table == "" {
		return fmt.Errorf("incr prepare: table empty")
	} // if

	if this.Field == "" {
		return fmt.Errorf("incr prepare: field empty")
	} // if

	if this.Key == "" {
		return fmt.Errorf("incr prepare: key empty")
	} // if

	this.incr = this.Major().IncrBy(this.Ctx(), this.Key, int64(this.Incr))
	return nil
}

// Complete 完成處理
func (this *Incr[T]) Complete() error {
	value, err := this.incr.Result()

	if err != nil {
		return fmt.Errorf("incr complete: %w", err)
	} // if

	filter := bson.D{{Key: this.Field, Value: this.Key}}
	data := bson.M{this.Field: this.Key, "value": value} // 在mongo中, 遞增值固定儲存在value欄位
	opt := options.Replace().SetUpsert(true)

	if _, err = this.Minor().Table(this.Table).ReplaceOne(this.Ctx(), filter, data, opt); err != nil {
		return fmt.Errorf("incr complete: %w", err)
	} // if

	this.Data = T(value)
	return nil
}
