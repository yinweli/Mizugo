package entitys

// NewModule 建立模組資料
func NewModule(moduleID ModuleID) *Module {
	return &Module{
		moduleID: moduleID,
	}
}

// Module 模組資料
type Module struct {
	moduleID ModuleID // 模組編號
	entity   *Entity  // 實體物件
}

// Moduler 模組介面
type Moduler interface {
	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Entity 取得實體物件
	Entity() *Entity

	// setup 設定模組
	setup(entity *Entity)
}

// ModuleID 取得模組編號
func (this *Module) ModuleID() ModuleID {
	return this.moduleID
}

// Entity 取得實體物件
func (this *Module) Entity() *Entity {
	return this.entity
}

// setup 設定模組
func (this *Module) setup(entity *Entity) {
	this.entity = entity
}
