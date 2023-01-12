package testdata

import (
	"go.uber.org/goleak"
)

// TestLeak 洩漏測試
type TestLeak struct {
}

// GoLeak 執行執行緒洩漏測試
func (this *TestLeak) GoLeak(t goleak.TestingT, run bool) {
	if runGoLeak && run {
		goleak.VerifyNone(t)
	} // if
}

func init() {
	runGoLeak = false
}

var runGoLeak bool // 執行緒洩漏測試旗標
