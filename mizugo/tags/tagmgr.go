package tags

import (
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
)

// NewTagmgr 建立標籤管理器
func NewTagmgr() *Tagmgr {
	return &Tagmgr{
		data: map[string]*hashset.Set{},
	}
}

// Tagmgr 標籤管理器
type Tagmgr struct {
	data map[string]*hashset.Set // 標籤列表
	lock sync.Mutex              // 執行緒鎖
}

// Add 新增標籤
func (this *Tagmgr) Add(value any, tag ...string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range tag {
		set := this.find(itor)
		set.Add(value)
	} // for
}

// Del 刪除標籤
func (this *Tagmgr) Del(value any, tag ...string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range tag {
		set := this.find(itor)
		set.Remove(value)
	} // for
}

// Get 取得物件
func (this *Tagmgr) Get(tag string) []any {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.find(tag).Values()
}

// Tag 取得標籤
func (this *Tagmgr) Tag(value any) []string {
	this.lock.Lock()
	defer this.lock.Unlock()

	result := []string{}

	for tag, set := range this.data {
		if set.Contains(value) {
			result = append(result, tag)
		} // if
	} // for

	return result
}

// find 尋找標籤列表
func (this *Tagmgr) find(tag string) *hashset.Set {
	result, ok := this.data[tag]

	if ok == false {
		result = hashset.New()
		this.data[tag] = result
	} // if

	return result
}
