package entitys

import (
	"github.com/yinweli/Mizugo/v2/mizugos/nets"
)

const ( // 內部事件名稱
	EventCapacity = 1000           // 事件容量
	EventDispose  = "dispose"      // 結束事件, Finalize 時首先發布, 參數為 nil
	EventShutdown = "shutdown"     // 關閉事件, Finalize 時第二個發布, 參數為 nil (此時連線已中斷)
	EventRecv     = nets.EventRecv // 接收訊息事件, 當接收訊息後觸發, 參數是訊息物件
	EventSend     = nets.EventSend // 傳送訊息事件, 當傳送訊息後觸發, 參數是訊息物件
)

// EntityID 實體編號
type EntityID = uint64

// ModuleID 模組編號
type ModuleID = uint64

// Wrong 錯誤處理函式類型
type Wrong func(err error)

// Do 執行處理
func (this Wrong) Do(err error) {
	if this != nil {
		this(err)
	} // if
}
