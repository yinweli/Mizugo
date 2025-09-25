package nets

import (
	"sync"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
)

// NewNetmgr 建立網路管理器
func NewNetmgr() *Netmgr {
	return &Netmgr{
		connectmgr: newConnectmgr(),
		listenmgr:  newListenmgr(),
		sessionmgr: newSessionmgr(),
	}
}

// Netmgr 網路管理器
//
// 用於集中管理「連接器(Connecter)」, 「接聽器(Listener)」與「會話(Sessioner)」
// 會話僅在內部維護, 不對外直接提供存取或操作
//
// 使用流程:
//   - AddConnectTCP / AddListenTCP → wrapperBind → bind 成功後加入 sessionmgr
//   - Stop → 依序停止 connectmgr / listenmgr / sessionmgr (細節見各 clear 實作)
//
// 使用情境:
//   - 使用 AddConnectTCP 啟動連線, 或使用 AddListenTCP 啟動監聽
//   - 使用 DelConnect / DelListen 移除並停止特定資源
//   - 使用 Stop 一次性關閉所有資源(連線, 接聽, 會話)
//
// Bind 常見流程:
//   - 建立並配置實體(處理器, 會話, 模組)
//   - 執行實體初始化
//   - 設定會話(編碼/解碼, 事件發布, 錯誤處理, 封包大小, 擁有者)
//   - 定義錯誤處理邏輯
//
// Unbind 常見流程:
//   - 釋放實體相關資源
//   - 移除實體資料
type Netmgr struct {
	connectmgr *connectmgr // 連接管理器
	listenmgr  *listenmgr  // 接聽管理器
	sessionmgr *sessionmgr // 會話管理器
}

// ConnectID 連接編號
type ConnectID = int64

// ListenID 接聽編號
type ListenID = int64

// AddConnectTCP 新增一個 TCP 連接器並立即嘗試連線
func (this *Netmgr) AddConnectTCP(ip, port string, timeout time.Duration, bind Bind, unbind Unbind, wrong Wrong) ConnectID {
	connect := NewTCPConnect(ip, port, timeout)
	connect.Connect(this.wrapperBind(bind), this.wrapperUnbind(unbind), wrong)
	return this.connectmgr.add(connect)
}

// DelConnect 刪除連接
func (this *Netmgr) DelConnect(connectID ConnectID) {
	this.connectmgr.del(connectID)
}

// GetConnect 取得連接
func (this *Netmgr) GetConnect(connectID ConnectID) Connecter {
	return this.connectmgr.get(connectID)
}

// AddListenTCP 新增一個 TCP 接聽器並立即開始監聽
func (this *Netmgr) AddListenTCP(ip, port string, bind Bind, unbind Unbind, wrong Wrong) ListenID {
	listen := NewTCPListen(ip, port)
	listen.Listen(this.wrapperBind(bind), this.wrapperUnbind(unbind), wrong)
	return this.listenmgr.add(listen)
}

// DelListen 刪除接聽
func (this *Netmgr) DelListen(listenID ListenID) {
	this.listenmgr.del(listenID)
}

// GetListen 取得接聽
func (this *Netmgr) GetListen(listenID ListenID) Listener {
	return this.listenmgr.get(listenID)
}

// Stop 停止 Netmgr 管理的所有資源
func (this *Netmgr) Stop() {
	this.connectmgr.clear()
	this.listenmgr.clear()
	this.sessionmgr.clear()
}

// Status 取得當前統計資料(connect 數量, listen 數量, session 數量)
func (this *Netmgr) Status() (connect, listen, session int) {
	return this.connectmgr.count(), this.listenmgr.count(), this.sessionmgr.count()
}

// wrapperBind 將外部傳入的 bind 包裝為 Netmgr 內部使用的初始化流程
func (this *Netmgr) wrapperBind(bind Bind) Bind {
	return func(session Sessioner) bool {
		if bind.Do(session) {
			this.sessionmgr.add(session)
			return true
		} // if

		return false
	}
}

// wrapperUnbind 將外部傳入的 unbind 包裝為 Netmgr 內部使用的釋放流程
func (this *Netmgr) wrapperUnbind(unbind Unbind) Unbind {
	return func(session Sessioner) {
		unbind.Do(session)
		this.sessionmgr.del(session)
	}
}

// newConnectmgr 建立連接管理器
func newConnectmgr() *connectmgr {
	return &connectmgr{
		data: map[ConnectID]Connecter{},
	}
}

// connectmgr 連接管理器
type connectmgr struct {
	connectID ConnectID               // 連接編號
	data      map[ConnectID]Connecter // 資料列表
	lock      sync.RWMutex            // 執行緒鎖
}

// add 新增連接
func (this *connectmgr) add(connect Connecter) ConnectID {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.connectID++
	this.data[this.connectID] = connect
	return this.connectID
}

// del 刪除連接
func (this *connectmgr) del(connectID ConnectID) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.data, connectID)
}

// clear 清除連接
func (this *connectmgr) clear() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.data = map[ConnectID]Connecter{}
}

// get 取得連接
func (this *connectmgr) get(connectID ConnectID) Connecter {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data[connectID]
}

// count 取得連接數量
func (this *connectmgr) count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.data)
}

// newListenmgr 建立接聽管理器
func newListenmgr() *listenmgr {
	return &listenmgr{
		data: map[ListenID]Listener{},
	}
}

// listenmgr 接聽管理器
type listenmgr struct {
	listenID ListenID              // 接聽編號
	data     map[ListenID]Listener // 資料列表
	lock     sync.RWMutex          // 執行緒鎖
}

// add 新增接聽
func (this *listenmgr) add(listen Listener) ListenID {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.listenID++
	this.data[this.listenID] = listen
	return this.listenID
}

// del 刪除接聽
func (this *listenmgr) del(listenID ListenID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if listen, ok := this.data[listenID]; ok {
		_ = listen.Stop()
		delete(this.data, listenID)
	} // if
}

// clear 清除接聽
func (this *listenmgr) clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.data {
		_ = itor.Stop()
	} // for

	this.data = map[ListenID]Listener{}
}

// get 取得接聽
func (this *listenmgr) get(listenID ListenID) Listener {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data[listenID]
}

// count 取得接聽數量
func (this *listenmgr) count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.data)
}

// newSessionmgr 建立會話管理器
func newSessionmgr() *sessionmgr {
	return &sessionmgr{
		data: hashset.New(),
	}
}

// sessionmgr 會話管理器
type sessionmgr struct {
	data *hashset.Set // 資料列表
	lock sync.RWMutex // 執行緒鎖
}

// add 新增會話
func (this *sessionmgr) add(session Sessioner) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.data.Add(session)
}

// del 刪除會話
func (this *sessionmgr) del(session Sessioner) {
	this.lock.Lock()
	defer this.lock.Unlock()
	session.Stop()
	this.data.Remove(session)
}

// clear 清除會話
func (this *sessionmgr) clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.data.Values() {
		itor.(Sessioner).Stop()
	} // for

	this.data = hashset.New()
}

// count 取得會話數量
func (this *sessionmgr) count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.data.Size()
}
