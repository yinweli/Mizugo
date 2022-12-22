package labels

import (
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
)

// NewLabelobj 建立標籤物件
func NewLabelobj() *Labelobj {
	return &Labelobj{
		data: hashset.New(),
	}
}

// Labelobj 標籤物件
type Labelobj struct {
	data *hashset.Set // 標籤列表
	lock sync.RWMutex // 執行緒鎖
}

// add 新增標籤
func (this *Labelobj) add(label ...string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range label {
		this.data.Add(itor)
	} // for
}

// del 刪除標籤
func (this *Labelobj) del(label ...string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range label {
		this.data.Remove(itor)
	} // for
}

// erase 清除標籤
func (this *Labelobj) erase() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data = hashset.New()
}

// Label 取得標籤
func (this *Labelobj) Label() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := []string{}

	for _, label := range this.data.Values() {
		result = append(result, label.(string))
	} // for

	return result
}
