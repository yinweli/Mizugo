package entitys

import (
	"fmt"
	"sort"
	"sync"
)

// NewModulemap 建立模組列表
func NewModulemap() *Modulemap {
	return &Modulemap{
		data: map[ModuleID]Moduler{},
	}
}

// Modulemap 模組列表, 負責新增/刪除/取得模組等功能
type Modulemap struct {
	data map[ModuleID]Moduler // 模組列表
	lock sync.RWMutex         // 執行緒鎖
}

// Add 新增模組
func (this *Modulemap) Add(module Moduler) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	moduleID := module.ModuleID()

	if _, ok := this.data[moduleID]; ok {
		return fmt.Errorf("modulemap add: duplicate module: %v", moduleID)
	} // if

	this.data[moduleID] = module
	return nil
}

// Del 刪除模組
func (this *Modulemap) Del(moduleID ModuleID) Moduler {
	this.lock.Lock()
	defer this.lock.Unlock()

	if module, ok := this.data[moduleID]; ok {
		delete(this.data, moduleID)
		return module
	} // if

	return nil
}

// Get 取得模組
func (this *Modulemap) Get(moduleID ModuleID) Moduler {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data[moduleID]
}

// All 取得模組列表
func (this *Modulemap) All() []Moduler {
	this.lock.RLock()
	defer this.lock.RUnlock()
	result := []Moduler{}

	for _, itor := range this.data {
		result = append(result, itor)
	} // for

	sort.Slice(result, func(r, l int) bool {
		return result[r].ModuleID() < result[l].ModuleID()
	})
	return result
}

// Count 取得模組數量
func (this *Modulemap) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.data)
}
