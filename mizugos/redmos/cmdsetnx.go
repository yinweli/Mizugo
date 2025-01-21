package redmos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetNX 設值行為, 當索引不存在才執行, 以索引值與資料到主要/次要資料庫中儲存資料, 使用上有以下幾點須注意
//   - 泛型類型T必須是結構, 並且不能是指標
//   - 資料結構如果包含 Save 結構或是符合 Saver 介面, 會套用儲存判斷機制, 減少不必要的儲存操作
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 執行前設定好 MajorEnable, MinorEnable
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Expire, 若不設置或是設置為0表示不逾期, 如果設為-1或是 RedisTTL 則表示不更動逾期時間
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data 並且不能為nil
type SetNX[T any] struct {
	Behave                     // 行為物件
	MajorEnable bool           // 啟用主要資料庫
	MinorEnable bool           // 啟用次要資料庫
	Meta        Metaer         // 元資料
	Expire      time.Duration  // 逾期時間, 若為0表示不逾期, 若為-1或是 RedisTTL 則表示不更動逾期時間
	Key         string         // 索引值
	Data        *T             // 資料物件
	cmd         *redis.BoolCmd // 命令結果
}

// Prepare 前置處理
func (this *SetNX[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("setnx prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("setnx prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("setnx prepare: data nil")
	} // if

	if save, ok := any(this.Data).(Saver); ok && save.GetSave() == false {
		return nil
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		data, err := json.Marshal(this.Data)

		if err != nil {
			return fmt.Errorf("setnx prepare: %w: %v", err, this.Key)
		} // if

		this.cmd = this.Major().SetNX(this.Ctx(), key, data, this.Expire)
	} // if

	if this.MinorEnable && this.Meta.MinorTable() == "" {
		return fmt.Errorf("setnx prepare: table empty")
	} // if

	return nil
}

// Complete 完成處理
func (this *SetNX[T]) Complete() error {
	if save, ok := any(this.Data).(Saver); ok && save.GetSave() == false {
		return nil
	} // if

	if this.MajorEnable {
		data, err := this.cmd.Result()

		if err != nil {
			return fmt.Errorf("setnx complete: %w: %v", err, this.Key)
		} // if

		if data == false {
			return fmt.Errorf("setnx complete: save to redis failed: %v", this.Key)
		} // if
	} // if

	if this.MinorEnable {
		key := this.Meta.MinorKey(this.Key)
		table := this.Meta.MinorTable()
		this.Minor().Operate(table, mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{MongoKey: key}).
			SetReplacement(&MinorData[T]{
				K: key,
				D: this.Data,
			}))
	} // if

	return nil
}
