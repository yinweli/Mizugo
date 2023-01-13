package utils

import (
	"sync"
)

// SyncOnce 單次執行器
type SyncOnce struct {
	once sync.Once // 單次執行物件
	done bool      // 執行旗標
}

// Do 單次執行
func (this *SyncOnce) Do(f func()) {
	this.once.Do(func() {
		f()
		this.done = true
	})
}

// Done 取得執行旗標
func (this *SyncOnce) Done() bool {
	return this.done
}

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
