package redmos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetNX 設值行為, 當索引不存在才執行, 以索引字串與資料到主要/次要資料庫中儲存資料, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 資料結構如果符合 Saver 介面, 會套用儲存判斷機制, 減少不必要的儲存操作
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 需要事先建立好與 Metaer 介面符合的元資料結構, 並填寫到 Meta
//   - 執行前設定好 Expire, 若不設置或是設置為0表示不過期
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data 並且不能為nil
type SetNX[T any] struct {
	Behave                     // 行為物件
	MajorEnable bool           // 啟用主要資料庫
	MinorEnable bool           // 啟用次要資料庫
	Meta        Metaer         // 元資料
	Expire      time.Duration  // 過期時間, 若為0表示不過期
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

	if save, ok := any(this.Data).(Saver); ok && save.Save() == false {
		return nil
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		data, err := json.Marshal(this.Data)

		if err != nil {
			return fmt.Errorf("setnx prepare: %w: %v", err, this.Key)
		} //

		this.cmd = this.Major().SetNX(this.Ctx(), key, data, this.Expire)
	} // if

	if this.MinorEnable {
		if this.Meta.MinorTable() == "" {
			return fmt.Errorf("setnx prepare: table empty")
		} // if

		if this.Meta.MinorField() == "" {
			return fmt.Errorf("setnx prepare: field empty")
		} // if
	} // if

	return nil
}

// Complete 完成處理
func (this *SetNX[T]) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("setnx complete: meta nil")
	} // if

	if save, ok := any(this.Data).(Saver); ok && save.Save() == false {
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
		field := this.Meta.MinorField()
		filter := bson.D{{Key: field, Value: key}}
		opt := options.Replace().SetUpsert(true)

		if _, err := this.Minor().Table(table).ReplaceOne(this.Ctx(), filter, this.Data, opt); err != nil {
			return fmt.Errorf("setnx complete: %w: %v", err, this.Key)
		} // if
	} // if

	return nil
}
