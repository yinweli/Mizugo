package entitys

import (
	"fmt"
	"net"
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
		Labelobj: labels.NewLabelobj(),
		entityID: entityID,
	}
}

// Entity 實體資料, mizugo中用於儲存對象的基礎物件, 對象可以是個連線, 也可以用於表示遊戲物件
//
// 建立實體時, 需要遵循以下流程
//   - 到實體管理器新增實體
//   - 取得實體物件
//   - 如果實體需要使用模組相關功能, 則呼叫 SetModulemgr 設置 entitys.Modulemgr
//   - 如果實體需要使用事件相關功能, 則呼叫 SetEventmgr 設置 events.Eventmgr
//   - 如果實體將要代表某個連線, 則呼叫 則需要設置 procs.Processor 與 nets.Sessioner
//   - 新增模組到實體中
//   - 執行實體的初始化處理
//
// 結束實體時, 會依序發布 EventDispose 與 EventShutdown 事件
//
// 實體可以新增模組, 用於分類與實作遊戲功能/訊息處理等; 如果想要新增模組, 需要在實體初始化之前完成;
// 需要在實體初始化之前設置 entitys.Modulemgr
//
// 使用者可以到實體訂閱事件, 事件可以被訂閱多次, 發布事件時會每個訂閱者都會執行一次;
// 需要在實體初始化之前設置 events.Eventmgr
//
// 使用者可以到實體發布事件, 發布事件有以下模式
//   - 單次事件: 事件只執行一次
//   - 延遲事件: 事件只執行一次, 且會延遲一段時間才發布
//   - 定時事件: 事件會定時執行, 由於不能刪除定時事件, 因此發布定時事件前請多想想
//
// 實體會發布以下內部事件
//   - EventUpdate: 定時事件, 實體定時觸發, 參數是nil, 間隔時間為 UpdateInterval
//   - EventDispose: 結束事件, 實體結束時第一個執行, 參數是nil
//   - EventShutdown: 關閉事件, 實體結束時第二個執行, 參數是nil, 這時連線已經中斷
//   - EventRecv: 接收訊息事件, 當接收訊息後觸發, 參數是訊息物件
//   - EventSend: 傳送訊息事件, 當傳送訊息後觸發, 參數是訊息物件
//
// 實體可以設置處理功能, 負責封包編/解碼, 管理訊息處理函式;
// 需要在實體初始化之前設置 procs.Processor
//
// 實體可以設置會話功能, 負責網路相關功能;
// 需要在實體初始化之前設置 nets.Sessioner
type Entity struct {
	*labels.Labelobj                                  // 標籤物件
	entityID         EntityID                         // 實體編號
	once             utils.SyncOnce                   // 單次執行物件
	modulemgr        utils.SyncAttr[*Modulemgr]       // 模組管理器
	eventmgr         utils.SyncAttr[*events.Eventmgr] // 事件管理器
	process          utils.SyncAttr[procs.Processor]  // 處理物件
	session          utils.SyncAttr[nets.Sessioner]   // 會話物件
}

// ===== 基礎功能 =====

// Initialize 初始化處理
func (this *Entity) Initialize(wrong Wrong) (err error) {
	if this.once.Done() {
		return fmt.Errorf("entity initialize: already initialize")
	} // if

	this.once.Do(func() {
		modulemgr := this.modulemgr.Get()

		if modulemgr == nil {
			err = fmt.Errorf("entity initialize: modulemgr nil")
			return
		} // if

		module := modulemgr.All()

		for _, itor := range module {
			if awaker, ok := itor.(Awaker); ok {
				if err = awaker.Awake(); err != nil {
					err = fmt.Errorf("entity initialize: %w", err)
					return
				} // if
			} // if
		} // for

		for _, itor := range module {
			if starter, ok := itor.(Starter); ok {
				if err = starter.Start(); err != nil {
					err = fmt.Errorf("entity initialize: %w", err)
					return
				} // if
			} // if
		} // for

		eventmgr := this.eventmgr.Get()

		if eventmgr == nil {
			err = fmt.Errorf("entity initialize: eventmgr nil")
			return
		} // if

		eventmgr.Sub(EventRecv, func(param any) {
			if err := this.process.Get().Process(param); err != nil {
				wrong.Do(fmt.Errorf("entity recv: %w", err))
			} // if
		})
		eventmgr.Sub(EventShutdown, func(_ any) {
			if session := this.session.Get(); session != nil {
				session.Stop()
			} // if
		})
		eventmgr.PubFixed(EventUpdate, nil, UpdateInterval)

		if err = eventmgr.Initialize(); err != nil {
			err = fmt.Errorf("entity initialize: %w", err)
			return
		} // if
	})

	return err
}

