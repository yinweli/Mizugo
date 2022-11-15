package entitys

// NewModule 建立模組資料
func NewModule(moduleID ModuleID, name string) *Module {
	return &Module{
		moduleID: moduleID,
		name:     name,
	}
}

// Module 模組資料
type Module struct {
	moduleID ModuleID // 模組編號
	name     string   // 模組名稱
	entity   *Entity  // 實體物件
}

// ModuleID 模組編號
type ModuleID int64

// ModuleInterface 模組介面
type ModuleInterface interface {
	ModuleID() ModuleID  // 取得模組編號
	Name() string        // 取得模組名稱
	Entity() *Entity     // 取得實體物件
	Host(entity *Entity) // 設定宿主實體
}

// ModuleAwake 模組awake介面
type ModuleAwake interface {
	Awake()
}

// ModuleStart 模組start介面
type ModuleStart interface {
	Start()
}

// ModuleDispose 模組dispose介面
type ModuleDispose interface {
	Dispose()
}

// ModuleUpdate 模組update介面
type ModuleUpdate interface {
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
