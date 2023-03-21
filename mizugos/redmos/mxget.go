package redmos

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get 取值行為, 以索引字串到主要資料庫中取得資料, 使用上有以下幾點須注意
//   - 此行為結構需與泛型共同運作, 填入的泛型類型T需要是結構型別, 請不要填入指標型別
//   - 執行前必須設定好 Key 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Result 來判斷取值是否成功
//   - 執行後可用 Data 來取得資料
//   - 在內部執行過程中, 索引字串會被轉為小寫
type Get[T any] struct {
	Behave
	Key    string           // 索引值
	Result bool             // 執行結果
	Data   *T               // 資料物件
	get    *redis.StringCmd // 取值命令結果
}

// Prepare 前置處理
func (this *Get[T]) Prepare() error {
	if this.Key == "" {
		return fmt.Errorf("get prepare: key empty")
	} // if

	key := FormatKey(this.Key)
	this.Result = false
	this.get = this.Major().Get(this.Ctx(), key)
	return nil
}

// Complete 完成處理
func (this *Get[T]) Complete() error {
	value, err := this.get.Result()

	if err != redis.Nil && err != nil {
		return fmt.Errorf("get complete: %w", err)
	} // if

	if value != RedisNil {
		if this.Data == nil {
			this.Data = new(T)
		} // if

		if err = json.Unmarshal([]byte(value), this.Data); err != nil {
			return fmt.Errorf("get complete: %w", err)
		} // if

		this.Result = true
	} // if

	return nil
}

// Set 設值行為, 以索引字串與資料到主要/次要資料庫中儲存資料, 使用上有以下幾點須注意
//   - 此行為結構需與泛型共同運作, 填入的泛型類型T需要是結構型別, 請不要填入指標型別
//   - 由於會儲存到次要資料庫中, 因此泛型類型T的成員都需要設定好`bson:name`屬性
//   - 執行前必須設定好 Table, Field, Key 並且不能為空字串
//   - 執行前必須設定好 Data, 並且不能為nil
//   - 在內部執行過程中, 索引欄位與索引字串會被轉為小寫
type Set[T any] struct {
	Behave
	Table string           // 表格名稱
	Field string           // 索引名稱
	Key   string           // 索引值
	Data  *T               // 資料物件
	set   *redis.StatusCmd // 設值命令結果
}

// Prepare 前置處理
func (this *Set[T]) Prepare() error {
	if this.Table == "" {
		return fmt.Errorf("set prepare: table empty")
	} // if

	if this.Field == "" {
		return fmt.Errorf("set prepare: field empty")
	} // if

	if this.Key == "" {
		return fmt.Errorf("set prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("set prepare: data nil")
	} // if

	value, err := json.Marshal(this.Data)

	if err != nil {
		return fmt.Errorf("set prepare: %w", err)
	} //

	key := FormatKey(this.Key)
	this.set = this.Major().Set(this.Ctx(), key, value, 0)
	return nil
}

// Complete 完成處理
func (this *Set[T]) Complete() error {
	value, err := this.set.Result()

	if err != nil {
		return fmt.Errorf("set complete: %w", err)
	} // if

	if value != RedisOk {
		return fmt.Errorf("set complete: save to redis failed")
	} // if

	field := FormatField(this.Field)
	key := FormatKey(this.Key)
	filter := bson.D{{Key: field, Value: key}}
	opt := options.Replace().SetUpsert(true)

	if _, err = this.Minor().Table(this.Table).ReplaceOne(this.Ctx(), filter, this.Data, opt); err != nil {
		return fmt.Errorf("set complete: %w", err)
	} // if

	return nil
}
