package utils

import (
	"sync"
)

// SyncAttr 同步屬性器
type SyncAttr[T any] struct {
	attr T            // 屬性物件
	lock sync.RWMutex // 執行緒鎖
}

// Set 設定屬性物件
func (this *SyncAttr[T]) Set(attr T) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.attr = attr
}

// Get 取得屬性物件
func (this *SyncAttr[T]) Get() T {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.attr
}
