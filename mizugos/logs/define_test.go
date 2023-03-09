package logs

import (
	"fmt"
)

// newLoggerTester 建立日誌測試器
func newLoggerTester(initialize bool) *loggerTester {
	return &loggerTester{
		initialize: initialize,
	}
}

// loggerTester 日誌測試器
type loggerTester struct {
	initialize bool
}

func (this *loggerTester) Initialize() error {
	if this.initialize {
		return nil
	} else {
		return fmt.Errorf("failed")
	} // if
}

func (this *loggerTester) Finalize() {
}

func (this *loggerTester) Debug(_ string) Stream {
	return &streamTester{}
}

func (this *loggerTester) Info(_ string) Stream {
	return &streamTester{}
}

func (this *loggerTester) Warn(_ string) Stream {
	return &streamTester{}
}

func (this *loggerTester) Error(_ string) Stream {
	return &streamTester{}
}

// streamTester 記錄測試器
type streamTester struct {
}

func (this *streamTester) Message(_ string, _ ...any) Stream {
	return this
}

func (this *streamTester) Caller(_ int) Stream {
	return this
}

func (this *streamTester) KV(_ string, _ any) Stream {
	return this
}

func (this *streamTester) Error(_ error) Stream {
	return this
}

func (this *streamTester) EndError(_ error) {
}

func (this *streamTester) End() {
}
