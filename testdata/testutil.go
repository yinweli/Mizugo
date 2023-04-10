package testdata

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"os"
	"time"
)

// WaitTimeout 等待超時時間, 可以輸入等待的毫秒數量, 如果未輸入毫秒數量, 則會等待200毫秒
func WaitTimeout(count ...time.Duration) {
	if len(count) > 0 {
		time.Sleep(time.Millisecond * count[0])
	} else {
		time.Sleep(time.Millisecond * 200)
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

// RandString 取得隨機字串
func RandString(length int) string {
	// 為了避免循環引用問題, 所以這裡實現了隨機字串相關機制, 而非使用 utils.RandString

	builder := bytes.Buffer{}
	max := big.NewInt(int64(len(RandStringLetter)))

	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, max)
		builder.WriteByte(RandStringLetter[int(index.Int64())])
	} // for

	return builder.String()
}
