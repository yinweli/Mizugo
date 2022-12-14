package events

import (
	"sync/atomic"
)

// Fixed 定時控制器
type Fixed struct {
	finish atomic.Bool // 結束旗標
}

// Stop 停止定時事件
func (this *Fixed) Stop() {
	this.finish.Store(true)
}

// State 取得定時狀態
func (this *Fixed) State() bool {
	return this.finish.Load()
}
