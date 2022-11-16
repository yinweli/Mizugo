package entitys

import (
	"time"
)

const eventBufferSize = 1000       // 事件緩衝區設為1000, 避免因為爆滿而卡住
const eventAwake = "awake"         // awake事件, 模組初始化時第一個被執行; 參數類型為Moduler
const eventStart = "start"         // start事件, 模組初始化時第二個被執行; 參數類型為Moduler
const eventDispose = "dispose"     // dispose事件, 模組結束時執行; 參數類型為Moduler
const eventUpdate = "update"       // update事件, 模組定時事件; 參數類型為Moduler
const updateInterval = time.Second // update事件間隔時間
