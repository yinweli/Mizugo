package entitys

import (
	"sync"

	"github.com/yinweli/Mizugo/cores/nets"
)

// NewSessionmgr 建立會話管理器
func NewSessionmgr() *Sessionmgr {
	return &Sessionmgr{}
}

// Sessionmgr 會話管理器
type Sessionmgr struct {
	session nets.Sessioner // 會話物件
	once    sync.Once      // 單次執行緒鎖
}

// Set 設定會話物件
func (this *Sessionmgr) Set(session nets.Sessioner) {
	this.once.Do(func() {
		this.session = session
	})
}

// Get 取得會話物件
func (this *Sessionmgr) Get() nets.Sessioner {
	return this.session
}