// Finalize 結束處理
func (this *Entity) Finalize() {
	if this.once.Done() == false {
		return
	} // if

	if eventmgr := this.eventmgr.Get(); eventmgr != nil {
		eventmgr.PubOnce(EventDispose, nil)
		eventmgr.PubOnce(EventShutdown, nil)
		eventmgr.Finalize()
	} // if
}

// Bundle 取得綁定資料
func (this *Entity) Bundle() *nets.Bundle {
	return &nets.Bundle{
		Encode:  this.process.Get().Encode,
		Decode:  this.process.Get().Decode,
		Publish: this.eventmgr.Get().PubOnce,
	}
}

// EntityID 取得實體編號
func (this *Entity) EntityID() EntityID {
	return this.entityID
}

// Enable 取得啟用旗標
func (this *Entity) Enable() bool {
	return this.once.Done()
}

// ===== 模組功能 =====

// SetModulemgr 設定模組管理器, 初始化完成後就不能設定模組物件
func (this *Entity) SetModulemgr(modulemgr *Modulemgr) error {
	if this.once.Done() {
		return fmt.Errorf("entity set modulemgr: overdue")
	} // if

	this.modulemgr.Set(modulemgr)
	return nil
}

// GetModulemgr 取得模組管理器
func (this *Entity) GetModulemgr() *Modulemgr {
	return this.modulemgr.Get()
}

// AddModule 新增模組, 初始化完成後就不能新增模組
func (this *Entity) AddModule(module Moduler) error {
	if this.once.Done() {
		return fmt.Errorf("entity add module: overdue")
	} // if

	if err := this.modulemgr.Get().Add(module); err != nil {
		return fmt.Errorf("entity add module: %w", err)
	} // if

	module.initialize(this)
	return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) Moduler {
	return this.modulemgr.Get().Get(moduleID)
}

// ===== 事件功能 =====

// SetEventmgr 設定事件管理器, 初始化完成後就不能設定事件物件
func (this *Entity) SetEventmgr(eventmgr *events.Eventmgr) error {
	if this.once.Done() {
		return fmt.Errorf("entity set eventmgr: overdue")
	} // if

	this.eventmgr.Set(eventmgr)
	return nil
}

// GetEventmgr 取得事件管理器
func (this *Entity) GetEventmgr() *events.Eventmgr {
	return this.eventmgr.Get()
}

// Subscribe 訂閱事件
func (this *Entity) Subscribe(name string, process events.Process) string {
	return this.eventmgr.Get().Sub(name, process)
}

// Unsubscribe 取消訂閱事件
func (this *Entity) Unsubscribe(subID string) {
	this.eventmgr.Get().Unsub(subID)
}

// PublishOnce 發布單次事件
func (this *Entity) PublishOnce(name string, param any) {
	this.eventmgr.Get().PubOnce(name, param)
}

// PublishDelay 發布延遲事件, 事件會延遲一段時間才發布, 但仍是單次事件
func (this *Entity) PublishDelay(name string, param any, delay time.Duration) {
	this.eventmgr.Get().PubDelay(name, param, delay)
}

// PublishFixed 發布定時事件, 請注意! 由於不能刪除定時事件, 因此發布定時事件前請多想想
func (this *Entity) PublishFixed(name string, param any, interval time.Duration) {
	this.eventmgr.Get().PubFixed(name, param, interval)
}

// ===== 處理功能 =====

// SetProcess 設定處理物件, 初始化完成後就不能設定處理物件
func (this *Entity) SetProcess(process procs.Processor) error {
	if this.once.Done() {
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

// ===== 會話功能 =====

// SetSession 設定會話物件, 初始化完成後就不能設定會話物件
func (this *Entity) SetSession(session nets.Sessioner) error {
	if this.once.Done() {
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
