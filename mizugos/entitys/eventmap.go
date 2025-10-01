package entitys

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/pools"
)

// NewEventmap 建立事件列表
func NewEventmap() *Eventmap {
	ctx, cancel := context.WithCancel(context.Background())
	return &Eventmap{
		ctx:    ctx,
		cancel: cancel,
		notify: make(chan notify, EventCapacity),
		pubsub: newPubsub(),
	}
}

// Eventmap 事件列表, 提供了事件訂閱/發布等功能, 事件列表本身是執行緒安全的,
// 但是發布的事件會在單一的事件執行緒中被執行, 以此保證了事件有序執行
type Eventmap struct {
	ctx    context.Context    // ctx物件
	cancel context.CancelFunc // 取消物件
	notify chan notify        // 通知通道
	pubsub *pubsub            // 訂閱/發布資料
	once   helps.SyncOnce     // 單次執行物件
	close  atomic.Bool        // 關閉旗標
}

// Initialize 初始化處理
func (this *Eventmap) Initialize() error {
	if this.once.Done() {
		return fmt.Errorf("eventmap initialize: already initialize")
	} // if

	this.once.Do(func() {
		pools.DefaultPool.Submit(func() {
			for {
				select {
				case n := <-this.notify:
					if n.pub {
						this.pubsub.pub(n.name, n.param)
					} else {
						this.pubsub.unsub(n.name)
					} // if

				case <-this.ctx.Done():
					goto Finish
				} // select
			} // for

		Finish:
			close(this.notify)

			for n := range this.notify { // 把剩餘的事件都做完
				if n.pub {
					this.pubsub.pub(n.name, n.param)
				} // if
			} // for
		})
	})

	return nil
}

// Finalize 結束處理
func (this *Eventmap) Finalize() {
	if this.once.Done() == false {
		return
	} // if

	this.close.Store(true)
	this.cancel()
}

// Sub 訂閱事件
func (this *Eventmap) Sub(name string, process Process) string {
	return this.pubsub.sub(name, process)
}

// Unsub 取消訂閱事件
func (this *Eventmap) Unsub(subID string) {
	if this.close.Load() {
		return
	} // if

	this.notify <- notify{
		pub:  false,
		name: subID,
	}
}

// PubOnce 發布單次事件
func (this *Eventmap) PubOnce(name string, param any) {
	if this.close.Load() {
		return
	} // if

	this.notify <- notify{
		pub:   true,
		name:  name,
		param: param,
	}
}

// PubDelay 發布延遲事件, 事件會延遲一段時間才發布, 但仍是單次事件
func (this *Eventmap) PubDelay(name string, param any, delay time.Duration) {
	if this.close.Load() {
		return
	} // if

	pools.DefaultPool.Submit(func() {
		timeout := time.NewTimer(delay)

		for {
			select {
			case <-timeout.C:
				if this.close.Load() == false {
					this.notify <- notify{
						pub:   true,
						name:  name,
						param: param,
					}
				} // if

			case <-this.ctx.Done():
				timeout.Stop()
				return
			} // select
		} // for
	})
}

// PubFixed 發布定時事件, 請注意! 由於不能刪除定時事件, 因此發布定時事件前請多想想
func (this *Eventmap) PubFixed(name string, param any, interval time.Duration) {
	if this.close.Load() {
		return
	} // if

	pools.DefaultPool.Submit(func() {
		timeout := time.NewTicker(interval)

		for {
			select {
			case <-timeout.C:
				if this.close.Load() == false {
					this.notify <- notify{
						pub:   true,
						name:  name,
						param: param,
					}
				} // if

			case <-this.ctx.Done():
				timeout.Stop()
				return
			} // select
		} // for
	})
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

	process := []Process{}

	if cell, ok := this.data[name]; ok {
		for _, itor := range cell {
			process = append(process, itor)
		} // for
	} // if

	this.lock.RUnlock()

	for _, itor := range process {
		itor.Do(param)
	} // for
}

// subIDEncode 編碼訂閱索引
func subIDEncode(name string, serial int64) string {
	builder := strings.Builder{}
	builder.WriteString(name)
	builder.WriteString("@")
	builder.WriteString(strconv.FormatInt(serial, 10))
	return builder.String()
}

// subIDDecode 解碼訂閱索引
func subIDDecode(subID string) (name string, serial int64, ok bool) {
	if before, after, ok := strings.Cut(subID, "@"); ok {
		if serial, err := strconv.ParseInt(after, 10, 64); err == nil {
			return before, serial, true
		} // if
	} // if

	return "", 0, false
}
