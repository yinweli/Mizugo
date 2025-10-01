package entitys

// NewModule 建立模組資料
func NewModule(moduleID ModuleID) *Module {
	return &Module{
		moduleID: moduleID,
	}
}

// Module 模組資料, 用於分類與實作遊戲功能/訊息處理等
//
// 使用時, 需要遵循以下流程
//   - 定義模組結構, 並且把 Module 作為模組的第一個成員
//   - (可選)繼承 Awaker 介面, 並實作介面中的 Awake 函式, 其中可以填寫模組的初始化
//   - (可選)繼承 Starter 介面, 並實作介面中的 Start 函式, 其中可以填寫模組的初始化
//   - 建立模組資料, 並把模組加入到實體中; 需要在實體初始化之前完成
//
// 兩個初始化介面的執行順序為 Awaker → Starter
type Module struct {
	moduleID ModuleID // 模組編號
	entity   *Entity  // 實體物件
}

// Moduler 模組介面
type Moduler interface {
	// initialize 初始化處理
	initialize(entity *Entity)

	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Entity 取得實體物件
	Entity() *Entity
}

// Awaker 模組喚醒介面
type Awaker interface {
	// Awake 喚醒處理, 模組初始化時第一個被執行
	Awake() error
}

// Starter 模組啟動介面
type Starter interface {
	// Start 啟動處理, 模組初始化時第二個被執行
	Start() error
}

// initialize 初始化處理
func (this *Module) initialize(entity *Entity) {
	this.entity = entity
}

// ModuleID 取得模組編號
func (this *Module) ModuleID() ModuleID {
	return this.moduleID
}

// Entity 取得實體物件
func (this *Module) Entity() *Entity {
	return this.entity
}
