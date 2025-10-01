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

// Eventmap 事件列表
//
// 提供事件的訂閱(Sub), 取消訂閱(Unsub), 與發布(PubOnce / PubDelay / PubFixed)等功能
//
// 併發/有序性:
//   - Eventmap 的公開方法具備執行緒安全性
//   - 事件的實際調用會在「單一事件執行緒」中依序處理, 確保同名事件的執行順序
//   - 訂閱(Sub)立即生效; 取消訂閱(Unsub)與發布(PubXXX)會送入「事件執行緒」序列化處理
//
// 生命週期:
//   - Initialize(): 啟動事件執行緒, 僅能呼叫一次, 重複呼叫會回傳錯誤
//   - Finalize(): 關閉事件執行緒, 會嘗試處理通道中已存在的剩餘事件, 新的發布/取消訂閱請求會被忽略
type Eventmap struct {
	ctx    context.Context    // ctx物件
	cancel context.CancelFunc // 取消物件
	notify chan notify        // 通知通道
	pubsub *pubsub            // 訂閱/發布資料
	once   helps.SyncOnce     // 單次執行物件
	close  atomic.Bool        // 關閉旗標
}

// notify 事件通知
type notify struct {
	pub   bool   // 發布旗標, true表示為發布事件, false則為取消訂閱
	name  string // 事件名稱
	param any    // 事件參數/事件索引
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

// Sub 訂閱事件, 並指定處理函式; 回傳訂閱 ID, 可用於後續 Unsub
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

// PubOnce 發布單次事件, 該事件會被送入事件執行緒並依序處理
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

// PubDelay 發布延遲事件, 該事件會延遲一段時間才發布, 僅發布一次, 不會重複
func (this *Eventmap) PubDelay(name string, param any, delay time.Duration) {
	if this.close.Load() {
		return
	} // if

	pools.DefaultPool.Submit(func() {
		select {
		case <-time.After(delay):
			if this.close.Load() == false {
				this.notify <- notify{
					pub:   true,
					name:  name,
					param: param,
				}
			}
		case <-this.ctx.Done():
			return
		} // switch
	})
}

// PubFixed 發布定時事件, 該事件會週期性地發布, 由於不能刪除定時事件, 發布前請先評估需求
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
//
// 實作處理函式時務必避免長時間阻塞，否則會影響同名事件的整體處理延遲
type Process func(param any)

// Do 執行處理
func (this Process) Do(param any) {
	if this != nil {
		this(param)
	} // if
}

// newPubsub 建立訂閱/發布資料
func newPubsub() *pubsub {
	return &pubsub{
		data: map[string]list{},
	}
}

// pubsub 訂閱/發布資料
//
// 管理事件名稱到處理函式的映射, 並以遞增序號為訂閱者分配唯一索引
type pubsub struct {
	serial int64           // 事件序號
	data   map[string]list // 處理列表
	lock   sync.RWMutex    // 執行緒鎖
}

// list 處理列表, 代表一個事件名稱底下的處理函式集合, 以遞增序號作為訂閱者索引
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

			if len(cell) == 0 {
				delete(this.data, name)
			} // if
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
