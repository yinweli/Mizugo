package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// QPeek 取得佇列內容行為, 以索引值到主要資料庫中取得佇列內容, 但是不會改變佇列, 使用上有以下幾點須注意
//   - 泛型類型T必須是結構, 並且不能是指標
//   - 執行前設定好 Meta, 這需要事先建立好與 Metaer 介面符合的元資料結構
//   - 執行前設定好 Key 並且不能為空字串
//   - 執行後可用 Data 來取得資料列表
type QPeek[T any] struct {
	Behave                       // 行為物件
	Meta   Metaer                // 元資料
	Key    string                // 索引值
	Data   []*T                  // 資料列表
	cmd    *redis.StringSliceCmd // 命令結果
}

// Prepare 前置處理
func (this *QPeek[T]) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("qpeek prepare: meta nil")
	} // if

	if this.Key == "" {
		return fmt.Errorf("qpeek prepare: key empty")
	} // if

	key := this.Meta.MajorKey(this.Key)
	this.cmd = this.Major().LRange(this.Ctx(), key, 0, -1)
	return nil
}

// Complete 完成處理
func (this *QPeek[T]) Complete() error {
	result, err := this.cmd.Result()

	if err != nil && errors.Is(err, redis.Nil) == false {
		return fmt.Errorf("qpeek complete: %w: %v", err, this.Key)
	} // if

	for _, itor := range result {
		data := new(T)

		if err = json.Unmarshal([]byte(itor), data); err != nil {
			return fmt.Errorf("qpeek complete: %w: %v", err, this.Key)
		} // if

		this.Data = append(this.Data, data)
	} // for

	return nil
}
