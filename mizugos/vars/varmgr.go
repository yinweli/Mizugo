package vars

import (
	"sync"
)

// NewVarmgr 建立變數管理器
func NewVarmgr() *Varmgr {
	return &Varmgr{
		data: map[string]any{},
	}
}

// Varmgr 變數管理器, 讓使用者可以用字串作為索引, 存取全域資料
type Varmgr struct {
	data map[string]any // 變數列表
	lock sync.RWMutex   // 執行緒鎖
}

// Reset 重置變數管理器
func (this *Varmgr) Reset() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data = map[string]any{}
}

// Set 設定變數
func (this *Varmgr) Set(name string, data any) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data[name] = data
}

// Del 刪除變數
func (this *Varmgr) Del(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.data, name)
}

// Get 取得變數
func (this *Varmgr) Get(name string) any {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if data, ok := this.data[name]; ok {
		return data
	} // if

	return nil
}
