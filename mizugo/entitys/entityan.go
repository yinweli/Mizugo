package entitys

import (
	"fmt"
	"sort"
	"sync"
)

// NewEntityan 建立實體管理器
func NewEntityan() *Entityan {
	return &Entityan{
		data: map[EntityID]*Entity{},
	}
}

// Entityan 實體管理器
type Entityan struct {
	data map[EntityID]*Entity // 實體列表
	lock sync.Mutex           // 執行緒鎖
}

// Clear 清除實體
func (this *Entityan) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.data {
		itor.finalize()
	} // for

	this.data = map[EntityID]*Entity{}
}

// Add 新增實體
func (this *Entityan) Add(entity *Entity) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	entityID := entity.EntityID()

	if _, ok := this.data[entityID]; ok {
		return fmt.Errorf("entityr add: duplicate entityID")
	} // if

	this.data[entityID] = entity
	entity.initialize()
	return nil
}

// Del 刪除實體
func (this *Entityan) Del(entityID EntityID) *Entity {
	this.lock.Lock()
	defer this.lock.Unlock()

	if entity, ok := this.data[entityID]; ok {
		delete(this.data, entityID)
		entity.finalize()
		return entity
	} // if

	return nil
}

// Get 取得實體
func (this *Entityan) Get(entityID EntityID) *Entity {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.data[entityID]
}

// All 取得實體列表
func (this *Entityan) All() []*Entity {
	this.lock.Lock()
	defer this.lock.Unlock()
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
func (this *Entityan) Count() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return len(this.data)
}
