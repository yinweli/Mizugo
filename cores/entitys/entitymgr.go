package entitys

import (
	"fmt"
	"sort"
	"sync"
)

// NewEntitymgr 建立實體管理器
func NewEntitymgr() *Entitymgr {
	return &Entitymgr{
		data: map[EntityID]*Entity{},
	}
}

// Entitymgr 實體管理器
type Entitymgr struct {
	data map[EntityID]*Entity // 實體列表
	lock sync.RWMutex         // 執行緒鎖
}

// Add 新增實體
func (this *Entitymgr) Add(entity *Entity) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	entityID := entity.EntityID()

	if _, ok := this.data[entityID]; ok {
		return fmt.Errorf("entitymgr add: duplicate entity: %v", entityID)
	} // if

	this.data[entityID] = entity
	entity.initialize()
	return nil
}

// Del 刪除實體
func (this *Entitymgr) Del(entityID EntityID) *Entity {
	this.lock.Lock()
	defer this.lock.Unlock()

	if entity, ok := this.data[entityID]; ok {
		delete(this.data, entityID)
		entity.finalize()
		return entity
	} // if

	return nil
}

// Clear 清除實體
func (this *Entitymgr) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.data {
		itor.finalize()
	} // for

	this.data = map[EntityID]*Entity{}
}

// Get 取得實體
func (this *Entitymgr) Get(entityID EntityID) *Entity {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.data[entityID]
}

// All 取得實體列表
func (this *Entitymgr) All() []*Entity {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := []*Entity{}

	for _, itor := range this.data {
		result = append(result, itor)
	} // for

	sort.Slice(result, func(r, l int) bool {
		return result[r].EntityID() < result[l].EntityID()
	})
	return result
}

// Count 取得實體數量
func (this *Entitymgr) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return len(this.data)
}
