package entitys

import (
    "fmt"
    `sync/atomic`

    `github.com/yinweli/Mizugo/mizugo/events`
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID, name string) *Entity {
    return &Entity{
        entityID: entityID,
        name:     name,
        moduler:  NewModuler(),
        event:    events.NewEvent(processEvent),
    }
}

// Entity 實體資料
type Entity struct {
    entityID EntityID      // 實體編號
    name     string        // 實體名稱
    enable   atomic.Bool   // 啟用旗標
    moduler  *Moduler      // 模組管理器
    event    *events.Event // 事件管理器
}

// EntityID 實體編號
type EntityID int64

// EntityID 取得實體編號
func (this *Entity) EntityID() EntityID {
    return this.entityID
}

// Name 取得實體名稱
func (this *Entity) Name() string {
    return this.name
}

// AddModule 新增模組
func (this *Entity) AddModule(module IModule) error {
    if err := this.moduler.Add(module); err != nil {
        return fmt.Errorf("entity add module: %w", err)
    } // if

    module.Host(this)

    if this.enable.Load() {
        this.event.InvokeAwake(module)
        this.event.InvokeStart(module)
        this.event.InvokeUpdate(module, updateInterval)
    } // if

    return nil
}

// DelModule 刪除模組
func (this *Entity) DelModule(moduleID ModuleID) IModule {
    if module := this.moduler.Del(moduleID); module != nil {
        this.event.InvokeDispose(module)
        return module
    } // if

    return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) IModule {
    return this.moduler.Get(moduleID)
}

// initialize 初始化處理 TODO: 單元測試
func (this *Entity) initialize() {
    if this.enable.CompareAndSwap(false, true) {
        this.event.Initialize()
        module := this.moduler.All()

        for _, itor := range module {
            this.event.InvokeAwake(itor)
        } // for

        for _, itor := range module {
            this.event.InvokeStart(itor)
        } // for

        for _, itor := range module {
            this.event.InvokeUpdate(itor, updateInterval)
        } // for
    } // if
}

// finalize 結束處理 TODO: 單元測試
func (this *Entity) finalize() {
    this.enable.Store(false)
    this.event.Finalize()
}

// processEvent 事件處理 TODO: 單元測試
func processEvent(event any) {
    if e, ok := event.(*events.Awake); ok {
        if module, ok := e.Param.(ModuleAwake); ok {
            module.Awake()
        } // if
    } // if

    if e, ok := event.(*events.Start); ok {
        if module, ok := e.Param.(ModuleStart); ok {
            module.Start()
        } // if
    } // if

    if e, ok := event.(*events.Dispose); ok {
        if module, ok := e.Param.(ModuleDispose); ok {
            module.Dispose()
        } // if
    } // if

    if e, ok := event.(*events.Update); ok {
        if module, ok := e.Param.(ModuleUpdate); ok {
            module.Update()
        } // if
    } // if
}