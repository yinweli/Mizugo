package testdata

import (
	"go.uber.org/goleak"
)

// Leak 執行洩漏測試, 測試是否有執行緒未被關閉, 但是會有誤判的狀況;
// 依靠 leak 來決定是否要執行測試
func Leak(t goleak.TestingT, test bool) {
	if leak && test {
		goleak.VerifyNone(t)
	} // if
}

var leak = false // 洩漏測試旗標, 預設為關閉
