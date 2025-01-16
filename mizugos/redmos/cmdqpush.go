package redmos

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// QPush 推入佇列行為, 以索引值與資料到主要/次要資料中儲存佇列, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 執行前設定好 MajorEnable, MinorEnable
//   - 執行前設定好 Name 並且不能為空字串
//   - 執行前設定好 Data 並且不能為nil
type QPush[T any] struct {
	Behave                    // 行為物件
	MajorEnable bool          // 啟用主要資料庫
	MinorEnable bool          // 啟用次要資料庫
	Name        string        // 佇列名稱, 同時會是主要資料庫中的索引值與次要資料庫中的表格名稱
	Data        *T            // 資料物件
	cmd         *redis.IntCmd // 命令結果
}

// Prepare 前置處理
func (this *QPush[T]) Prepare() error {
	if this.Name == "" {
		return fmt.Errorf("qpush prepare: name empty")
	} // if

	if this.Data == nil {
		return fmt.Errorf("qpush prepare: data nil")
	} // if

	if this.MajorEnable {
		data, err := json.Marshal(this.Data)

		if err != nil {
			return fmt.Errorf("qpush prepare: %w: %v", err, this.Name)
		} // if

		this.cmd = this.Major().RPush(this.Ctx(), this.Name, data)
	} // if

	return nil
}

// Complete 完成處理
func (this *QPush[T]) Complete() error {
	if this.MajorEnable {
		count, err := this.cmd.Result()

		if err != nil {
			return fmt.Errorf("qpush complete: %w: %v", err, this.Name)
		} // if

		if count == 0 {
			return fmt.Errorf("qpush complete: save to redis failed: %v", this.Name)
		} // if
	} // if

	if this.MinorEnable {
		doc, ok := any(this.Data).(bson.M)

		if ok == false {
			raw, err := bson.Marshal(this.Data)

			if err != nil {
				return fmt.Errorf("qpush complete: failed to marshal data to BSON: %w: %v", err, this.Name)
			} // if

			if err = bson.Unmarshal(raw, &doc); err != nil {
				return fmt.Errorf("qpush complete: failed to unmarshal BSON to bson.M: %w: %v", err, this.Name)
			} // if
		} // if

		doc[QPushTime] = time.Now().Unix()
		this.Minor().Operate(this.Name, mongo.NewInsertOneModel().SetDocument(doc))
	} // if

	return nil
}
