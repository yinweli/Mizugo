package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get 取值行為, 以索引字串到主要/次要資料庫中取得資料, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 需要事先建立好與 Metaer 介面符合的元資料結構, 並填寫到 Meta
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Data 來取得資料
type Get[T any] struct {
	Behave                  // 行為物件
	Meta   Metaer           // 元資料
	Key    string           // 索引值
	Data   *T               // 資料物件
	get    *redis.StringCmd // 命令結果
}

// Prepare 前置處理
func (this *Get[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("get prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("get prepare: key empty")
	} // if

	major, minor := this.Meta.Enable()

	if major {
		key := this.Meta.MajorKey(this.Key)
		this.get = this.Major().Get(this.Ctx(), key)
	} // if

	if minor {
		if this.Meta.MinorTable() == "" {
			return fmt.Errorf("get prepare: table empty")
		} // if

		if this.Meta.MinorField() == "" {
			return fmt.Errorf("get prepare: field empty")
		} // if
	} // if

	return nil
}

// Complete 完成處理
func (this *Get[T]) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("get complete: meta nil")
	} // if

	major, minor := this.Meta.Enable()

	if major {
		data, err := this.get.Result()

		if err != nil && err != redis.Nil {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		if data != RedisNil {
			if this.Data == nil {
				this.Data = new(T)
			} // if

			if err = json.Unmarshal([]byte(data), this.Data); err != nil {
				return fmt.Errorf("get complete: %w: %v", err, this.Key)
			} // if

			minor = false // 如果主要資料庫讀取成功, 就不必執行次要資料庫了
		} // if
	} // if

	if minor {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		field := this.Meta.MinorField()
		filter := bson.D{{Key: field, Value: key}}
		result := this.Minor().Table(table).FindOne(this.Ctx(), filter)
		err := result.Err()
		empty := errors.Is(err, mongo.ErrNoDocuments)

		if err != nil && empty == false {
			return fmt.Errorf("get complete: %w: %v", err, this.Key)
		} // if

		if empty == false {
			if this.Data == nil {
				this.Data = new(T)
			} // if

			if err = result.Decode(this.Data); err != nil {
				return fmt.Errorf("get complete: %w: %v", err, this.Key)
			} // if
		} // if
	} // if

	return nil
}
