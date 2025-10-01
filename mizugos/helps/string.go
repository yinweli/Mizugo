package helps

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

const (
	StrAlphaLower  = "abcdefghijklmnopqrstuvwxyz"              // 小寫英文字母列表
	StrAlphaUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"              // 大寫英文字母列表
	StrNumber      = "0123456789"                              // 數字列表
	StrNumberAlpha = StrAlphaLower + StrAlphaUpper + StrNumber // 小寫英文字母 + 大寫英文字母 + 數字列表
)

// StringDisplayLength 計算字串的顯示長度
func StringDisplayLength(input string) int {
	return runewidth.StringWidth(input)
}

// StrPercentage 取得百分比字串(含 % 符號, 保留兩位小數)
func StrPercentage(count, total int) string {
	if total == 0 {
		return "0.00%"
	} // if

	return fmt.Sprintf("%.2f%%", float64(count)/float64(total)*float64(PercentRatio100))
}
