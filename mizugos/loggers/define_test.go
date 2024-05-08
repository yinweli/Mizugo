package loggers

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

func (this *loggerTester) Get() Retain {
	return &retainTester{}
}

// retainTester 儲存測試器
type retainTester struct {
}

func (this *retainTester) Clear() Retain {
	return this
}

func (this *retainTester) Flush() Retain {
	return this
}

func (this *retainTester) Debug(_ string) Stream {
	return &streamTester{retain: this}
}

func (this *retainTester) Info(_ string) Stream {
	return &streamTester{retain: this}
}

func (this *retainTester) Warn(_ string) Stream {
	return &streamTester{retain: this}
}

func (this *retainTester) Error(_ string) Stream {
	return &streamTester{retain: this}
}

// streamTester 記錄測試器
type streamTester struct {
	retain Retain
}

func (this *streamTester) Message(_ string, _ ...any) Stream {
	return this
}

func (this *streamTester) KV(_ string, _ any) Stream {
	return this
}

func (this *streamTester) Caller(_ int, _ ...bool) Stream {
	return this
}

func (this *streamTester) Error(_ error) Stream {
	return this
}

func (this *streamTester) End() Retain {
	return this.retain
}

func (this *streamTester) EndFlush() {
}
