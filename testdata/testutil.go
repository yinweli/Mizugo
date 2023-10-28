package testdata

import (
	"os"
	"time"
)

// WaitTimeout 等待超時時間, 可以輸入等待時間, 預設等待200毫秒
func WaitTimeout(duration ...time.Duration) {
	if len(duration) > 0 {
		time.Sleep(duration[0])
	} else {
		time.Sleep(Timeout)
	} // if
}

// CompareFile 比對檔案內容, 預期資料來自位元陣列
func CompareFile(path string, expected []byte) bool {
	if actual, err := os.ReadFile(path); err == nil {
		if string(expected) == string(actual) {
			return true
		} // if
	} // if

	return false
}
