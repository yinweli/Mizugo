package events

import (
	"sync"
)

// NewPubsub 建立訂閱/發布資料
func NewPubsub() *Pubsub {
	return &Pubsub{
		data: map[string][]Process{},
	}
}

// Pubsub 訂閱/發布資料
type Pubsub struct {
	data map[string][]Process // 處理列表
	lock sync.Mutex           // 執行緒鎖
}

// Process 處理函式類型
type Process func(param any)

// Sub 訂閱
func (this *Pubsub) Sub(name string, process Process) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if event, ok := this.data[name]; ok {
		this.data[name] = append(event, process)
	} else {
		this.data[name] = []Process{process}
	} // if
}

// Pub 發布
func (this *Pubsub) Pub(name string, param any) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if event, ok := this.data[name]; ok {
		for _, itor := range event {
			itor(param)
		} // for
	} // if
}
