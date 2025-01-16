package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get 取值行為, 以索引值到主要/次要資料庫中取得資料, 不會影響主要資料庫中的資料的逾期時間, 使用上有以下幾點須注意
//   - 泛型類型T必須是結構, 並且不能是指標
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 執行前設定好 MajorEnable, MinorEnable
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Data 來取得資料
type Get[T any] struct {
	Behave                       // 行為物件
	MajorEnable bool             // 啟用主要資料庫
	MinorEnable bool             // 啟用次要資料庫
	Meta        Metaer           // 元資料
	Key         string           // 索引值
	Data        *T               // 資料物件
	cmd         *redis.StringCmd // 命令結果
}

// Prepare 前置處理
func (this *Get[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("get prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("get prepare: key empty")
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		this.cmd = this.Major().Get(this.Ctx(), key)
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("get prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *Get[T]) Complete() error {
	if this.MajorEnable {
		data, err := this.cmd.Result()

		if err != nil && errors.Is(err, redis.Nil) == false {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		if data != RedisNil {
			if this.Data == nil {
				this.Data = new(T)
			} // if

			if err = json.Unmarshal([]byte(data), this.Data); err != nil {
				return fmt.Errorf("get complete: %w: %v", err, this.Key)
			} // if

			return nil // 如果主要資料庫讀取成功, 就不必執行次要資料庫了
		} // if
	} // if

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		result := this.Minor().Collection(table).FindOne(this.Ctx(), bson.M{MongoKey: key})
		err := result.Err()
		empty := errors.Is(err, mongo.ErrNoDocuments)

		if err != nil && empty == false {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		if empty == false {
			if this.Data == nil {
				this.Data = new(T)
			} // if

			if err = result.Decode(&MinorData[T]{
				D: this.Data,
			}); err != nil {
				return fmt.Errorf("get complete: %w: %v", err, this.Key)
			} // if
		} // if
	} // if

	return nil
}
