package events

import (
    `sync/atomic`
    `time`
)

const eventSize = 10 // 事件緩衝區長度

// NewEvent 建立事件管理器
func NewEvent(process Process) *Event {
    return &Event{
        event:   make(chan any, eventSize),
        process: process,
    }
}

// Event 事件管理器
type Event struct {
    enable  atomic.Bool // 啟用旗標
    event   chan any    // 事件通道
    process Process     // 事件處理函式
}

// Process 事件處理函式類型
type Process func(event any)

// Awake awake事件
type Awake struct {
    Param any // 參數物件
}

// Start start事件
type Start struct {
    Param any // 參數物件
}

// Dispose dispose事件
type Dispose struct {
    Param any // 參數物件
}

// Update update事件
type Update struct {
    Param any // 參數物件
}

// Initialize 初始化處理
func (this *Event) Initialize() {
    go func() {
        this.enable.Store(true)

        for {
            select {
            case event := <-this.event:
                if this.process != nil {
                    this.process(event)
                } // if

            default:
                if this.enable.Load() == false {
                    return
                } // if
            } // select
        } // for
    }()
}

// Finalize 結束處理
func (this *Event) Finalize() {
    this.enable.Store(false)
}

// InvokeAwake 執行awake事件
func (this *Event) InvokeAwake(param any) {
    this.event <- &Awake{
        Param: param,
    }
}

// InvokeStart 執行start事件
func (this *Event) InvokeStart(param any) {
    this.event <- &Start{
        Param: param,
    }
}

// InvokeDispose 執行dispose事件
func (this *Event) InvokeDispose(param any) {
    this.event <- &Dispose{
        Param: param,
    }
}

// InvokeUpdate 執行update事件
//   會建立一個執行緒定時新增事件, 因此使用時要注意不能太過份
//   如果事件管理器結束時, 所有已建立的執行緒都會跟著結束
func (this *Event) InvokeUpdate(param any, interval time.Duration) {
    go func() {
        timeout := time.After(interval)

        for {
            select {
            case <-timeout:
                this.event <- &Update{
                    Param: param,
                }

            default:
                if this.enable.Load() == false {
                    return
                } // if
            } // select
        } // for
    }()
}