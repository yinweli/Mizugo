package entitys

import (
	"sync"

	"github.com/yinweli/Mizugo/core/nets"
)

// NewReactmgr 建立反應管理器
func NewReactmgr() *Reactmgr {
	return &Reactmgr{}
}

// Reactmgr 反應管理器
type Reactmgr struct {
	react nets.Reactor // 反應物件
	once  sync.Once    // 單次執行緒鎖
}

// Set 設定反應物件
func (this *Reactmgr) Set(react nets.Reactor) {
	this.once.Do(func() {
		this.react = react
	})
}

// Get 取得反應物件
func (this *Reactmgr) Get() nets.Reactor {
	return this.react
}
