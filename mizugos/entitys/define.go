package entitys

import (
	"time"
)

const eventSize = 1000             // 事件通道大小設為1000, 避免因為爆滿而卡住
const updateInterval = time.Second // update事件間隔時間

// 預設事件名稱
const (
	EventUpdate    = "update"    // update事件, 模組定時事件
	EventDispose   = "dispose"   // dispose事件, 模組結束時執行
	EventAfterSend = "afterSend" // afterSend事件, 傳送訊息結束後執行
	EventAfterRecv = "afterRecv" // afterRecv事件, 接收訊息結束後執行
	EventFinalize  = "finalize"  // finalize事件, 實體結束時執行; 這個事件只能實體內部使用
)

// EntityID 實體編號
type EntityID int64

// ModuleID 模組編號
type ModuleID int64

// Awaker awake事件介面
type Awaker interface {
	// Awake 模組初始化時第一個被執行
	Awake()
}

// Starter start事件介面
type Starter interface {
	// Start 模組初始化時第二個被執行
	Start()
}
