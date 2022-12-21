package msgs

import (
	"sync"
)

// NewMsgmgr 建立訊息管理器
func NewMsgmgr() *Msgmgr {
	return &Msgmgr{
		data: map[MessageID]Process{},
	}
}

// Msgmgr 訊息管理器
type Msgmgr struct {
	data map[MessageID]Process // 訊息列表
	lock sync.RWMutex          // 執行緒鎖
}

// Add 新增訊息處理函式
func (this *Msgmgr) Add(messageID MessageID, process Process) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data[messageID] = process
}

// Del 刪除訊息處理函式
func (this *Msgmgr) Del(messageID MessageID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.data, messageID)
}

// Get 取得訊息處理函式
func (this *Msgmgr) Get(messageID MessageID) Process {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.data[messageID]
}
