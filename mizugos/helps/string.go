package helps

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	StrNumber      = "0123456789"                              // 數字列表
	StrNumberAlpha = StrNumber + StrAlphaLower + StrAlphaUpper // 數字+小寫英文字母+大寫英文字母列表
	StrAlphaLower  = "abcdefghijklmnopqrstuvwxyz"              // 小寫英文字母列表
	StrAlphaUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"              // 大寫英文字母列表
	StrPercent     = 100                                       // 百分比乘數
)

// StringJoin 以分隔字串來組合字串列表, 實際上就是使用了 strings.Join 只是轉換了參數順序而已
func StringJoin(sep string, element ...string) string {
	return strings.Join(element, sep)
}

// StringDisplayLength 計算字串的顯示長度, 計算的標準如下
//   - 中日韓字: 2顯示長度
//   - 英數及其他字: 1顯示長度
func StringDisplayLength(input string) int {
	length := 0

	for _, r := range input {
		if unicode.Is(unicode.Han, r) ||
			unicode.Is(unicode.Hangul, r) ||
			unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Hiragana, r) {
			length += 2
		} else {
			length++
		} // if
	} // for

	return length
}

// StrPercentage 取得百分比字串
func StrPercentage(count, total int) string {
	return fmt.Sprintf("%.2f%%", float64(count)/float64(total)*StrPercent)
}
