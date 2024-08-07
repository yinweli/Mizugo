package procs

import (
	"sync"
)

// NewProcmgr 建立管理器
func NewProcmgr() *Procmgr {
	return &Procmgr{
		data: map[int32]Process{},
	}
}

// Procmgr 管理器, 負責管理訊息處理函式
type Procmgr struct {
	data map[int32]Process // 處理列表
	lock sync.RWMutex      // 執行緒鎖
}

// Add 新增處理函式
func (this *Procmgr) Add(messageID int32, process Process) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.data[messageID] = process
}

// Del 刪除處理函式
func (this *Procmgr) Del(messageID int32) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.data, messageID)
}

// Get 取得處理函式
func (this *Procmgr) Get(messageID int32) Process {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data[messageID]
}
