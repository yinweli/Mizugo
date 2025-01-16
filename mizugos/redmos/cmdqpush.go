package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// QPush 推入佇列行為, 以索引值與資料到主要資料庫中儲存佇列, 使用上有以下幾點須注意
//   - 泛型類型T必須是結構, 並且不能是指標
//   - 執行前設定好 MinorEnable; 由於佇列行為只會在主要資料庫中執行, 因此次要資料庫僅用於備份
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data 並且不能為nil
type QPush[T any] struct {
	Behave                            // 行為物件
	MinorEnable bool                  // 啟用次要資料庫
	Meta        Metaer                // 元資料
	Key         string                // 索引值
	Data        *T                    // 資料物件
	cmd         *redis.IntCmd         // 命令結果
	cmdQueued   *redis.StringSliceCmd // 佇列內容結果
}

// QueueData 佇列資料
type QueueData[T any] struct {
	Value []*T `bson:"value"` // 內容列表
}

// Prepare 前置處理
func (this *QPush[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("qpush prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("qpush prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("qpush prepare: data nil")
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("qpush prepare: table empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	data, err := json.Marshal(this.Data)

	if err != nil {
		return fmt.Errorf("qpush prepare: %w: %v", err, this.Key)
	} // if

	this.cmd = this.Major().RPush(this.Ctx(), key, data)
	this.cmdQueued = this.Major().LRange(this.Ctx(), key, 0, -1)
	return nil
}

// Complete 完成處理
func (this *QPush[T]) Complete() error {
	count, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("qpush complete: %w: %v", err, this.Key)
	} // if

	if count == 0 {
		return fmt.Errorf("qpush complete: save to redis failed: %v", this.Key)
	} // if

	if this.MinorEnable {
		list, err := this.cmdQueued.Result()

		if err != nil && errors.Is(err, redis.Nil) == false {
			return fmt.Errorf("qpush complete: %w: %v", err, this.Key)
		} // if

		queue := &QueueData[T]{}

		for _, itor := range list {
			d := new(T)

			if err = json.Unmarshal([]byte(itor), d); err != nil {
				return fmt.Errorf("qpush complete: %w: %v", err, this.Key)
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
