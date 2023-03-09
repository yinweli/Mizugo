package entitys

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos/nets"
)

const ( // 內部事件名稱
	EventUpdate   = "update"       // 定時事件, 實體定時觸發, 參數是nil, 間隔時間為updateInterval
	EventDispose  = "dispose"      // 結束事件, 實體結束時第一個執行, 參數是nil
	EventShutdown = "shutdown"     // 關閉事件, 實體結束時第二個執行, 參數是nil, 這時連線已經中斷
	EventRecv     = nets.EventRecv // 接收訊息事件, 當接收訊息後觸發, 參數是訊息物件
	EventSend     = nets.EventSend // 傳送訊息事件, 當傳送訊息後觸發, 參數是訊息物件
)

const (
	UpdateInterval = time.Second // 定時事件間隔時間
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
