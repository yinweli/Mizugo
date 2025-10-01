package helps

import (
	"sync"
	"sync/atomic"
)

// SyncOnce 單次執行器, 這是執行緒安全的執行, 用於需要保證某個動作只會執行一次的情境, 例如初始化設定或建立全域資源
//
// 範例:
//
//	once := helps.SyncOnce{}
//
//	// 初始化邏輯, 只會被執行一次
//	once.Do(func() {
//	    fmt.Println("初始化設定")
//	})
//
//	// 檢查是否已被執行
//	if once.Done() {
//	    fmt.Println("已初始化完成")
//	} // if
type SyncOnce struct {
	once sync.Once   // 單次執行物件
	done atomic.Bool // 執行旗標
}

// Do 執行傳入的函式 f, 但僅會執行一次
//
// 多個 goroutine 同時呼叫時, 也只會有一個 goroutine 執行, 之後的呼叫都不會再執行
//
// 回傳 true 表示本次呼叫真的執行了 f, false 表示 f 已經執行過
//
// 注意: 若 f 為 nil, 仍會被視為已執行過
func (this *SyncOnce) Do(f func()) (do bool) {
	this.once.Do(func() {
		do = true
		this.done.Store(true)

		if f != nil {
			f()
		} // if
	})

	return do
}

// Done 回傳是否已經執行過
//
// 無論 f 是否為 nil, 只要 Do 曾被呼叫過, 就會回傳 true
func (this *SyncOnce) Done() bool {
	return this.done.Load()
}

// SyncAttr 屬性存取器, 這是執行緒安全的屬性存取器, 適合用於多個 goroutine 同時讀寫同一個變數時, 確保資料一致性
//
// 範例:
//
//	count := helps.SyncAttr[int]{}
//
//	// 設定屬性
//	count.Set(10)
//
//	// 取得屬性
//	value := count.Get()
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

// Get 取得屬性物件, 回傳屬性的複本
func (this *SyncAttr[T]) Get() T {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.attr
}
