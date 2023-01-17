package procs

import (
	"sync"
)

// 管理器, 負責管理訊息處理函式

// newProcmgr 建立管理器
func newProcmgr() *procmgr {
	return &procmgr{
		data: map[MessageID]Process{},
	}
}

// procmgr 管理器
type procmgr struct {
	data map[MessageID]Process // 處理列表
	lock sync.RWMutex          // 執行緒鎖
}

// Add 新增處理函式
func (this *procmgr) Add(messageID MessageID, process Process) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data[messageID] = process
}

// Del 刪除處理函式
func (this *procmgr) Del(messageID MessageID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.data, messageID)
}

// Get 取得處理函式
func (this *procmgr) Get(messageID MessageID) Process {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.data[messageID]
}
