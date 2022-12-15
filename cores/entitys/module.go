package entitys

import (
	"github.com/yinweli/Mizugo/cores/events"
)

// NewModule 建立模組資料
func NewModule(moduleID ModuleID) *Module {
	return &Module{
		moduleID: moduleID,
	}
}

// Moduler 模組介面
type Moduler interface {
	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Entity 取得實體物件
	Entity() *Entity

	// Internal 取得內部物件
	Internal() *Internal
}

// Module 模組資料
type Module struct {
	moduleID ModuleID // 模組編號
	internal Internal // 內部物件
}

// ModuleID 模組編號
type ModuleID int64

// ModuleID 取得模組編號
func (this *Module) ModuleID() ModuleID {
	return this.moduleID
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
