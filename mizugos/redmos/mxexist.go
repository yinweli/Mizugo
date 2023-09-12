package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Exist 索引查詢行為, 以索引列表到主要資料庫中查詢索引是否存在, 使用上有以下幾點須注意
//   - 需要事先建立好與 Metaer 介面符合的元資料結構, 並填寫到 Meta
//   - 執行前設定好 Key 並且不能為空列表
//   - 執行後可用 Count 來取得存在的索引數量
type Exist struct {
	Behave               // 行為物件
	Meta   Metaer        // 元資料
	Key    []string      // 索引列表
	Count  int           // 存在的索引數量
	exist  *redis.IntCmd // 命令結果
}

// Prepare 前置處理
func (this *Exist) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("exist prepare: meta nil")
	} // if

	major, _ := this.Meta.Enable()

	if major {
		if len(this.Key) == 0 {
			return fmt.Errorf("exist prepare: key empty")
		} // if

		key := make([]string, 0, len(this.Key))

		for _, itor := range this.Key {
			key = append(key, this.Meta.MajorKey(itor))
		} // for

		this.Count = 0
		this.exist = this.Major().Exists(this.Ctx(), key...)
	} // if

	return nil
}

// Complete 完成處理
func (this *Exist) Complete() error {
	if this.Meta == nil {
		return fmt.Errorf("exist complete: meta nil")
	} // if

	major, _ := this.Meta.Enable()

	if major {
		count, err := this.exist.Result()

		if err != nil {
			return fmt.Errorf("exist complete: %w: %v", err, this.Key)
		} // if

		this.Count = int(count)
	} // if

	return nil
}
