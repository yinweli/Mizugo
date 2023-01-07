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
type Netmgr struct {
	connectmgr *connectmgr // 連接管理器
	listenmgr  *listenmgr  // 接聽管理器
	sessionmgr *sessionmgr // 會話管理器
}

// Status 狀態資料
type Status struct {
	Connect int // 連接數量
	Listen  int // 接聽數量
	Session int // 會話數量
}

// ConnectID 連接編號
type ConnectID = int64

// ListenID 接聽編號
type ListenID = int64

// AddConnectTCP 新增連接(TCP)
func (this *Netmgr) AddConnectTCP(ip, port string, timeout time.Duration, inform Inform) ConnectID {
	connect := NewTCPConnect(ip, port, timeout)
	connect.Connect(this.wrapperInform(inform))
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

// AddListenTCP 新增接聽(TCP)
func (this *Netmgr) AddListenTCP(ip, port string, inform Inform) ListenID {
	listen := NewTCPListen(ip, port)
	listen.Listen(this.wrapperInform(inform))
	return this.listenmgr.add(listen)
}

// DelListen 刪除連接
func (this *Netmgr) DelListen(listenID ListenID) {
	this.listenmgr.del(listenID)
}

// GetListen 取得連接
func (this *Netmgr) GetListen(listenID ListenID) Listener {
	return this.listenmgr.get(listenID)
}

// Stop 停止網路
func (this *Netmgr) Stop() {
	this.connectmgr.clear()
	this.listenmgr.clear()
	this.sessionmgr.clear()
}

// Status 取得狀態資料
func (this *Netmgr) Status() *Status {
	return &Status{
		Connect: this.connectmgr.count(),
		Listen:  this.listenmgr.count(),
		Session: this.sessionmgr.count(),
	}
}

// wrapperInform 包裝通知物件
func (this *Netmgr) wrapperInform(inform Inform) Inform {
	return Inform{
		Error: inform.Error,
		Bind: func(session Sessioner) {
			this.sessionmgr.add(session)
			inform.Bind.Do(session)
		},
		Unbind: func(session Sessioner) {
			inform.Unbind.Do(session)
			this.sessionmgr.del(session)
		},
		Encode:    inform.Encode,
		Decode:    inform.Decode,
		Receive:   inform.Receive,
		AfterSend: inform.AfterSend,
		AfterRecv: inform.AfterRecv,
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

	if _, ok := this.data[connectID]; ok {
		delete(this.data, connectID)
	} // if
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
