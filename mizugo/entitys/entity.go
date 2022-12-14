package entitys

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/mizugo/events"
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID) *Entity {
	return &Entity{
		entityID:  entityID,
		modulemgr: NewModulemgr(),
		eventmgr:  events.NewEventmgr(eventSize),
	}
}

// Entity 實體資料
type Entity struct {
	entityID  EntityID         // 實體編號
	modulemgr *Modulemgr       // 模組管理器
	eventmgr  *events.Eventmgr // 事件管理器
	enable    atomic.Bool      // 啟用旗標
}

// EntityID 實體編號
type EntityID int64

// EntityID 取得實體編號
func (this *Entity) EntityID() EntityID {
	return this.entityID
}

// AddModule 新增模組
func (this *Entity) AddModule(module Moduler) error {
	if this.enable.Load() {
		return fmt.Errorf("entity add module: overdue")
	} // if

	if err := this.modulemgr.Add(module); err != nil {
		return fmt.Errorf("entity add module: %w", err)
	} // if

	module.Internal().entity = this
	return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) Moduler {
	return this.modulemgr.Get(moduleID)
}

// SubEvent 訂閱事件, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Entity) SubEvent(name string, process events.Process) {
	this.eventmgr.Sub(name, process)
}

// PubOnceEvent 發布單次事件
func (this *Entity) PubOnceEvent(name string, param any) {
	this.eventmgr.PubOnce(name, param)
}

// PubFixedEvent 發布定時事件, 回傳用於停止定時事件的控制物件
func (this *Entity) PubFixedEvent(name string, param any, interval time.Duration) *events.Fixed {
	return this.eventmgr.PubFixed(name, param, interval)
}

// Enable 取得啟用旗標
func (this *Entity) Enable() bool {
	return this.enable.Load()
}

// initialize 初始化處理
func (this *Entity) initialize() {
	if this.enable.CompareAndSwap(false, true) {
		this.eventmgr.Sub(eventAwake, processAwake)
		this.eventmgr.Sub(eventStart, processStart)
		this.eventmgr.Sub(eventDispose, processDispose)
		this.eventmgr.Sub(eventUpdate, processUpdate)
		this.eventmgr.Initialize()
		module := this.modulemgr.All()

		for _, itor := range module {
			this.eventmgr.PubOnce(eventAwake, itor)
		} // for

		for _, itor := range module {
			this.eventmgr.PubOnce(eventStart, itor)
		} // for

		for _, itor := range module {
			itor.Internal().update = this.eventmgr.PubFixed(eventUpdate, itor, updateInterval)
		} // for
	} // if
}

// finalize 結束處理
func (this *Entity) finalize() {
	for _, itor := range this.modulemgr.All() {
		this.eventmgr.PubOnce(eventDispose, itor)
		itor.Internal().updateStop()
	} // for

	this.enable.Store(false)
	this.eventmgr.Finalize()
}

// processAwake 處理awake事件
func processAwake(param any) {
	if module, ok := param.(Awaker); ok {
		module.Awake()
	} // if
}

// processStart 處理start事件
func processStart(param any) {
	if module, ok := param.(Starter); ok {
		module.Start()
	} // if
}

// processDispose 處理dispose事件
func processDispose(param any) {
	if module, ok := param.(Disposer); ok {
		module.Dispose()
	} // if
}

// processUpdate 處理update事件
func processUpdate(param any) {
	if module, ok := param.(Updater); ok {
		module.Update()
	} // if
}
