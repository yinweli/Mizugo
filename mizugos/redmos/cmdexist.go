package redmos

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Exist 查詢行為
//
// 以索引列表(Key)查詢主要資料庫中對應鍵是否存在, 並回傳存在的數量(Count)
//
// 事前準備:
//   - 設定 Meta: 需為符合 Metaer 介面的元資料物件(至少需提供 MajorKey)
//   - 設定 Key: 不可為空列表
//
// 注意:
//   - 本行為僅使用主要資料庫, 次要資料庫不參與
type Exist struct {
	Behave               // 行為物件
	Meta   Metaer        // 元資料
	Key    []string      // 索引列表
	Count  int           // 存在的索引數量
	cmd    *redis.IntCmd // 命令結果
}

// Prepare 前置處理
func (this *Exist) Prepare() error {
	if this.Meta == nil {
		return fmt.Errorf("exist prepare: meta nil")
	} // if

	if len(this.Key) == 0 {
		return fmt.Errorf("exist prepare: key empty")
	} // if

	key := make([]string, 0, len(this.Key))

	for _, itor := range this.Key {
		key = append(key, this.Meta.MajorKey(itor))
	} // for

	this.Count = 0
	this.cmd = this.Major().Exists(this.Ctx(), key...)
	return nil
}

// Complete 完成處理
func (this *Exist) Complete() error {
	count, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("exist complete: %w: %v", err, this.Key)
	} // if

	this.Count = int(count)
	return nil
}
