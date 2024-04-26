package entitys

import (
	"fmt"
	"net"
	"time"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/labels"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/procs"
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID) *Entity {
	return &Entity{
		Label:     labels.NewLabel(),
		entityID:  entityID,
		modulemap: NewModulemap(),
		eventmap:  NewEventmap(),
	}
}

// Entity 實體資料, mizugo中用於儲存對象的基礎物件, 對象可以是個連線, 也可以用於表示遊戲物件
//
// 使用者可以到實體新增模組, 用於分類與實作遊戲功能/訊息處理等
//
// 使用者可以到實體訂閱事件, 用於接收實體發布的事件
//
// 使用者可以到實體發布事件, 發布事件有以下模式
//   - 單次事件: 事件只執行一次
//   - 延遲事件: 事件只執行一次, 且會延遲一段時間才發布
//   - 定時事件: 事件會定時執行, 由於不能刪除定時事件, 因此發布定時事件前請多想想
//
// 實體會發布以下內部事件
//   - EventDispose: 結束事件, 實體結束時第一個執行, 參數是nil
//   - EventShutdown: 關閉事件, 實體結束時第二個執行, 參數是nil, 這時連線已經中斷
//   - EventRecv: 接收訊息事件, 當接收訊息後觸發, 參數是訊息物件
//   - EventSend: 傳送訊息事件, 當傳送訊息後觸發, 參數是訊息物件
//
// 實體可以設置處理功能, 負責訊息處理功能; 需要在實體初始化之前設置 procs.Processor
//
// 實體可以設置會話功能, 負責網路相關功能; 需要在實體初始化之前設置 nets.Sessioner
//
// 結束實體時, 需要執行 Finalize
type Entity struct {
	*labels.Label                                 // 標籤物件
	entityID      EntityID                        // 實體編號
	modulemap     *Modulemap                      // 模組列表
	eventmap      *Eventmap                       // 事件列表
	once          helps.SyncOnce                  // 單次執行物件
	process       helps.SyncAttr[procs.Processor] // 處理物件
	session       helps.SyncAttr[nets.Sessioner]  // 會話物件
}

// ===== 基礎功能 =====

// Initialize 初始化處理
func (this *Entity) Initialize(wrong Wrong) (err error) {
	if this.once.Done() {
		return fmt.Errorf("entity initialize: already initialize")
	} // if

	this.once.Do(func() {
		module := this.modulemap.All()

		for _, itor := range module {
			if awaker, ok := itor.(Awaker); ok {
				if err = awaker.Awake(); err != nil {
					return
				} // if
			} // if
		} // for

		for _, itor := range module {
			if starter, ok := itor.(Starter); ok {
				if err = starter.Start(); err != nil {
					return
				} // if
			} // if
		} // for

		if err = this.eventmap.Initialize(); err != nil {
			return
		} // if

		this.eventmap.Sub(EventRecv, func(param any) {
			if err = this.process.Get().Process(param); err != nil {
				wrong.Do(fmt.Errorf("entity recv: %w", err))
			} // if
		})
		this.eventmap.Sub(EventShutdown, func(_ any) {
			if session := this.session.Get(); session != nil {
				session.Stop()
			} // if
		})
	})

	if err != nil {
		return fmt.Errorf("entity initialize: %w", err)
	} // if

	return nil
}

// Finalize 結束處理
func (this *Entity) Finalize() {
	if this.once.Done() == false {
		return
	} // if

	this.eventmap.PubOnce(EventDispose, nil)
	this.eventmap.PubOnce(EventShutdown, nil)
	this.eventmap.Finalize()
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

// AddModule 新增模組, 初始化完成後就不能新增模組
func (this *Entity) AddModule(module Moduler) error {
	if this.once.Done() {
		return fmt.Errorf("entity add module: overdue")
	} // if

	if err := this.modulemap.Add(module); err != nil {
		return fmt.Errorf("entity add module: %w", err)
	} // if

	module.initialize(this)
	return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) Moduler {
	return this.modulemap.Get(moduleID)
}

// ===== 事件功能 =====

// Subscribe 訂閱事件
func (this *Entity) Subscribe(name string, process Process) string {
	return this.eventmap.Sub(name, process)
}

// Unsubscribe 取消訂閱事件
func (this *Entity) Unsubscribe(subID string) {
	this.eventmap.Unsub(subID)
}

// PublishOnce 發布單次事件
func (this *Entity) PublishOnce(name string, param any) {
	this.eventmap.PubOnce(name, param)
}

// PublishDelay 發布延遲事件, 事件會延遲一段時間才發布, 但仍是單次事件
func (this *Entity) PublishDelay(name string, param any, delay time.Duration) {
	this.eventmap.PubDelay(name, param, delay)
}

// PublishFixed 發布定時事件, 請注意! 由於不能刪除定時事件, 因此發布定時事件前請多想想
func (this *Entity) PublishFixed(name string, param any, interval time.Duration) {
	this.eventmap.PubFixed(name, param, interval)
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

// AddMessage 新增訊息處理
func (this *Entity) AddMessage(messageID int32, process procs.Process) {
	this.process.Get().Add(messageID, process)
}

// DelMessage 刪除訊息處理
func (this *Entity) DelMessage(messageID int32) {
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

// Stop 停止會話, 不會等待會話內部循環結束
func (this *Entity) Stop() {
	this.session.Get().Stop()
}

// StopWait 停止會話, 會等待會話內部循環結束
func (this *Entity) StopWait() {
	this.session.Get().StopWait()
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
