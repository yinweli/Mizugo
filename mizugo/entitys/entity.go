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
        event:    events.NewEvent(eventSize),
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
        this.event.Execute(events.Awake, module)
        this.event.Execute(events.Start, module)
    } // if

    return nil
}

// DelModule 刪除模組
func (this *Entity) DelModule(moduleID ModuleID) IModule {
    if module := this.moduler.Del(moduleID); module != nil {
        this.event.Execute(events.Dispose, module)
        return module
    } // if

    return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) IModule {
    return this.moduler.Get(moduleID)
}

// TODO: 通知封包事件
// TODO: 註冊封包處理
// TODO: 通知訊息事件
// TODO: 註冊訊息處理

// initialize 初始化處理 TODO: 單元測試
func (this *Entity) initialize() {
    if this.enable.CompareAndSwap(false, true) {
        this.event.Initialize(eventInterval, this.processEvent)
        module := this.moduler.All()

        for _, itor := range module {
            this.event.Execute(events.Awake, itor)
        } // for

        for _, itor := range module {
            this.event.Execute(events.Start, itor)
        } // for
    } // if
}

// finalize 結束處理 TODO: 單元測試
func (this *Entity) finalize() {
    this.enable.Store(false)
    this.event.Finalize()
}

// processEvent 事件處理 TODO: 單元測試
func (this *Entity) processEvent(data events.Data) {
    // TODO: 事件處理
}
