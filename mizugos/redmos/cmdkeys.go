package redmos

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Keys 搜尋行為, 以匹配字串到主要資料庫中取得索引, 使用上有以下幾點須注意
//   - 執行前設定好 Pattern 並且不能為空字串
//   - 執行後可用 Data 來取得資料
//
// # Pattern匹配規則
//
// `*`: 匹配任意數量的字符(包括零個字符)
//   - 模式 `user:*` 匹配所有以 `user:` 開頭的鍵
//   - 模式 `*` 匹配所有鍵
//
// `?`: 匹配任意一個字符
//   - 模式 `user:???` 匹配所有以 `user:` 開頭, 後面跟三個字符的鍵, 如 `user:abc`
//
// `[]`: 匹配括號內的任意一個字符, 可以使用範圍來指定字符, 例如 [a-z] 匹配所有小寫字母
//   - 模式 `user:[abc]*` 匹配所有以 `user:` 開頭, 且接著是 a、b、或 c 的鍵, 如 `user:a123`、`user:b456`
//
// `[^]`: 匹配不在括號內的任意一個字符(否定模式)
//   - 模式 `user:[^abc]*` 匹配所有以 `user:` 開頭, 且後面跟的第一個字符不是 a、b、或 c 的鍵, 如 `user:d123`
//
// `\`: 用於轉義特殊字符, 使其作為普通字符匹配
//   - 模式 `user:\*` 匹配鍵 `user:*`, 而不是所有以 `user:` 開頭的鍵
//
// # 其他範例
//   - `user:*`:匹配所有以 `user:` 開頭的鍵
//   - `user:?*`:匹配所有以 `user:` 開頭, 且後面至少有一個字符的鍵
//   - `user:[abc]*`:匹配所有以 `user:` 開頭, 且接著是 a、b、或 c 的鍵
//   - `user:[^abc]*`:匹配所有以 `user:` 開頭, 且接著的第一個字符不是 a、b、或 c 的鍵
//
// # 注意事項
//   - Keys 命令的時間複雜度為O(N), 其中N是主要資料庫中的鍵的數量, 由於這個原因, Keys 命令在大數據集上運行時可能會導致性能問題, 因此在生產環境中應謹慎使用
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
	data, err := this.cmd.Result()

	if err != nil && errors.Is(err, redis.Nil) == false {
		return fmt.Errorf("keys complete: %w: %v", err, this.Pattern)
	} // if

	this.Data = data
	return nil
}
