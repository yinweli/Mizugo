package redmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// QPeek 佇列預覽行為
//
// 以索引鍵(Key)查詢主要資料庫的「佇列(List)」完整內容, 不會改變佇列狀態;
// 最後可選擇以 Done 回呼帶出結果
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MajorKey)
//   - 設定 Key: 不可為空字串
//   - (可選)設定 Done: 完成時的回呼函式, 參數為資料列表
//
// 注意:
//   - 本行為僅使用主要資料庫, 次要資料庫不參與
//   - 泛型 T 應為「值型別的結構(struct)」, 且不要以 *T 作為型別參數
//   - 若佇列過長可能造成效能問題, 建議用於小~中型佇列
type QPeek[T any] struct {
	Behave                       // 行為物件
	Meta   Metaer                // 元資料
	Key    string                // 索引值
	Data   []*T                  // 資料列表
	Done   func(data []*T)       // 完成回呼
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

	if this.Done != nil {
		this.Done(this.Data)
	} // if

	return nil
}
