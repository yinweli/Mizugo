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
}

// Awaker awake介面
type Awaker interface {
	// Awake 模組喚醒通知, 模組加入實體後第一個被觸發的通知
	Awake()
}

// Starter start介面
type Starter interface {
	// Start 模組啟動通知, 模組加入實體後第二個被觸發的通知
	Start()
}

// Disposer dispose介面
type Disposer interface {
	// Dispose 模組結束通知, 模組被移出實體後觸發的通知
	Dispose()
}

// Updater update介面
type Updater interface {
	// Update 模組定時通知, 模組加入實體後, 定時觸發的通知, 間隔時間定義在updateInterval
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
