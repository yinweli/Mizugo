package helps

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

const (
	StrNumber      = "0123456789"                              // 數字列表
	StrNumberAlpha = StrNumber + StrAlphaLower + StrAlphaUpper // 數字+小寫英文字母+大寫英文字母列表
	StrAlphaLower  = "abcdefghijklmnopqrstuvwxyz"              // 小寫英文字母列表
	StrAlphaUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"              // 大寫英文字母列表
)

// StringDisplayLength 計算字串的顯示長度
func StringDisplayLength(input string) int {
	return runewidth.StringWidth(input)
}

// StrPercentage 取得百分比字串
func StrPercentage(count, total int) string {
	return fmt.Sprintf("%.2f%%", float64(count)/float64(total)*PercentRatio100)
}
