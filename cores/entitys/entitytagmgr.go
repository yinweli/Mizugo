package entitys

import (
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
)

// TODO: 考慮是不是(再度)把標籤變成獨立功能
// TODO: 並且利用介面來跟entity分離
// TODO: 名稱是否可以改成 label LabelMgr / LabelObj

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

	entity.tag.Add(tag...)

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

	entity.tag.Del(tag...)

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
	return entity.Tag()
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

// newEntityTag 建立實體標籤資料
func newEntityTag() *entityTag {
	return &entityTag{
		data: hashset.New(),
	}
}

// entityTag 實體標籤資料
type entityTag struct {
	data *hashset.Set // 標籤列表
	lock sync.RWMutex // 執行緒鎖
}

// Add 新增標籤
func (this *entityTag) Add(tag ...string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range tag {
		this.data.Add(itor)
	} // for
}

// Del 刪除標籤
func (this *entityTag) Del(tag ...string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range tag {
		this.data.Remove(itor)
	} // for
}

// Tag 取得標籤
func (this *entityTag) Tag() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := []string{}

	for _, tag := range this.data.Values() {
		result = append(result, tag.(string))
	} // for

	return result
}
