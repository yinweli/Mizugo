package redmos

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get 取值行為, 以索引字串到主要資料庫中取得資料, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 需要事先建立好與 Metaer 介面符合的元資料結構, 並填寫到 Meta
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Result 來判斷取值是否成功
//   - 執行後可用 Data 來取得資料
type Get[T any] struct {
	Behave                  // 行為物件
	Meta   Metaer           // 元資料
	Key    string           // 索引值
	Result bool             // 執行結果
	Data   *T               // 資料物件
	get    *redis.StringCmd // 取值命令結果
}

// Prepare 前置處理
func (this *Get[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("get prepare: meta nil")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("get prepare: table empty")
	} // if

	if this.Meta.MinorField() == "" {
		return fmt.Errorf("get prepare: field empty")
	} // if

	if this.Key == "" {
		return fmt.Errorf("get prepare: key empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	this.Result = false
	this.get = this.Major().Get(this.Ctx(), key)
	return nil
}

// Complete 完成處理
func (this *Get[T]) Complete() error {
	data, err := this.get.Result()

	if err != redis.Nil && err != nil {
		return fmt.Errorf("get complete: %w: %v", err, this.Key)
	} // if

	if data != RedisNil {
		if this.Data == nil {
			this.Data = new(T)
		} // if

		if err = json.Unmarshal([]byte(data), this.Data); err != nil {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		this.Result = true
	} // if

	return nil
}

// Set 設值行為, 以索引字串與資料到主要/次要資料庫中儲存資料, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 需要事先建立好與 Metaer 介面符合的元資料結構, 並填寫到 Meta
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data 並且不能為nil
type Set[T any] struct {
	Behave                  // 行為物件
	Meta   Metaer           // 元資料
	Key    string           // 索引值
	Data   *T               // 資料物件
	set    *redis.StatusCmd // 設值命令結果
}

// Prepare 前置處理
func (this *Set[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("set prepare: meta nil")
	} // if

	if this.Meta.MinorTable() == "" {
		return fmt.Errorf("set prepare: table empty")
	} // if

	if this.Meta.MinorField() == "" {
		return fmt.Errorf("set prepare: field empty")
	} // if

	if this.Key == "" {
		return fmt.Errorf("set prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("set prepare: data nil")
	} // if

	key := this.Meta.MajorKey(this.Key)
	data, err := json.Marshal(this.Data)

	if err != nil {
		return fmt.Errorf("set prepare: %w: %v", err, this.Key)
	} //

	this.set = this.Major().Set(this.Ctx(), key, data, 0)
	return nil
}

// Complete 完成處理
func (this *Set[T]) Complete() error {
	data, err := this.set.Result()

	if err != nil {
		return fmt.Errorf("set complete: %w: %v", err, this.Key)
	} // if

	if data != RedisOk {
		return fmt.Errorf("set complete: save to redis failed: %v", this.Key)
	} // if

	key := this.Meta.MinorKey(this.Key)
	table := this.Meta.MinorTable()
	field := this.Meta.MinorField()
	filter := bson.D{{Key: field, Value: key}}
	opt := options.Replace().SetUpsert(true)

	if _, err = this.Minor().Table(table).ReplaceOne(this.Ctx(), filter, this.Data, opt); err != nil {
		return fmt.Errorf("set complete: %w: %v", err, this.Key)
	} // if

	return nil
}
