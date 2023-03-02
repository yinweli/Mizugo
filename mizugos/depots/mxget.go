package depots

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Getter 取值行為, 以索引字串到主要資料庫中取得資料, 使用上有以下幾點須注意
//   - 此行為結構需與泛型共同運作, 填入的泛型類型 T 需要是結構型別, 請不要填入指標型別
//   - 使用前必須設定好 Key 並且不能為空字串
//   - 取得的資料會填入 Data 成員中, 如果 Data 成員為nil, 則內部程序會建立一個來填寫
//   - 當取值完成時, 可用 Result 成員來判斷取值是否成功
type Getter[T any] struct {
	Behave
	Key    string           // 索引字串
	Data   *T               // 資料物件
	Result bool             // 執行結果
	get    *redis.StringCmd // 取值命令結果
}

// Prepare 前置處理
func (this *Getter[T]) Prepare() error {
	if this.Key == "" {
		return fmt.Errorf("getter prepare: key empty")
	} // if

	key := FormatKey(this.Key)
	this.Result = false
	this.get = this.Major().Get(this.Ctx(), key)
	return nil
}

// Complete 完成處理
func (this *Getter[T]) Complete() error {
	value, err := this.get.Result()

	if err != redis.Nil && err != nil {
		return fmt.Errorf("getter complete: %w", err)
	} // if

	if value != RedisNil {
		if this.Data == nil {
			this.Data = new(T)
		} // if

		if err = json.Unmarshal([]byte(value), this.Data); err != nil {
			return fmt.Errorf("getter complete: %w", err)
		} // if

		this.Result = true
	} // if

	return nil
}

// Setter 設值行為, 以索引字串與資料到主要/次要資料庫中儲存資料, 使用上有以下幾點須注意
//   - 此行為結構需與泛型共同運作, 填入的泛型類型 T 需要是結構型別, 請不要填入指標型別
//   - 由於會儲存到次要資料庫中, 因此泛型類型 T 的成員都需要設定好 `bson:name` 屬性
//   - 使用前必須設定好 Field, Key 並且不能為空字串
//   - 使用前必須設定好 Data, 並且不能為nil
type Setter[T any] struct {
	Behave
	Field string           // 索引欄位
	Key   string           // 索引字串
	Data  *T               // 資料物件
	set   *redis.StatusCmd // 設值命令結果
}

// Prepare 前置處理
func (this *Setter[T]) Prepare() error {
	if this.Field == "" {
		return fmt.Errorf("setter prepare: field empty")
	} // if

	if this.Key == "" {
		return fmt.Errorf("setter prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("setter prepare: data nil")
	} // if

	value, err := json.Marshal(this.Data)

	if err != nil {
		return fmt.Errorf("setter prepare: %w", err)
	} //

	key := FormatKey(this.Key)
	this.set = this.Major().Set(this.Ctx(), key, value, 0)
	return nil
}

// Complete 完成處理
func (this *Setter[T]) Complete() error {
	value, err := this.set.Result()

	if err != nil {
		return fmt.Errorf("setter complete: %w", err)
	} // if

	if value != RedisOk {
		return fmt.Errorf("setter complete: save to redis failed")
	} // if

	field := FormatField(this.Field)
	key := FormatKey(this.Key)
	filter := bson.D{{Key: field, Value: key}}
	opt := options.Replace().SetUpsert(true)

	if _, err = this.Minor().ReplaceOne(this.Ctx(), filter, this.Data, opt); err != nil {
		return fmt.Errorf("setter complete: %w", err)
	} // if

	return nil
}

// TODO: 考慮一下怎麼讓redis的set跟mongo的upsert一起執行, 並且要等待全部執行完成後才回傳
