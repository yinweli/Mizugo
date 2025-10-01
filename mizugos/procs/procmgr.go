package procs

import (
	"sync"
)

// NewProcmgr 建立訊息處理管理器
func NewProcmgr() *Procmgr {
	return &Procmgr{
		data: map[int32]Process{},
	}
}

// Procmgr 訊息處理管理器
//
// 職責:
//   - 儲存 messageID 與對應的 Process 處理函式
//   - 在訊息解碼完成後, 由 Processor 呼叫 Procmgr 來取得對應的 Process, 並執行實際的業務邏輯
type Procmgr struct {
	data map[int32]Process // 處理列表
	lock sync.RWMutex      // 執行緒鎖
}

// Add 新增處理函式, 如果指定的 messageID 已存在, 新的 process 會覆蓋舊的 process
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
