package entitys

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/cores/events"
	"github.com/yinweli/Mizugo/cores/nets"
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
	entityID  EntityID         // 實體編號
	enable    atomic.Bool      // 啟用旗標
	session   SessionAttr      // 會話物件
	modulemgr *Modulemgr       // 模組管理器
	eventmgr  *events.Eventmgr // 事件管理器
}

// EntityID 取得實體編號
func (this *Entity) EntityID() EntityID {
	return this.entityID
}

// Enable 取得啟用旗標
func (this *Entity) Enable() bool {
	return this.enable.Load()
}

// SetSession 設定會話物件
func (this *Entity) SetSession(session nets.Sessioner) error {
	if this.enable.Load() {
		return fmt.Errorf("entity set session: overdue")
	} // if

	this.session.Set(session)
	return nil
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

// initialize 初始化處理
func (this *Entity) initialize() {
	if this.enable.CompareAndSwap(false, true) {
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

// TODO: 新增封包管理器, 在bind的時候, 以封包管理器當作參數建立合適的reactor(或是實體本身就有合適的介面可以塞到reactor中)
