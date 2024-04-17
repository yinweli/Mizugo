package labels

import (
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
)

// NewLabelmgr 建立標籤管理器
func NewLabelmgr() *Labelmgr {
	return &Labelmgr{
		data: map[string]*hashset.Set{},
	}
}

// Labelmgr 標籤管理器, 讓使用者可以在可管理物件上新增/刪除標籤, 也可以通過標籤來取得管理中的物件;
// 建立可管理物件需要定義標籤結構, 並且把 Label 作為結構的第一個成員;
// 標籤必須是字串, 標籤可以關連到多個物件, 而物件也可以擁有多個標籤
type Labelmgr struct {
	data map[string]*hashset.Set // 標籤列表
	lock sync.RWMutex            // 執行緒鎖
}

// Labeler 標籤介面
type Labeler interface {
	// add 新增標籤
	add(label ...string)

	// del 刪除標籤
	del(label ...string)

	// erase 清除標籤
	erase()

	// label 取得標籤
	label() []string
}

// Add 新增標籤
func (this *Labelmgr) Add(obj any, label ...string) {
	if labelObj, ok := obj.(Labeler); ok && len(label) != 0 {
		labelObj.add(label...)

		this.lock.Lock()
		defer this.lock.Unlock()

		for _, itor := range label {
			set := this.find(itor)
			set.Add(obj)
		} // for
	} // if
}

// Del 刪除標籤
func (this *Labelmgr) Del(obj any, label ...string) {
	if labelObj, ok := obj.(Labeler); ok && len(label) != 0 {
		labelObj.del(label...)

		this.lock.Lock()
		defer this.lock.Unlock()

		for _, itor := range label {
			set := this.find(itor)
			set.Remove(obj)
		} // for
	} // if
}

// Erase 清除標籤
func (this *Labelmgr) Erase(obj any) {
	if labelObj, ok := obj.(Labeler); ok {
		this.lock.Lock()
		defer this.lock.Unlock()

		for _, itor := range labelObj.label() {
			set := this.find(itor)
			set.Remove(obj)
		} // for

		labelObj.erase()
	} // if
}

// Get 取得物件
func (this *Labelmgr) Get(label string) []any {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.find(label).Values()
}

// Label 取得標籤
func (this *Labelmgr) Label(obj any) []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if labelObj, ok := obj.(Labeler); ok {
		return labelObj.label()
	} // if

	return []string{}
}

// find 尋找標籤列表
func (this *Labelmgr) find(label string) *hashset.Set {
	result, ok := this.data[label]

	if ok == false {
		result = hashset.New()
		this.data[label] = result
	} // if

	return result
}
