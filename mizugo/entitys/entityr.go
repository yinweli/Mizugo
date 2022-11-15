package entitys

import (
	"fmt"
	"sort"
	"sync"
)

// NewEntityr 建立實體管理器
func NewEntityr() *Entityr {
	return &Entityr{
		data: map[EntityID]*Entity{},
	}
}

// Entityr 實體管理器
type Entityr struct {
	data map[EntityID]*Entity // 實體列表
	lock sync.Mutex           // 執行緒鎖
}

// Add 新增實體
func (this *Entityr) Add(entity *Entity) error {
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
func (this *Entityr) Del(entityID EntityID) *Entity {
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
func (this *Entityr) Get(entityID EntityID) *Entity {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.data[entityID]
}

// All 取得實體列表
func (this *Entityr) All() []*Entity {
	result := []*Entity{}

	for _, itor := range this.data {
		result = append(result, itor)
	} // for

	sort.Slice(result, func(r, l int) bool {
		return result[r].EntityID() < result[l].EntityID()
	})
	return result
}