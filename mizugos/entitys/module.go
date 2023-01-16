package entitys

// 模組, mizugo中用於存放遊戲功能的資料或是處理函式
// * 建立模組流程
//   建立模組結構, 並且把Module作為模組的第一個成員
//   建立Awake函式, 填入初始化步驟到Awake函式中
//   建立Start函式, 填入初始化步驟到Start函式中
//   以此模組結構宣告模組物件, 並加入到實體中
// * 內部事件
//   實體提供了內部事件可供訂閱, 內部事件請參考define.go中的說明

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
	// initialize 初始化處理
	initialize(entity *Entity)

	// ModuleID 取得模組編號
	ModuleID() ModuleID

	// Entity 取得實體物件
	Entity() *Entity

	// Awake 模組初始化時第一個被執行
	Awake() error

	// Start 模組初始化時第二個被執行
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
