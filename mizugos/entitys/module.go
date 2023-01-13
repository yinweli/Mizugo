package entitys

// 模組, mizugo中用於存放遊戲功能的資料或是處理函式
// * 建立模組流程
//   建立模組結構, 並且把Module作為模組的第一個成員
//   如果需要執行Awake事件, 則替模組結構建立Awake函式
//     Awake() error
//   如果需要執行Start事件, 則替模組結構建立Start函式
//     Start() error
//   以此模組結構宣告模組物件, 並加入到實體中
// * 模組事件
//   - Awake: 通過為模組結構添加Awake函式完成, 是模組初始化第一個被執行的事件
//   - Start: 通過為模組結構添加Start函式完成, 是模組初始化第二個被執行的事件
//   實際上實體會先執行底下所有模組的Awake事件, 然後再執行所有模組的Start事件
// * 內部事件
//   實體提供了以下內部事件可供訂閱
//   - update: 每updateInterval觸發一次
//   - dispose: 實體結束時執行
//   - afterSend: 傳送訊息結束後執行
//   - afterRecv: 接收訊息結束後執行

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
