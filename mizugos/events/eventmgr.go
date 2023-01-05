package events

import (
	"sync"
	"sync/atomic"
	"time"
)

// NewEventmgr 建立事件管理器
func NewEventmgr(capacity int) *Eventmgr {
	return &Eventmgr{
		event:  make(chan event, capacity),
		pubsub: newPubsub(),
	}
}

// Eventmgr 事件管理器
type Eventmgr struct {
	event  chan event  // 事件通道
	pubsub *pubsub     // 訂閱/發布資料
	finish atomic.Bool // 結束旗標
}

// Process 處理函式類型
type Process func(param any)

// Initialize 初始化處理, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Eventmgr) Initialize() {
	go func() {
		for {
			select {
			case e := <-this.event:
				this.pubsub.pub(e.name, e.param)

			default:
				if this.finish.Load() {
					goto Finish
				} // if
			} // select
		} // for

	Finish:
		close(this.event)

		for e := range this.event { // 把剩餘的事件都做完
			this.pubsub.pub(e.name, e.param)
		} // for
	}()
}

// Finalize 結束處理
func (this *Eventmgr) Finalize() {
	this.finish.Store(true)
}

// Sub 訂閱事件, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Eventmgr) Sub(name string, process Process) {
	this.pubsub.sub(name, process)
}

// PubOnce 發布單次事件
func (this *Eventmgr) PubOnce(name string, param any) {
	if this.finish.Load() {
		return
	} // if

	this.event <- event{
		name:  name,
		param: param,
	}
}

// PubFixed 發布定時事件
func (this *Eventmgr) PubFixed(name string, param any, interval time.Duration) {
	if this.finish.Load() {
		return
	} // if

	go func() {
		timeout := time.NewTicker(interval)

		for {
			select {
			case <-timeout.C:
				if this.finish.Load() == false {
					this.event <- event{
						name:  name,
						param: param,
					}
				} // if

			default:
				if this.finish.Load() {
					timeout.Stop()
					return
				} // if
			} // select
		} // for
	}()
}

// newPubsub 建立訂閱/發布資料
func newPubsub() *pubsub {
	return &pubsub{
		data: map[string]Process{},
	}
}

// pubsub 訂閱/發布資料
type pubsub struct {
	data map[string]Process // 處理列表
	lock sync.RWMutex       // 執行緒鎖
}

// sub 訂閱
func (this *pubsub) sub(name string, process Process) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data[name] = process
}

// pub 發布
func (this *pubsub) pub(name string, param any) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if process, ok := this.data[name]; ok {
		process(param)
	} // if
}

// event 事件資料
type event struct {
	name  string // 事件名稱
	param any    // 事件參數
}
