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
	session   utils.SyncAttr[nets.Sessioner]  // 會話物件
	process   utils.SyncAttr[procs.Processor] // 處理物件
	modulemgr *Modulemgr                      // 模組管理器
	eventmgr  *events.Eventmgr                // 事件管理器
	labelobj  *labels.Labelobj                // 標籤物件
}

// ===== 基礎功能 =====

// Initialize 初始化處理
func (this *Entity) Initialize() error {
	if this.enable.CompareAndSwap(false, true) == false {
		return fmt.Errorf("entity initialize: already initialize")
	} // if

	module := this.modulemgr.All()

	for _, itor := range module {
		if awaker, ok := itor.(Awaker); ok {
			awaker.Awake()
		} // if
	} // for

	for _, itor := range module {
		if starter, ok := itor.(Starter); ok {
			starter.Start()
		} // if
	} // for

	this.eventmgr.Sub(EventFinalize, func(_ any) {
		if session := this.session.Get(); session != nil {
			session.Stop()
		} // if
	})
	this.eventmgr.PubFixed(EventUpdate, nil, updateInterval)
	this.eventmgr.Initialize()
	return nil
}

// Finalize 結束處理
func (this *Entity) Finalize() {
	if this.enable.CompareAndSwap(true, false) == false {
		return
	} // if

	this.eventmgr.PubOnce(EventDispose, nil)
	this.eventmgr.PubOnce(EventFinalize, nil)
	this.eventmgr.Finalize()
}

// Bundle 取得綁定資料
func (this *Entity) Bundle() nets.Bundle {
	return nets.Bundle{
		Encode:  this.process.Get().Encode,
		Decode:  this.process.Get().Decode,
		Receive: this.process.Get().Process,
		AfterSend: func() {
			this.eventmgr.PubOnce(EventAfterSend, nil)
		},
		AfterRecv: func() {
			this.eventmgr.PubOnce(EventAfterRecv, nil)
		},
	}
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

// RemoteAddr 取得遠端位址
func (this *Entity) RemoteAddr() net.Addr {
	return this.session.Get().RemoteAddr()
}

// LocalAddr 取得本地位址
func (this *Entity) LocalAddr() net.Addr {
	return this.session.Get().LocalAddr()
}

// ===== 處理功能 =====

// SetProcess 設定處理物件, 初始化完成後就不能設定處理物件
func (this *Entity) SetProcess(process procs.Processor) error {
	if this.enable.Load() {
		return fmt.Errorf("entity set process: overdue")
	} // if

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

// SubEvent 訂閱事件
func (this *Entity) SubEvent(name string, process events.Process) (eventID events.EventID, err error) {
	if name == EventFinalize {
		return eventID, fmt.Errorf("entity sub event: can't sub finalize")
	} // if

	return this.eventmgr.Sub(name, process), nil
}

// UnsubEvent 取消訂閱事件
func (this *Entity) UnsubEvent(eventID events.EventID) {
	this.eventmgr.Unsub(eventID)
}

// PubOnceEvent 發布單次事件
func (this *Entity) PubOnceEvent(name string, param any) {
	this.eventmgr.PubOnce(name, param)
}

// PubFixedEvent 發布定時事件; 請注意! 由於不能刪除定時事件, 因此發布定時事件前請多想想
func (this *Entity) PubFixedEvent(name string, param any, interval time.Duration) {
	this.eventmgr.PubFixed(name, param, interval)
}
