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

// Module 模組資料
type Module struct {
	moduleID ModuleID      // 模組編號
	name     string        // 模組名稱
	entity   *Entity       // 實體物件
	fixed    *events.Fixed // 定時事件物件
}

// ModuleID 模組編號
type ModuleID int64

// Moduler 模組介面
type Moduler interface {
	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Name 取得模組名稱
	Name() string

	// Entity 取得實體物件
	Entity() *Entity

	// Host 設定宿主實體
	Host(entity *Entity)

	// Fixed 設定定時控制器
	Fixed(fixed *events.Fixed)

	// FixedStop 停止定時控制器
	FixedStop()
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
	return this.entity
}

// Host 設定宿主實體
func (this *Module) Host(entity *Entity) {
	this.entity = entity
}

// Fixed 設定定時控制器
func (this *Module) Fixed(fixed *events.Fixed) {
	this.fixed = fixed
}

// FixedStop 停止定時控制器
func (this *Module) FixedStop() {
	if this.fixed != nil {
		this.fixed.Stop()
	} // if
}
