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

// QPopAll 彈出全部佇列行為, 以索引值到主要資料庫中取得全部佇列, 使用上有以下幾點須注意
//   - 泛型類型T必須是結構, 並且不能是指標
//   - 執行前設定好 MinorEnable; 由於佇列行為只會在主要資料庫中執行, 因此次要資料庫僅用於備份
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Data 來取得資料列表
type QPopAll[T any] struct {
	Behave                            // 行為物件
	MinorEnable bool                  // 啟用次要資料庫
	Meta        Metaer                // 元資料
	Key         string                // 索引值
	Data        *QueueData[T]         // 資料物件
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

	if this.Data == nil {
		this.Data = new(QueueData[T])
	} // if

	for _, itor := range result {
		data := new(T)

		if err = json.Unmarshal([]byte(itor), data); err != nil {
			return fmt.Errorf("qpopall complete: %w: %v", err, this.Key)
		} // if

		this.Data.Data = append(this.Data.Data, data)
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

	return nil
}
