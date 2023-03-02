package redmos

import (
	"strings"
)

// FormatField 格式化欄位, 把輸入的索引轉為小寫
func FormatField(field string) string {
	return strings.ToLower(field)
}

// FormatKey 格式化索引, 把輸入的多個索引用':'連接起來, 並且轉為小寫
func FormatKey(key ...string) string {
	return strings.ToLower(strings.Join(key, ":"))
}
