package events

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const separateSubID = "@" // 訂閱索引分隔字串

// NewEventmgr 建立事件管理器
func NewEventmgr(capacity int) *Eventmgr {
	return &Eventmgr{
		notify: make(chan notify, capacity),
		pubsub: newPubsub(),
	}
}

// Eventmgr 事件管理器
type Eventmgr struct {
	notify chan notify // 通知通道
	pubsub *pubsub     // 訂閱/發布資料
	finish atomic.Bool // 結束旗標
}

// Process 處理函式類型
type Process func(param any)

// Do 執行處理
func (this Process) Do(param any) {
	if this != nil {
		this(param)
	} // if
}

// notify 通知資料
type notify struct {
	pub   bool   // 發布旗標, true表示為發布事件, false則為取消訂閱
	name  string // 事件名稱
	param any    // 事件參數/事件索引
}

// Initialize 初始化處理, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Eventmgr) Initialize() {
	go func() {
		for {
			select {
			case n := <-this.notify:
				if n.pub {
					this.pubsub.pub(n.name, n.param)
				} else {
					this.pubsub.unsub(n.name)
				} // if

			default:
				if this.finish.Load() {
					goto Finish
				} // if
			} // select
		} // for

	Finish:
		close(this.notify)

		for n := range this.notify { // 把剩餘的事件都做完
			if n.pub {
				this.pubsub.pub(n.name, n.param)
			} // if
		} // for
	}()
}

// Finalize 結束處理
func (this *Eventmgr) Finalize() {
	this.finish.Store(true)
}

// Sub 訂閱事件
func (this *Eventmgr) Sub(name string, process Process) string {
	return this.pubsub.sub(name, process)
}

// Unsub 取消訂閱事件
func (this *Eventmgr) Unsub(subID string) {
	if this.finish.Load() {
		return
	} // if

	this.notify <- notify{
		pub:  false,
		name: subID,
	}
}

// PubOnce 發布單次事件
func (this *Eventmgr) PubOnce(name string, param any) {
	if this.finish.Load() {
		return
	} // if

	this.notify <- notify{
		pub:   true,
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
					this.notify <- notify{
						pub:   true,
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
		data: map[string]list{},
	}
}

// pubsub 訂閱/發布資料
type pubsub struct {
	serial int64           // 事件序號
	data   map[string]list // 處理列表
	lock   sync.RWMutex    // 執行緒鎖
}

// list 處理列表
type list map[int64]Process

// sub 訂閱
func (this *pubsub) sub(name string, process Process) string {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.serial++

	if cell, ok := this.data[name]; ok {
		cell[this.serial] = process
	} else {
		this.data[name] = list{this.serial: process}
	} // if

	return subIDEncode(name, this.serial)
}

// unsub 取消訂閱
func (this *pubsub) unsub(subID string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if name, serial, ok := subIDDecode(subID); ok {
		if cell, ok := this.data[name]; ok {
			delete(cell, serial)
		} // if
	} // if
}

// pub 發布
func (this *pubsub) pub(name string, param any) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if cell, ok := this.data[name]; ok {
		for _, itor := range cell {
			itor.Do(param)
		} // for
	} // if
}

// subIDEncode 編碼訂閱索引
func subIDEncode(name string, serial int64) string {
	builder := strings.Builder{}
	builder.WriteString(name)
	builder.WriteString(separateSubID)
	builder.WriteString(strconv.FormatInt(serial, 10))
	return builder.String()
}

// subIDDecode 解碼訂閱索引
func subIDDecode(subID string) (name string, serial int64, ok bool) {
	if before, after, ok := strings.Cut(subID, separateSubID); ok {
		if serial, err := strconv.ParseInt(after, 10, 64); err == nil {
			return before, serial, true
		} // if
	} // if

	return "", 0, false
}
