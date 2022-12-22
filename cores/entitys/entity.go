package entitys

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/cores/events"
	"github.com/yinweli/Mizugo/cores/msgs"
	"github.com/yinweli/Mizugo/cores/nets"
	"github.com/yinweli/Mizugo/cores/utils"
)

// newEntity 建立實體資料
func newEntity(entityID EntityID) *Entity {
	return &Entity{
		entityID:  entityID,
		modulemgr: NewModulemgr(),
		eventmgr:  events.NewEventmgr(eventSize),
	}
}

// Entity 實體資料
type Entity struct {
	entityID  EntityID                       // 實體編號
	modulemgr *Modulemgr                     // 模組管理器
	eventmgr  *events.Eventmgr               // 事件管理器
	enable    atomic.Bool                    // 啟用旗標
	session   utils.SyncAttr[nets.Sessioner] // 會話物件
	process   utils.SyncAttr[msgs.Processor] // 處理物件
}

// ===== 基礎功能 =====

// Initialize 初始化處理
func (this *Entity) Initialize() error {
	if this.enable.CompareAndSwap(false, true) == false {
		return fmt.Errorf("entity initialize: already initialize")
	} // if

	this.eventmgr.Sub(eventAwake, this.eventAwake)
	this.eventmgr.Sub(eventStart, this.eventStart)
	this.eventmgr.Sub(eventDispose, this.eventDispose)
	this.eventmgr.Sub(eventUpdate, this.eventUpdate)
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

	return nil
}

// Finalize 結束處理, 請不要重複使用結束的實體物件
func (this *Entity) Finalize() error {
	if this.enable.CompareAndSwap(true, false) == false {
		return fmt.Errorf("entity initialize: already finalize or not initialize")
	} // if

	for _, itor := range this.modulemgr.All() {
		itor.Internal().updateStop()
		this.eventmgr.PubOnce(eventDispose, itor)
	} // for

	this.eventmgr.Finalize()
	return nil
}

// EntityID 取得實體編號
func (this *Entity) EntityID() EntityID {
	return this.entityID
}

// Enable 取得啟用旗標
func (this *Entity) Enable() bool {
	return this.enable.Load()
}

// ===== 模組功能 =====

// AddModule 新增模組, 初始化完成後就不能新增模組
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

// ===== 事件功能 =====

// SubEvent 訂閱事件, 初始化完成後就不能訂閱事件
func (this *Entity) SubEvent(name string, process events.Process) error {
	if this.enable.Load() {
		return fmt.Errorf("entity sub event: overdue")
	} // if

	this.eventmgr.Sub(name, process)
	return nil
}

// PubOnceEvent 發布單次事件
func (this *Entity) PubOnceEvent(name string, param any) {
	this.eventmgr.PubOnce(name, param)
}

// PubFixedEvent 發布定時事件, 回傳用於停止定時事件的控制物件
func (this *Entity) PubFixedEvent(name string, param any, interval time.Duration) *events.Fixed {
	return this.eventmgr.PubFixed(name, param, interval)
}

// ===== 會話功能 =====

// SetSession 設定會話物件, 初始化完成後就不能設定會話物件
func (this *Entity) SetSession(session nets.Sessioner) error {
	if this.enable.Load() {
		return fmt.Errorf("entity set session: overdue")
	} // if

	this.session.Set(session)
	return nil
}

// GetSession 取得會話物件
func (this *Entity) GetSession() nets.Sessioner {
	return this.session.Get()
}

// ===== 處理功能 =====

// SetProcess 設定處理物件, 初始化完成後就不能設定處理物件
func (this *Entity) SetProcess(process msgs.Processor) error {
	if this.enable.Load() {
		return fmt.Errorf("entity set process: overdue")
	} // if

	this.process.Set(process)
	return nil
}

// GetProcess 取得處理物件
func (this *Entity) GetProcess() msgs.Processor {
	return this.process.Get()
}

// ===== 內部功能 =====

// eventAwake 處理awake事件
func (this *Entity) eventAwake(param any) {
	if module, ok := param.(Awaker); ok {
		module.Awake()
	} // if
}

// eventStart 處理start事件
func (this *Entity) eventStart(param any) {
	if module, ok := param.(Starter); ok {
		module.Start()
	} // if
}

// eventDispose 處理dispose事件
func (this *Entity) eventDispose(param any) {
	if module, ok := param.(Disposer); ok {
		module.Dispose()
	} // if
}

// eventUpdate 處理update事件
func (this *Entity) eventUpdate(param any) {
	if module, ok := param.(Updater); ok {
		module.Update()
	} // if
}
