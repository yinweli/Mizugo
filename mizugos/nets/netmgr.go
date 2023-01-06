package nets

import (
	"fmt"
	"sync"
)

// NewNetmgr 建立網路管理器
func NewNetmgr() *Netmgr {
	return &Netmgr{
		listenmgr:  newListenmgr(),
		sessionmgr: newSessionmgr(),
	}
}

// Netmgr 網路管理器
type Netmgr struct {
	listenmgr  *listenmgr  // 接聽管理器
	sessionmgr *sessionmgr // 會話管理器
}

// Status 狀態資料
type Status struct {
	Listen  []string // 接聽列表
	Session int      // 會話數量
}

// AddConnect 新增連接
func (this *Netmgr) AddConnect(connecter Connecter, binder Binder) {
	connecter.Connect(func(session Sessioner, err error) {
		if err != nil {
			binder.Error(fmt.Errorf("netmgr connect: %v: %w", connecter.Address(), err))
			return
		} // if

		go session.Start(this.sessionmgr.add(session), binder)
	})
}

// AddListen 新增接聽
func (this *Netmgr) AddListen(listener Listener, binder Binder) ListenID {
	listener.Listen(func(session Sessioner, err error) {
		if err != nil {
			binder.Error(fmt.Errorf("netmgr listen: %v: %w", listener.Address(), err))
			return
		} // if

		go session.Start(this.sessionmgr.add(session), binder)
	})
	return this.listenmgr.add(listener)
}

// DelListen 刪除接聽
func (this *Netmgr) DelListen(listenID ListenID) {
	this.listenmgr.del(listenID)
}

// GetSession 取得會話
func (this *Netmgr) GetSession(sessionID SessionID) Sessioner {
	return this.sessionmgr.get(sessionID)
}

// StopSession 停止會話
func (this *Netmgr) StopSession(sessionID SessionID) {
	this.sessionmgr.del(sessionID)
}

// Stop 停止網路
func (this *Netmgr) Stop() {
	this.listenmgr.clear()
	this.sessionmgr.clear()
}

// Status 取得狀態資料
func (this *Netmgr) Status() *Status {
	return &Status{
		Listen:  this.listenmgr.address(),
		Session: this.sessionmgr.count(),
	}
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
	data     map[ListenID]Listener // 接聽列表
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

// address 取得接聽位址列表
func (this *listenmgr) address() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	result := []string{}

	for _, itor := range this.data {
		result = append(result, itor.Address())
	} // for

	return result
}

// newSessionmgr 建立會話管理器
func newSessionmgr() *sessionmgr {
	return &sessionmgr{
		data: map[SessionID]Sessioner{},
	}
}

// sessionmgr 會話管理器
type sessionmgr struct {
	sessionID SessionID               // 會話編號
	data      map[SessionID]Sessioner // 會話列表
	lock      sync.RWMutex            // 執行緒鎖
}

// add 新增會話
func (this *sessionmgr) add(session Sessioner) SessionID {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.sessionID++
	this.data[this.sessionID] = session
	return this.sessionID
}

// del 刪除會話
func (this *sessionmgr) del(sessionID SessionID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if session, ok := this.data[sessionID]; ok {
		session.Stop()
		delete(this.data, sessionID)
	} // if
}

// clear 清除會話
func (this *sessionmgr) clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.data {
		itor.Stop()
	} // for

	this.data = map[SessionID]Sessioner{}
}

// get 取得會話
func (this *sessionmgr) get(sessionID SessionID) Sessioner {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.data[sessionID]
}

// count 取得會話數量
func (this *sessionmgr) count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return len(this.data)
}
