package redmos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetEx 設值設定逾期時間行為, 以索引字串與資料到主要/次要資料庫中儲存資料, 主要資料庫中的資料將會有逾期時間, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 資料結構如果包含 Save 結構或是符合 Saver 介面, 會套用儲存判斷機制, 減少不必要的儲存操作
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 執行前設定好 MajorEnable, MinorEnable
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Expire, 若不設置或是設置為0表示不逾期, 如果設為-1或是 KeepTTL 則表示不更動逾期時間
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行前設定好 Data 並且不能為nil
type SetEx[T any] struct {
	Behave                       // 行為物件
	MajorEnable bool             // 啟用主要資料庫
	MinorEnable bool             // 啟用次要資料庫
	Meta        Metaer           // 元資料
	Expire      time.Duration    // 逾期時間, 若為0表示不逾期, 若為-1或是 KeepTTL 則表示不更動逾期時間
	Key         string           // 索引值
	Data        *T               // 資料物件
	cmd         *redis.StatusCmd // 命令結果
}

// Prepare 前置處理
func (this *SetEx[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("setex prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("setex prepare: key empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("setex prepare: data nil")
	} // if

	if save, ok := any(this.Data).(Saver); ok && save.GetSave() == false {
		return nil
	} // if

	if this.MajorEnable {
		key := this.Meta.MajorKey(this.Key)
		data, err := json.Marshal(this.Data)

		if err != nil {
			return fmt.Errorf("setex prepare: %w: %v", err, this.Key)
		} // if

		this.cmd = this.Major().Set(this.Ctx(), key, data, this.Expire)
	} // if

	if this.MinorEnable {
		if this.Meta.MinorTable() == "" {
			return fmt.Errorf("setex prepare: table empty")
		} // if

		if this.Meta.MinorField() == "" {
			return fmt.Errorf("setex prepare: field empty")
		} // if
	} // if

	return nil
}

// Complete 完成處理
func (this *SetEx[T]) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("setex complete: meta nil")
	} // if

	if save, ok := any(this.Data).(Saver); ok && save.GetSave() == false {
		return nil
	} // if

	if this.MajorEnable {
		data, err := this.cmd.Result()

		if err != nil {
			return fmt.Errorf("setex complete: %w: %v", err, this.Key)
		} // if

		if data != RedisOk {
			return fmt.Errorf("setex complete: save to redis failed: %v", this.Key)
		} // if
	} // if

	if this.MinorEnable {
		table := this.Meta.MinorTable()
		field := this.Meta.MinorField()
		key := this.Meta.MinorKey(this.Key)
		filter := bson.D{{Key: field, Value: key}}
		this.Minor().Operate(table, mongo.NewReplaceOneModel().SetUpsert(true).SetFilter(filter).SetReplacement(this.Data))
	} // if

	return nil
}