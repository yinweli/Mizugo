package entitys

import (
    "fmt"
    `sync/atomic`

    `github.com/yinweli/Mizugo/mizugo/event`
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID, name string) *Entity {
    return &Entity{
        entityID: entityID,
        name:     name,
        moduler:  NewModuler(),
        event:    event.NewEvent(eventBufferSize),
    }
}

// Entity 實體資料
type Entity struct {
    entityID EntityID     // 實體編號
    name     string       // 實體名稱
    startup  atomic.Bool  // 啟動旗標
    moduler  *Moduler     // 模組管理器
    event    *event.Event // 事件管理器
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

    if this.startup.Load() {
        this.event.Add(&eventAwake{
            module: module,
        })
        this.event.Add(&eventStart{
            module: module,
        })
    } // if

    return nil
}

// DelModule 刪除模組
func (this *Entity) DelModule(moduleID ModuleID) IModule {
    if module := this.moduler.Del(moduleID); module != nil {
        this.event.Add(&eventDispose{
            module: module,
        })
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

// begin 啟動實體
func (this *Entity) begin() {
    if this.startup.CompareAndSwap(false, true) {
        this.event.Begin(this.proc)
        module := this.moduler.All()

        for _, itor := range module {
            this.event.Add(&eventAwake{
                module: itor,
            })
        } // for

        for _, itor := range module {
            this.event.Add(&eventStart{
                module: itor,
            })
        } // for
    } // if
}

// end 結束實體
func (this *Entity) end() {
    this.startup.Store(false)
    this.event.End()
}

// proc 事件處理
func (this *Entity) proc(event any) {

}
