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

// Awaker awake介面
type Awaker interface {
	// Awake awake事件, 模組初始化時第一個被執行
	Awake()
}

// Starter start介面
type Starter interface {
	// Start start事件, 模組初始化時第二個被執行
	Start()
}

// Disposer dispose介面
type Disposer interface {
	// Dispose dispose事件, 模組結束時執行
	Dispose()
}

// Updater update介面
type Updater interface {
	// Update update事件, 模組定時事件, 間隔時間定義在updateInterval
	Update()
}
