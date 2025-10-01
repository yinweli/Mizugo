package redmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
)

// QPopAll 佇列全部彈出行為
//
// 以索引鍵(Key)將主要資料庫中的「佇列(List)」取出全部元素並清空佇列, 並可選擇將「彈出後的佇列狀態」備份到次要資料庫;
// 若成功完成後, 會將彈出元素儲存至 Data, 最後可選擇以 Done 回呼帶出結果
//
// 事前準備:
//   - (可選)設定 MinorEnable: true 表示完成時會將「彈出後的佇列快照」寫入次要資料庫
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MajorKey；若啟用備份，還需提供 MinorKey 與 MinorTable)
//   - 設定 Key: 不可為空字串
//   - (可選)設定 Done: 完成時的回呼函式, 參數為資料列表
//
// 注意:
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
//   - 若佇列過長可能造成效能問題, 建議用於小~中型佇列
type QPopAll[T any] struct {
	Behave                            // 行為物件
	MinorEnable bool                  // 啟用次要資料庫
	Meta        Metaer                // 元資料
	Key         string                // 索引值
	Data        []*T                  // 資料列表
	Done        func(data []*T)       // 完成回呼
	cmd         *redis.StringSliceCmd // 命令結果
	cmdRename   *redis.StatusCmd      // 更名命令結果
	cmdDelete   *redis.IntCmd         // 刪除命令結果
}

// Prepare 前置處理
func (this *QPopAll[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("qpopall prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("qpopall prepare: key empty")
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("qpopall prepare: table empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	rename := fmt.Sprintf("%v-%v-%v", key, time.Now().Unix(), helps.RandStringDefault())
	this.cmdRename = this.Major().Rename(this.Ctx(), key, rename)
	this.cmd = this.Major().LRange(this.Ctx(), rename, 0, -1)
	this.cmdDelete = this.Major().Del(this.Ctx(), rename)
	return nil
}

// Complete 完成處理
func (this *QPopAll[T]) Complete() error {
	if _, err := this.cmdRename.Result(); err != nil {
		return fmt.Errorf("qpopall complete: rename failed: %w: %v", err, this.Key)
	} // if

	if _, err := this.cmdDelete.Result(); err != nil {
		return fmt.Errorf("qpopall complete: delete failed: %w: %v", err, this.Key)
	} // if

	result, err := this.cmd.Result()

	if err != nil && errors.Is(err, redis.Nil) == false {
		return fmt.Errorf("qpopall complete: %w: %v", err, this.Key)
	} // if

	for _, itor := range result {
		data := new(T)

		if err = json.Unmarshal([]byte(itor), data); err != nil {
			return fmt.Errorf("qpopall complete: %w: %v", err, this.Key)
		} // if

		this.Data = append(this.Data, data)
	} // for

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{MongoKey: key}).
			SetReplacement(&MinorData[QueueData[T]]{
				K: key,
				D: &QueueData[T]{}, // 將次要資料庫中的儲存資料清空
			}))
	} // if

	if this.Done != nil {
		this.Done(this.Data)
	} // if

	return nil
}
