package entitys

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yinweli/Mizugo/mizugo/events"
)

// NewEntity 建立實體資料
func NewEntity(entityID EntityID, name string) *Entity {
	return &Entity{
		entityID: entityID,
		name:     name,
		modulean: NewModulean(),
		eventan:  events.NewEventan(eventBufferSize),
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
		this.eventan.PubOnce(eventAwake, module)
		this.eventan.PubOnce(eventStart, module)
		module.Fixed(this.eventan.PubFixed(eventUpdate, module, updateInterval))
	} // if

	return nil
}

// DelModule 刪除模組
func (this *Entity) DelModule(moduleID ModuleID) Moduler {
	if module := this.modulean.Del(moduleID); module != nil {
		this.eventan.PubOnce(eventDispose, module)
		module.FixedStop()
		return module
	} // if

	return nil
}

// GetModule 取得模組
func (this *Entity) GetModule(moduleID ModuleID) Moduler {
	return this.modulean.Get(moduleID)
}

// SubEvent 訂閱事件, 由於初始化完成後就會開始處理事件, 因此可能需要在初始化之前做完訂閱事件
func (this *Entity) SubEvent(name string, process events.Process) {
	this.eventan.Sub(name, process)
}

// PubOnceEvent 發布單次事件
func (this *Entity) PubOnceEvent(name string, param any) {
	this.eventan.PubOnce(name, param)
}

// PubFixedEvent 發布定時事件, 回傳用於停止定時事件的控制物件
func (this *Entity) PubFixedEvent(name string, param any, interval time.Duration) *events.Fixed {
	return this.eventan.PubFixed(name, param, interval)
}

// initialize 初始化處理
func (this *Entity) initialize() {
	if this.enable.CompareAndSwap(false, true) {
		this.eventan.Sub(eventAwake, processAwake)
		this.eventan.Sub(eventStart, processStart)
		this.eventan.Sub(eventDispose, processDispose)
		this.eventan.Sub(eventUpdate, processUpdate)
		this.eventan.Initialize()
		module := this.modulean.All()

		for _, itor := range module {
			this.eventan.PubOnce(eventAwake, itor)
		} // for

		for _, itor := range module {
			this.eventan.PubOnce(eventStart, itor)
		} // for

		for _, itor := range module {
			itor.Fixed(this.eventan.PubFixed(eventUpdate, itor, updateInterval))
		} // for
	} // if
}

// finalize 結束處理
func (this *Entity) finalize() {
	for _, itor := range this.modulean.All() {
		this.eventan.PubOnce(eventDispose, itor)
		itor.FixedStop()
	} // for

	this.enable.Store(false)
	this.eventan.Finalize()
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
