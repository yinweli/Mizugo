package entitys

import (
	"fmt"
	"sort"
	"sync"
)

// NewModulean 建立模組管理器
func NewModulean() *Modulean {
	return &Modulean{
		data: map[ModuleID]Moduler{},
	}
}

// Modulean 模組管理器
type Modulean struct {
	data map[ModuleID]Moduler // 模組列表
	lock sync.Mutex           // 執行緒鎖
}

// Add 新增模組
func (this *Modulean) Add(module Moduler) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	moduleID := module.ModuleID()

	if _, ok := this.data[moduleID]; ok {
		return fmt.Errorf("modulean add: duplicate moduleID")
	} // if

	this.data[moduleID] = module
	return nil
}

// Del 刪除模組
func (this *Modulean) Del(moduleID ModuleID) Moduler {
	this.lock.Lock()
	defer this.lock.Unlock()

	if module, ok := this.data[moduleID]; ok {
		delete(this.data, moduleID)
		return module
	} // if

	return nil
}

// Get 取得模組
func (this *Modulean) Get(moduleID ModuleID) Moduler {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.data[moduleID]
}

// All 取得模組列表
func (this *Modulean) All() []Moduler {
	result := []Moduler{}

	for _, itor := range this.data {
		result = append(result, itor)
	} // for

	sort.Slice(result, func(r, l int) bool {
		return result[r].ModuleID() < result[l].ModuleID()
	})
	return result
}
