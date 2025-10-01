package redmos

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Keys 索引搜尋行為
//
// 以匹配字串(Pattern)在主要資料庫中搜尋符合條件的鍵, 並回傳結果列表(Data)
//
// 事前準備:
//   - 設定 Pattern: 不可為空字串, 使用 Redis 的全域匹配語法
//
// 注意:
//   - 本行為僅使用主要資料庫, 次要資料庫不參與
//   - 內部以 Redis `KEYS` 指令實作, 時間複雜度 O(N); 在鍵數量龐大時可能造成阻塞, 不建議於生產環境大量使用
//   - 執行成功後, 結果會寫入 Data; 當無匹配結果時回傳空列表
//
// Pattern 匹配規則(節錄):
//   - `*`  : 匹配任意長度(含 0)的字元, 例 `user:*`
//   - `?`  : 匹配任意單一字元, 例 `user:???` → `user:abc`
//   - `[]` : 字元集合/範圍, 例 `user:[a-c]*` → `user:a123`/`b456`
//   - `[^]`: 否定集合, 例 `user:[^abc]*` → 第一個字元非 a/b/c
//   - `\`  : 逸出特殊字元, 例 `user:\*` 僅匹配鍵名 `user:*`
//
// 範例:
//   - `user:*`      : 匹配 `user:` 開頭的鍵
//   - `user:?*`     : 匹配 `user:` 後至少還有 1 個字元的鍵
//   - `sess:[0-9]*` : 匹配 `sess:` 後為數字的鍵
type Keys struct {
	Behave                        // 行為物件
	Pattern string                // 匹配字串
	Data    []string              // 資料物件
	cmd     *redis.StringSliceCmd // 命令結果
}

// Prepare 前置處理
func (this *Keys) Prepare() error {
	if this.Pattern == "" {
		return fmt.Errorf("keys prepare: pattern empty")
	} // if

	this.cmd = this.Major().Keys(this.Ctx(), this.Pattern)
	return nil
}

// Complete 完成處理
func (this *Keys) Complete() error {
	result, err := this.cmd.Result()

	if err != nil && errors.Is(err, redis.Nil) == false {
		return fmt.Errorf("keys complete: %w: %v", err, this.Pattern)
	} // if

	this.Data = result
	return nil
}
