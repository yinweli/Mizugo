package entitys

import (
	"fmt"
	"sync/atomic"

	"github.com/yinweli/Mizugo/mizugo/events"
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID, name string) *Entity {
	return &Entity{
		entityID: entityID,
		name:     name,
		modulean: NewModulean(),
		eventan:  events.NewEventan(processEvent),
	}
}

// Entity 實體資料
type Entity struct {
	entityID EntityID        // 實體編號
	name     string          // 實體名稱
	enable   atomic.Bool     // 啟用旗標
	modulean *Modulean       // 模組管理器
	eventan  *events.Eventan // 事件管理器
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
func (this *Entity) AddModule(module Moduler) error {
	if err := this.modulean.Add(module); err != nil {
		return fmt.Errorf("entity add module: %w", err)
	} // if

	module.Host(this)

	if this.enable.Load() {
		this.eventan.InvokeAwake(module)
		this.eventan.InvokeStart(module)
		this.eventan.InvokeUpdate(module, updateInterval)
	} // if

	return nil
}

// DelModule 刪除模組
func (this *Entity) DelModule(moduleID ModuleID) Moduler {
	if module := this.modulean.Del(moduleID); module != nil {
		this.eventan.InvokeDispose(module)
		return module
	} // if

	return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) Moduler {
	return this.modulean.Get(moduleID)
}

// initialize 初始化處理 TODO: 單元測試
func (this *Entity) initialize() {
	if this.enable.CompareAndSwap(false, true) {
		this.eventan.Initialize()
		module := this.modulean.All()

		for _, itor := range module {
			this.eventan.InvokeAwake(itor)
		} // for

		for _, itor := range module {
			this.eventan.InvokeStart(itor)
		} // for

		for _, itor := range module {
			this.eventan.InvokeUpdate(itor, updateInterval)
		} // for
	} // if
}

// finalize 結束處理 TODO: 單元測試
func (this *Entity) finalize() {
	this.enable.Store(false)
	this.eventan.Finalize()
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
