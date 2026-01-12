package helps

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

const (
	StrNumber      = "0123456789"                              // 數字列表
	StrAlphaUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"              // 大寫英文字母列表
	StrAlphaLower  = "abcdefghijklmnopqrstuvwxyz"              // 小寫英文字母列表
	StrNumberAlpha = StrNumber + StrAlphaUpper + StrAlphaLower // 數字列表 + 大寫英文字母列表 + 小寫英文字母列表
)

// StringDisplayLength 計算字串的顯示長度
func StringDisplayLength(input string) int {
	return runewidth.StringWidth(input)
}

// StringPercentage 取得百分比字串(含 % 符號, 保留兩位小數)
func StringPercentage(count, total int) string {
	if total == 0 {
		return "0.00%"
	} // if

	return fmt.Sprintf("%.2f%%", float64(count)/float64(total)*float64(PercentRatio100))
}

// StringDuration 解析包含時間單位的字串並轉換為 time.Duration
//
// 字串由一個或多個「數值」與「單位」的配對組成, 配對之間可包含空白, 若字串開頭為 "-" 則表示負的時長
//
// 支援的時間單位:
//   - d  → 天 (24 小時)
//   - h  → 小時
//   - m  → 分鐘
//   - s  → 秒
//   - ms → 毫秒
//
// 範例:
//
//	"1d"        => 24 小時
//	"2h30m"     => 2 小時 30 分
//	"1d 2h 30m" => 26 小時 30 分 (允許空白分隔)
//	"-10s"      => 負 10 秒
//	"500ms"     => 500 毫秒
//	"1H 30M"    => 1 小時 30 分 (支援大寫)
//
// 注意: 數值必須為整數(因為使用了 regex `\d+`), 不支援小數點(例如 "1.5h" 無效, 需寫成 "1h 30m")
func StringDuration(input string) (result time.Duration, err error) {
	s := strings.TrimSpace(input)
	sign := time.Duration(1)

	if s != "" && s[0] == '-' {
		s = strings.TrimSpace(s[1:])
		sign = -1
	} // if

	for s != "" {
		loc := durationCompile.FindStringSubmatchIndex(s)

		if loc == nil {
			if remain := strings.TrimSpace(s); remain != "" { // 如果剩下的不是空白, 表示解析失敗
				return 0, fmt.Errorf("duration: invalid token %q", remain)
			} // if
		} // if

		value := s[loc[2]:loc[3]]
		number, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return 0, fmt.Errorf("duration: invalid number %q: %w", value, err)
		} // if

		unit := strings.ToLower(s[loc[4]:loc[5]])

		switch unit {
		case "d":
			result += time.Duration(number) * 24 * time.Hour

		case "h":
			result += time.Duration(number) * time.Hour

		case "m":
			result += time.Duration(number) * time.Minute

		case "s":
			result += time.Duration(number) * time.Second

		case "ms":
			result += time.Duration(number) * time.Millisecond

		default: // 理論上正則表達式限制了單位, 這裡幾乎不會發生, 但保留以防萬一
			return 0, fmt.Errorf("duration: invalid unit %q", unit)
		} // switch

		s = s[loc[1]:]
	} // for

	return result * sign, nil
}

var durationCompile = regexp.MustCompile(`(?i)^\s*(\d+)\s*(ms|s|m|h|d)`) // 時長正則表達式
