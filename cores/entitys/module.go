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

	// SetEntity 設定實體物件
	SetEntity(entity *Entity)

	// GetEntity 取得實體物件
	GetEntity() *Entity
}

// ModuleID 取得模組編號
func (this *Module) ModuleID() ModuleID {
	return this.moduleID
}

// SetEntity 設定實體物件
func (this *Module) SetEntity(entity *Entity) {
	this.entity = entity
}

// GetEntity 取得實體物件
func (this *Module) GetEntity() *Entity {
	return this.entity
}
