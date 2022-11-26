package entitys

import (
	"fmt"
	"sort"
	"sync"
)

// NewModulemgr 建立模組管理器
func NewModulemgr() *Modulemgr {
	return &Modulemgr{
		data: map[ModuleID]Moduler{},
	}
}

// Modulemgr 模組管理器
type Modulemgr struct {
	data map[ModuleID]Moduler // 模組列表
	lock sync.Mutex           // 執行緒鎖
}

// Add 新增模組
func (this *Modulemgr) Add(module Moduler) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	moduleID := module.ModuleID()

	if _, ok := this.data[moduleID]; ok {
		return fmt.Errorf("modulemgr add: duplicate moduleID")
	} // if

	this.data[moduleID] = module
	return nil
}

// Del 刪除模組
func (this *Modulemgr) Del(moduleID ModuleID) Moduler {
	this.lock.Lock()
	defer this.lock.Unlock()

	if module, ok := this.data[moduleID]; ok {
		delete(this.data, moduleID)
		return module
	} // if

	return nil
}

// Get 取得模組
func (this *Modulemgr) Get(moduleID ModuleID) Moduler {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.data[moduleID]
}

// All 取得模組列表
func (this *Modulemgr) All() []Moduler {
	this.lock.Lock()
	defer this.lock.Unlock()
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
func (this *Modulemgr) Count() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return len(this.data)
}
