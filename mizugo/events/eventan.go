package events

import (
	"sync/atomic"
	"time"
)

// NewEventan 建立事件管理器
func NewEventan(bufferSize int) *Eventan {
	return &Eventan{
		pubsub: NewPubsub(),
		event:  make(chan event, bufferSize),
	}
}

// Eventan 事件管理器
type Eventan struct {
	pubsub *Pubsub     // 訂閱/發布資料
	event  chan event  // 事件通道
	finish atomic.Bool // 結束旗標
}

// event 事件資料
type event struct {
	name  string // 事件名稱
	param any    // 事件參數
}

// Initialize 初始化處理, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Eventan) Initialize() {
	go func() {
		for this.finish.Load() == false {
			e := <-this.event
			this.pubsub.Pub(e.name, e.param)
		} // for

		// 當事件管理器要關閉時, 首先把事件通道關閉, 避免有更多的事件跑進來
		// 然後把剩餘的事件執行完畢後結束

		close(this.event)

		for e := range this.event {
			this.pubsub.Pub(e.name, e.param)
		} // for
	}()
}

// Finalize 結束處理
func (this *Eventan) Finalize() {
	if this.finish.CompareAndSwap(false, true) {
		this.finish.Store(true)
		this.event <- event{} // 新增一個空事件, 讓結束程序得以開始運作
	} // if
}

// Sub 訂閱事件, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Eventan) Sub(name string, process Process) {
	this.pubsub.Sub(name, process)
}

// PubOnce 發布單次事件
func (this *Eventan) PubOnce(name string, param any) {
	if this.finish.Load() == false {
		this.event <- event{
			name:  name,
			param: param,
		}
	} // if
}

// PubFixed 發布定時事件, 回傳用於停止定時事件的定時控制器
func (this *Eventan) PubFixed(name string, param any, interval time.Duration) *Fixed {
	fixed := &Fixed{}

	go func() {
		timeout := time.After(interval)

		for {
			select {
			case <-timeout:
				this.PubOnce(name, param)

			default:
				if fixed.State() {
					return
				} // if
			} // select
		} // for
	}()

	return fixed
}
