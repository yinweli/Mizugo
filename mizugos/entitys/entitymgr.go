package entitys

import (
	"cmp"
	"maps"
	"slices"
	"sync"
)

// NewEntitymgr 建立實體管理器
func NewEntitymgr() *Entitymgr {
	return &Entitymgr{
		data: map[EntityID]*Entity{},
	}
}

// Entitymgr 實體管理器, 負責新增/刪除/取得實體等功能
type Entitymgr struct {
	entityID EntityID             // 實體編號
	data     map[EntityID]*Entity // 實體列表
	lock     sync.RWMutex         // 執行緒鎖
}

// Add 新增實體
func (this *Entitymgr) Add() *Entity {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.entityID++
	entity := NewEntity(this.entityID)
	this.data[this.entityID] = entity
	return entity
}

// Del 刪除實體
func (this *Entitymgr) Del(entityID EntityID) *Entity {
	this.lock.Lock()
	defer this.lock.Unlock()

	if entity, ok := this.data[entityID]; ok {
		delete(this.data, entityID)
		return entity
	} // if

	return nil
}

// Clear 清除實體
func (this *Entitymgr) Clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
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
	result := slices.Collect(maps.Values(this.data))
	this.lock.RUnlock()
	slices.SortFunc(result, func(l, r *Entity) int {
		return cmp.Compare(l.EntityID(), r.EntityID())
	})
	return result
}

// Count 取得實體數量
func (this *Entitymgr) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.data)
}
