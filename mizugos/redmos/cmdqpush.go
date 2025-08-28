package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// QPush 佇列推入行為
//
// 以索引鍵(Key)將一筆資料(Data)推入主要資料庫的「佇列(List)」尾端, 並可選擇將「推入後的佇列快照」備份至次要資料庫
//
// 事前準備:
//   - (可選)設定 MinorEnable: true 表示在完成時會將「推入後的佇列快照」寫入次要資料庫
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MajorKey; 若啟用備份, 還需提供 MinorKey 與 MinorTable)
//   - 設定 Key: 不可為空字串
//   - 設定 Data: 不可為 nil
//
// 注意:
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
//   - 若佇列過長可能造成效能問題, 建議用於小~中型佇列
type QPush[T any] struct {
	Behave                            // 行為物件
	MinorEnable bool                  // 啟用次要資料庫
	Meta        Metaer                // 元資料
	Key         string                // 索引值
	Data        *T                    // 資料物件
	cmd         *redis.IntCmd         // 命令結果
	cmdQueued   *redis.StringSliceCmd // 佇列內容結果
}

// QueueData 佇列資料, 作為次要資料庫備份時的文件載體
type QueueData[T any] struct {
	Data []*T `bson:"value"` // 佇列列表
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

	if this.MinorEnable {
		this.cmdQueued = this.Major().LRange(this.Ctx(), key, 0, -1)
	} // if

	return nil
}

// Complete 完成處理
func (this *QPush[T]) Complete() error {
	_, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("qpush complete: %w: %v", err, this.Key)
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

			queue.Data = append(queue.Data, d)
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
