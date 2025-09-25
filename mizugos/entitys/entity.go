package entitys

import (
	"fmt"
	"net"
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/nets"
	"github.com/yinweli/Mizugo/v2/mizugos/procs"
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID) *Entity {
	return &Entity{
		entityID:  entityID,
		modulemap: NewModulemap(),
		eventmap:  NewEventmap(),
	}
}

// Entity 實體資料
//
// mizugo 中用於承載「對象」(如: 網路連線或遊戲物件) 的基礎單元, 具備三大能力:
//   - 模組管理: 彈性掛載自訂模組以實作功能/訊息處理等
//   - 事件系統: 支援單次/延遲/定時事件的發布與訂閱
//   - 訊息/會話: 可配置 procs.Processor 處理訊息, 並綁定 nets.Sessioner 進行網路通訊
//
// 生命週期:
//   - 建構: 呼叫 NewEntity 建立; 此時可 AddModule / SetProcess / SetSession
//   - 初始化: 呼叫 Initialize, 依序呼叫模組的 Awaker.Awake → Starter.Start, 並啟動事件系統
//   - 結束: 呼叫 Finalize, 依序發布 EventDispose → EventShutdown, 釋放資源
//
// Initialize 前應完成 AddModule / SetProcess / SetSession; Initialize 後將禁止再調整
//
// 事件模式:
//   - 單次事件: 只執行一次
//   - 延遲事件: 延遲一段時間後執行一次
//   - 定時事件: 以固定間隔重複執行; 注意: 定時事件無法刪除, 發布前請審慎評估
//
// 內建事件 (由實體在適當時機發布):
//   - EventDispose: 結束事件, Finalize 時首先發布, 參數為 nil
//   - EventShutdown: 關閉事件, Finalize 時第二個發布, 參數為 nil (此時連線已中斷)
//   - EventRecv: 接收訊息事件, 收到訊息時觸發, 參數為訊息物件
//   - EventSend: 傳送訊息事件, 送出訊息時觸發, 參數為訊息物件
type Entity struct {
	entityID  EntityID                        // 實體編號
	modulemap *Modulemap                      // 模組列表
	eventmap  *Eventmap                       // 事件列表
	once      helps.SyncOnce                  // 單次執行物件
	process   helps.SyncAttr[procs.Processor] // 處理物件
	session   helps.SyncAttr[nets.Sessioner]  // 會話物件
}

// ===== 基礎功能 =====

// Initialize 初始化處理, 僅會成功執行一次
//
// 流程:
//   - 對所有模組依序呼叫 Awaker.Awake → Starter.Start
//   - 初始化事件系統, 並訂閱 EventRecv, EventShutdown 事件供內部使用
//
// 失敗時(包括已初始化過)回傳錯誤
func (this *Entity) Initialize(wrong Wrong) (err error) {
	if this.once.Done() {
		return fmt.Errorf("entity initialize: already initialize")
	} // if

	this.once.Do(func() {
		module := this.modulemap.All()

		for _, itor := range module {
			if a, ok := itor.(Awaker); ok {
				if err = a.Awake(); err != nil {
					return
				} // if
			} // if
		} // for

		for _, itor := range module {
			if s, ok := itor.(Starter); ok {
				if err = s.Start(); err != nil {
					return
				} // if
			} // if
		} // for

		if err = this.eventmap.Initialize(); err != nil {
			return
		} // if

		this.eventmap.Sub(EventRecv, func(param any) {
			if err := this.process.Get().Process(param); err != nil {
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
//
// 流程:
//   - 僅在已完成 Initialize 的前提下生效
//   - 依序發布 EventDispose → EventShutdown
//   - 結束事件系統
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

// Enable 取得實體是否已完成初始化
func (this *Entity) Enable() bool {
	return this.once.Done()
}

// ===== 模組功能 =====

// AddModule 新增模組, 僅能在 Initialize 前加入
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

// Subscribe 訂閱事件並回傳訂閱 ID
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

// PublishDelay 發布延遲事件, 該事件會延遲一段時間才發布, 僅發布一次, 不會重複
func (this *Entity) PublishDelay(name string, param any, delay time.Duration) {
	this.eventmap.PubDelay(name, param, delay)
}

// PublishFixed 發布定時事件, 該事件會週期性地發布, 由於不能刪除定時事件, 發布前請先評估需求
func (this *Entity) PublishFixed(name string, param any, interval time.Duration) {
	this.eventmap.PubFixed(name, param, interval)
}

// ===== 處理功能 =====

// SetProcess 設定處理物件, 僅能在 Initialize 前設定
//
// 之後可透過 AddMessage / DelMessage / GetMessage 管理訊息處理流程
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

// GetMessage 取得訊息處理
func (this *Entity) GetMessage(messageID int32) procs.Process {
	return this.process.Get().Get(messageID)
}

// ===== 會話功能 =====

// SetSession 設定會話物件, 僅能在 Initialize 前設定
//
// 綁定後, 可使用 Stop / StopWait / Send / LocalAddr / RemoteAddr 進行網路操作
func (this *Entity) SetSession(session nets.Sessioner) error {
	if this.once.Done() {
		return fmt.Errorf("entity set session: overdue")
	} // if

	this.session.Set(session)
	return nil
}

// Stop 立即要求停止會話, 不等待會話內部循環結束
func (this *Entity) Stop() {
	this.session.Get().Stop()
}

// StopWait 停止會話並等待會話內部循環安全結束
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
