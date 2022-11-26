package utils

import (
	"time"
)

// NewWaitTimeout 建立逾時等待器
func NewWaitTimeout(timeout time.Duration) *WaitTimeout {
	return &WaitTimeout{
		timeout: timeout,
		done:    make(chan any, 1),
	}
}

// WaitTimeout 逾時等待器, 只能一次等待一個工作, 不能像sync.WaitGroup可以等待複數工作
type WaitTimeout struct {
	timeout time.Duration // 逾時時間
	done    chan any      // 結束通道
}

// Wait 等待處理, 回傳true表示未逾時, false則否
func (this *WaitTimeout) Wait() bool {
	select {
	case <-this.done:
		return true

	case <-time.After(this.timeout):
		return false
	} // select
}

// Done 等待結束
func (this *WaitTimeout) Done() {
	this.done <- nil
}
