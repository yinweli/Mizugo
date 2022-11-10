package event

import (
    `sync/atomic`
)

// NewEvent 建立事件管理器
func NewEvent(size int) *Event {
    return &Event{
        event: make(chan any, size),
    }
}

// Event 事件管理器
type Event struct {
    end   atomic.Bool // 結束旗標
    event chan any    // 事件通道
}

// Begin 啟動事件處理
func (this *Event) Begin(process func(event any)) {
    go func() {
        for {
            select {
            case event := <-this.event:
                process(event)

            default:
                if this.end.Load() {
                    this.end.Store(false)
                    return
                } // if
            } // select
        } // for
    }()
}

// End 結束事件處理
func (this *Event) End() {
    this.end.Store(true)
}

// Add 新增事件
func (this *Event) Add(event any) {
    this.event <- event
}
