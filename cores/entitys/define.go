package entitys

import (
	"time"
)

const eventSize = 1000             // 事件通道大小設為1000, 避免因為爆滿而卡住
const updateInterval = time.Second // update事件間隔時間

// 內部事件名稱
const (
	eventAwake   = "awake"   // awake事件, 模組初始化時第一個被執行; 參數類型為Moduler
	eventStart   = "start"   // start事件, 模組初始化時第二個被執行; 參數類型為Moduler
	eventDispose = "dispose" // dispose事件, 模組結束時執行; 參數類型為Moduler
	eventUpdate  = "update"  // update事件, 模組定時事件; 參數類型為Moduler
)

// EntityID 實體編號
type EntityID int64

// ModuleID 模組編號
type ModuleID int64

// Moduler 模組介面
type Moduler interface {
	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Entity 取得實體物件
	Entity() *Entity

	// Internal 取得內部物件
	Internal() *Internal
}

// Awaker awake模組事件介面
type Awaker interface {
	// Awake 模組初始化時第一個被執行
	Awake()
}

// Starter start模組事件介面
type Starter interface {
	// Start 模組初始化時第二個被執行
	Start()
}

// Disposer dispose模組事件介面
type Disposer interface {
	// Dispose 模組結束時執行
	Dispose()
}

// Updater update模組事件介面
type Updater interface {
	// Update 模組定時執行, 間隔時間定義在updateInterval
	Update()
}
