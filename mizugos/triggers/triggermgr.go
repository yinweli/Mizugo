package triggers

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
)

// NewTriggermgr 建立信號調度管理器
func NewTriggermgr() *Triggermgr {
	return &Triggermgr{
		data: map[string]*Trigger{},
		fptr: map[uintptr]bool{},
		done: make(chan bool, 1),
	}
}

// Triggermgr 信號調度管理器, 提供了新增, 鎖定和觸發信號的功能;
// 每個信號都有一個名稱和一個處理函式, 管理器會保證名稱和處理函式的唯一性;
// 當觸發信號時, 管理器會先鎖定該信號, 然後執行信號的處理函式,
// 在處理函式執行期間,其他試圖鎖定該信號的操作會被阻塞,避免出現競爭條件
type Triggermgr struct {
	data map[string]*Trigger // 信號列表
	fptr map[uintptr]bool    // 處理位址列表
	done chan bool           // 關閉用通道
	lock sync.RWMutex        // 執行緒鎖
}

// Finalize 結束處理
func (this *Triggermgr) Finalize() {
	this.done <- true
}

// Watch 設定監聽redis的頻道, 當頻道出現信號名稱時, 就觸發該信號
func (this *Triggermgr) Watch(client redis.UniversalClient, channelName string) {
	go func() {
		pubsub := client.Subscribe(ctxs.Get().Ctx(), channelName)
		defer func() {
			_ = pubsub.Close()
		}()
		channel := pubsub.Channel() // 請不要放到 for - range 之後, 會造成錯誤

		for {
			select {
			case name := <-channel:
				if name == nil { // 空訊息表示頻道已經關閉
					return
				} // if

				if trigger := this.Get(name.Payload); trigger != nil {
					trigger.Invoke()
				} // if

			case <-this.done:
				return
			} // select
		} // for
	}()
}

// Add 新增信號
func (this *Triggermgr) Add(name string, exec func()) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.data[name]; ok {
		return fmt.Errorf("triggermgr add: name duplicate")
	} // if

	fptr := reflect.ValueOf(exec).Pointer()

	if _, ok := this.fptr[fptr]; ok {
		return fmt.Errorf("triggermgr add: exec duplicate")
	} // if

	this.data[name] = &Trigger{
		name: name,
		exec: exec,
	}
	this.fptr[fptr] = true
	return nil
}

// Get 取得信號
func (this *Triggermgr) Get(name string) *Trigger {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data[name]
}

// Trigger 信號資料
type Trigger struct {
	name string       // 信號名稱
	exec func()       // 處理函式
	lock sync.RWMutex // 執行緒鎖
}

// Lock 鎖定信號
func (this *Trigger) Lock() {
	this.lock.RLock()
}

// Unlock 解鎖信號
func (this *Trigger) Unlock() {
	this.lock.RUnlock()
}

// Invoke 觸發信號
func (this *Trigger) Invoke() {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.exec != nil {
		this.exec()
	} // if
}
