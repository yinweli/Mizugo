package entitys

import (
    "fmt"
    `sort`
    "sync"
)

// NewModuler 建立模組管理器
func NewModuler() *Moduler {
    return &Moduler{
        data: map[ModuleID]IModule{},
    }
}

// Moduler 模組管理器
type Moduler struct {
    data map[ModuleID]IModule // 模組列表
    lock sync.Mutex           // 執行緒鎖
}

// IModule 模組介面
type IModule interface {
    ModuleID() ModuleID  // 取得模組編號
    Name() string        // 取得模組名稱
    Entity() *Entity     // 取得實體物件
    Host(entity *Entity) // 設定宿主實體
}

// Add 新增模組
func (this *Moduler) Add(module IModule) error {
    this.lock.Lock()
    defer this.lock.Unlock()
    moduleID := module.ModuleID()

    if _, ok := this.data[moduleID]; ok {
        return fmt.Errorf("moduler add: duplicate moduleID")
    } // if

    this.data[moduleID] = module
    return nil
}

// Del 刪除模組
func (this *Moduler) Del(moduleID ModuleID) IModule {
    this.lock.Lock()
    defer this.lock.Unlock()

    if module, ok := this.data[moduleID]; ok {
        delete(this.data, moduleID)
        return module
    } // if

    return nil
}

// Get 取得模組
func (this *Moduler) Get(moduleID ModuleID) IModule {
    this.lock.Lock()
    defer this.lock.Unlock()

    return this.data[moduleID]
}

// All 取得模組列表
func (this *Moduler) All() []IModule {
    result := []IModule{}

    for _, itor := range this.data {
        result = append(result, itor)
    } // for

    sort.Slice(result, func(r, l int) bool {
        return result[r].ModuleID() < result[l].ModuleID()
    })
    return result
}