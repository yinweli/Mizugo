package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// QPop 彈出佇列行為, 以索引值到主要資料庫中取得佇列, 使用上有以下幾點須注意
//   - 泛型類型T必須是結構, 並且不能是指標
//   - 執行前設定好 MinorEnable; 由於佇列行為只會在主要資料庫中執行, 因此次要資料庫僅用於備份
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Data 來取得資料
type QPop[T any] struct {
	Behave                            // 行為物件
	MinorEnable bool                  // 啟用次要資料庫
	Meta        Metaer                // 元資料
	Key         string                // 索引值
	Data        *T                    // 資料物件
	cmd         *redis.StringCmd      // 命令結果
	cmdQueued   *redis.StringSliceCmd // 佇列內容結果
}

// Prepare 前置處理
func (this *QPop[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("qpop prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("qpop prepare: key empty")
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("qpop prepare: table empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	this.cmd = this.Major().LPop(this.Ctx(), key)
	this.cmdQueued = this.Major().LRange(this.Ctx(), key, 0, -1)
	return nil
}

// Complete 完成處理
func (this *QPop[T]) Complete() error {
	data, err := this.cmd.Result()

	if err != nil && errors.Is(err, redis.Nil) == false {
		return fmt.Errorf("qpop complete: %w: %v", err, this.Key)
	} // if

	if data != RedisNil {
		if this.Data == nil {
			this.Data = new(T)
		} // if

		if err = json.Unmarshal([]byte(data), this.Data); err != nil {
			return fmt.Errorf("qpop complete: %w: %v", err, this.Key)
		} // if
	} // if

	if this.MinorEnable {
		list, err := this.cmdQueued.Result()

		if err != nil && errors.Is(err, redis.Nil) == false {
			return fmt.Errorf("qpop complete: %w: %v", err, this.Key)
		} // if

		queue := &QueueData[T]{}

		for _, itor := range list {
			d := new(T)

			if err = json.Unmarshal([]byte(itor), d); err != nil {
				return fmt.Errorf("qpop complete: %w: %v", err, this.Key)
			} // if

			queue.Value = append(queue.Value, d)
		} // for

		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{MongoKey: key}).
			SetReplacement(&MinorData[QueueData[T]]{
				K: key,
				D: queue,
			}))
	} // if

	return nil
}
