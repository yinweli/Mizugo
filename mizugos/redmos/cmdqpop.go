package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// QPop 佇列彈出行為
//
// 以索引鍵(Key)自主要資料庫的「佇列(List)」頭部彈出一筆元素, 並可選擇將「彈出後的佇列狀態」備份到次要資料庫;
// 若成功完成後, 會將彈出元素儲存至 Data
//
// 事前準備:
//   - (可選)設定 MinorEnable: true 表示完成時會將「彈出後的佇列快照」寫入次要資料庫
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MajorKey；若啟用備份，還需提供 MinorKey 與 MinorTable)
//   - 設定 Key: 不可為空字串
//   - (可選)設定 Data: 若為 nil, 執行時會自動建立 *T
//
// 注意:
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
//   - 若佇列過長可能造成效能問題, 建議用於小~中型佇列
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

	if this.MinorEnable {
		this.cmdQueued = this.Major().LRange(this.Ctx(), key, 0, -1)
	} // if

	return nil
}

// Complete 完成處理
func (this *QPop[T]) Complete() error {
	result, err := this.cmd.Result()

	if err != nil && errors.Is(err, redis.Nil) == false {
		return fmt.Errorf("qpop complete: %w: %v", err, this.Key)
	} // if

	if result != RedisNil {
		if this.Data == nil {
			this.Data = new(T)
		} // if

		if err = json.Unmarshal([]byte(result), this.Data); err != nil {
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
