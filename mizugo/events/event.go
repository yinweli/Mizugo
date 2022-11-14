package events

import (
    `sync/atomic`
    `time`
)

const (
    Awake   int = iota // awake事件類型
    Start              // start事件類型
    Update             // update事件類型
    Dispose            // dispose事件類型
)

// NewEvent 建立事件管理器
func NewEvent(size int) *Event {
    return &Event{
        event: make(chan Data, size),
    }
}

// Event 事件管理器
type Event struct {
    enable atomic.Bool // 啟用旗標
    event  chan Data   // 事件通道
}

// Data 事件資料
type Data struct {
    Type  int // 事件類型
    Param any // 事件參數
}

// Initialize 初始化處理
func (this *Event) Initialize(interval time.Duration, process func(data Data)) {
    go func() {
        this.enable.Store(true)
        timeout := time.After(interval)

        for {
            select {
            case data := <-this.event:
                process(data)

            case <-timeout:
                process(Data{
                    Type: Update,
                })

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

// Execute 執行事件
func (this *Event) Execute(type_ int, param any) {
    this.event <- Data{
        Type:  type_,
        Param: param,
    }
}
