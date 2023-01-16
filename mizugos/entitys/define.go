package entitys

import (
	"time"
)

const (
	updateInterval = time.Second // update事件間隔時間
)

// 內部事件
const (
	EventUpdate    = "update"    // update事件, 每updateInterval觸發一次
	EventDispose   = "dispose"   // dispose事件, 實體結束時執行
	EventAfterSend = "afterSend" // afterSend事件, 傳送訊息結束後執行
	EventAfterRecv = "afterRecv" // afterRecv事件, 接收訊息結束後執行
	EventReceive   = "receive"   // receive事件, 接收訊息時執行, 這個事件無法訂閱
	EventFinalize  = "finalize"  // finalize事件, 實體結束時執行, 這個事件無法訂閱
)

// EntityID 實體編號
type EntityID int64

// ModuleID 模組編號
type ModuleID int64

// Awaker awake事件介面
type Awaker interface {
	// Awake 模組初始化時第一個被執行
	Awake() error
}

// Starter start事件介面
type Starter interface {
	// Start 模組初始化時第二個被執行
	Start() error
}

// Wrong 錯誤處理函式類型
type Wrong func(err error)

// Do 執行處理
func (this Wrong) Do(err error) {
	if this != nil {
		this(err)
	} // if
}
