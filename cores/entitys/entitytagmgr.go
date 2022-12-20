package entitys

import (
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
)

// NewEntityTagmgr 建立實體標籤管理器
func NewEntityTagmgr() *EntityTagmgr {
	return &EntityTagmgr{
		data: map[string]*hashset.Set{},
	}
}

// EntityTagmgr 實體標籤管理器
type EntityTagmgr struct {
	data map[string]*hashset.Set // 標籤列表
	lock sync.RWMutex            // 執行緒鎖
}

// Add 新增標籤
func (this *EntityTagmgr) Add(entity *Entity, tag ...string) {
	if entity == nil || len(tag) == 0 {
		return
	} // if

	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range tag {
		set := this.find(itor)
		set.Add(entity)
	} // for
}

// Del 刪除標籤
func (this *EntityTagmgr) Del(entity *Entity, tag ...string) {
	if entity == nil || len(tag) == 0 {
		return
	} // if

	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range tag {
		set := this.find(itor)
		set.Remove(entity)
	} // for
}

// Get 取得實體
func (this *EntityTagmgr) Get(tag string) []*Entity {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := []*Entity{}

	for _, itor := range this.find(tag).Values() {
		result = append(result, itor.(*Entity))
	} // for

	return result
}

// Tag 取得標籤
func (this *EntityTagmgr) Tag(entity *Entity) []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := []string{}

	for tag, set := range this.data {
		if set.Contains(entity) {
			result = append(result, tag)
		} // if
	} // for

	return result
}

// find 尋找標籤列表
func (this *EntityTagmgr) find(tag string) *hashset.Set {
	result, ok := this.data[tag]

	if ok == false {
		result = hashset.New()
		this.data[tag] = result
	} // if

	return result
}
