package entitys

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/mizugos/events"
	"github.com/yinweli/Mizugo/mizugos/labels"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
	"github.com/yinweli/Mizugo/mizugos/utils"
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID) *Entity {
	return &Entity{
		entityID:  entityID,
		modulemgr: NewModulemgr(),
		eventmgr:  events.NewEventmgr(eventSize),
		labelobj:  labels.NewLabelobj(),
	}
}

// Entity 實體資料
type Entity struct {
	entityID  EntityID                        // 實體編號
	enable    atomic.Bool                     // 啟用旗標
	close     []func()                        // 結束處理列表
	session   utils.SyncAttr[nets.Sessioner]  // 會話物件
	process   utils.SyncAttr[procs.Processor] // 處理物件
	modulemgr *Modulemgr                      // 模組管理器
	eventmgr  *events.Eventmgr                // 事件管理器
	labelobj  *labels.Labelobj                // 標籤物件
}

// ===== 基礎功能 =====

// Initialize 初始化處理
func (this *Entity) Initialize(closes ...func()) error {
	if this.enable.CompareAndSwap(false, true) == false {
		return fmt.Errorf("entity initialize: already initialize")
	} // if

	this.close = closes
	this.eventmgr.Sub(eventAwake, func(param any) {
		if module, ok := param.(Awaker); ok {
			module.Awake()
		} // if
	})
	this.eventmgr.Sub(eventStart, func(param any) {
		if module, ok := param.(Starter); ok {
			module.Start()
		} // if
	})
	this.eventmgr.Sub(eventUpdate, func(param any) {
		if module, ok := param.(Updater); ok {
			module.Update()
		} // if
	})
	this.eventmgr.Sub(eventDispose, func(param any) {
		if module, ok := param.(Disposer); ok {
			module.Dispose()
		} // if
	})
	this.eventmgr.Sub(eventFinalize, func(_ any) {
		for _, itor := range this.close {
			itor()
		} // for

		if session := this.session.Get(); session != nil {
			session.Stop()
		} // if
	})
	this.eventmgr.Initialize()
	module := this.modulemgr.All()

	for _, itor := range module {
		this.eventmgr.PubOnce(eventAwake, itor)
	} // for

	for _, itor := range module {
		this.eventmgr.PubOnce(eventStart, itor)
	} // for

	for _, itor := range module {
		this.eventmgr.PubFixed(eventUpdate, itor, updateInterval)
	} // for

	return nil
}

// Finalize 結束處理
func (this *Entity) Finalize() {
	if this.enable.CompareAndSwap(true, false) == false {
		return
	} // if

	for _, itor := range this.modulemgr.All() {
		this.eventmgr.PubOnce(eventDispose, itor)
	} // for

	this.eventmgr.PubOnce(eventFinalize, nil)
	this.eventmgr.Finalize()
}

// EntityID 取得實體編號
func (this *Entity) EntityID() EntityID {
	return this.entityID
}

// Enable 取得啟用旗標
func (this *Entity) Enable() bool {
	return this.enable.Load()
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

// Send 傳送封包
func (this *Entity) Send(message any) {
	this.session.Get().Send(message)
}

// SessionID 取得會話編號
func (this *Entity) SessionID() nets.SessionID {
	return this.session.Get().SessionID()
}

// RemoteAddr 取得遠端位址
func (this *Entity) RemoteAddr() net.Addr {
	return this.session.Get().RemoteAddr()
}

// LocalAddr 取得本地位址
func (this *Entity) LocalAddr() net.Addr {
	return this.session.Get().LocalAddr()
}

// ===== 處理功能 =====

// SetProcess 設定處理物件
func (this *Entity) SetProcess(process procs.Processor) error {
	this.process.Set(process)
	return nil
}

// GetProcess 取得處理物件
func (this *Entity) GetProcess() procs.Processor {
	return this.process.Get()
}

// AddMessage 新增訊息處理
func (this *Entity) AddMessage(messageID procs.MessageID, process procs.Process) {
	this.process.Get().Add(messageID, process)
}

// DelMessage 刪除訊息處理
func (this *Entity) DelMessage(messageID procs.MessageID) {
	this.process.Get().Del(messageID)
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

	module.setup(this)
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

	for _, itor := range []string{
		eventAwake,
		eventStart,
		eventUpdate,
		eventDispose,
		eventFinalize,
	} {
		if name == itor {
			return fmt.Errorf("entity sub event: keyword")
		} // if
	} // for

	this.eventmgr.Sub(name, process)
	return nil
}

// PubOnceEvent 發布單次事件
func (this *Entity) PubOnceEvent(name string, param any) {
	this.eventmgr.PubOnce(name, param)
}

// PubFixedEvent 發布定時事件, 回傳用於停止定時事件的控制物件
func (this *Entity) PubFixedEvent(name string, param any, interval time.Duration) {
	this.eventmgr.PubFixed(name, param, interval)
}
