package entitys

import (
	"sync"

	"github.com/yinweli/Mizugo/cores/nets"
)

// SessionAttr 會話會話器
type SessionAttr struct {
	session nets.Sessioner // 會話物件
	lock    sync.RWMutex   // 執行緒鎖
}

// Set 設定會話物件
func (this *SessionAttr) Set(session nets.Sessioner) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session = session
}

// Get 取得會話物件
func (this *SessionAttr) Get() nets.Sessioner {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.session
}
