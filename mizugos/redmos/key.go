package redmos

import (
	"strings"
)

// MajorKey 取得主要資料庫索引, 會把輸入的多個索引用':'連接起來
func MajorKey(key ...string) string {
	return strings.Join(key, ":")
}
