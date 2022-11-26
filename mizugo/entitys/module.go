package entitys

import (
	"github.com/yinweli/Mizugo/mizugo/events"
)

// NewModule 建立模組資料
func NewModule(moduleID ModuleID, name string) *Module {
	return &Module{
		moduleID: moduleID,
		name:     name,
	}
}

// Moduler 模組介面
type Moduler interface {
	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Name 取得模組名稱
	Name() string

	// Entity 取得實體物件
	Entity() *Entity

	// Internal 取得內部物件
	Internal() *Internal
}

// Awaker awake介面
type Awaker interface {
	// Awake awake事件, 模組初始化時第一個被執行
	Awake()
}

// Starter start介面
type Starter interface {
	// Start start事件, 模組初始化時第二個被執行
	Start()
}

// Disposer dispose介面
type Disposer interface {
	// Dispose dispose事件, 模組結束時執行
	Dispose()
}

// Updater update介面
type Updater interface {
	// Update update事件, 模組定時事件, 間隔時間定義在updateInterval
	Update()
}

// Module 模組資料
type Module struct {
	moduleID ModuleID // 模組編號
	name     string   // 模組名稱
	internal Internal // 內部物件
}

// ModuleID 模組編號
type ModuleID int64

// ModuleID 取得模組編號
func (this *Module) ModuleID() ModuleID {
	return this.moduleID
}

// Name 取得模組名稱
func (this *Module) Name() string {
	return this.name
}

// Entity 取得實體物件
func (this *Module) Entity() *Entity {
	return this.internal.entity
}

// Internal 取得內部物件
func (this *Module) Internal() *Internal {
	return &this.internal
}

// Internal 內部資料
type Internal struct {
	entity *Entity       // 實體物件
	update *events.Fixed // update事件定時物件
}

// updateStop 停止update事件定時
func (this *Internal) updateStop() {
	if this.update != nil {
		this.update.Stop()
	} // if
}
