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
	listenmgr  *listenmgr  // 監聽管理器
	sessionmgr *sessionmgr // 會話管理器
}

// Status 狀態資料
type Status struct {
	Listen  []string // 監聽列表
	session int      // 會話數量
}

// Create 建立處理函式類型
type Create func() (coder Coder, reactor Reactor)

// Failed 錯誤處理函式類型
type Failed func(err error)

// AddConnect 新增連接
func (this *Netmgr) AddConnect(connecter Connecter, create Create, failed Failed) {
	go connecter.Connect(func(session Sessioner, err error) {
		if err != nil {
			failed(fmt.Errorf("netmgr connect: %s: %w", connecter.Address(), err))
			return
		} // if

		sessionID := this.sessionmgr.add(session)
		coder, reactor := create()
		go session.Start(sessionID, coder, reactor)
	})
}

// AddListen 新增監聽
func (this *Netmgr) AddListen(listener Listener, create Create, failed Failed) {
	this.listenmgr.add(listener)
	go listener.Listen(func(session Sessioner, err error) {
		if err != nil {
			failed(fmt.Errorf("netmgr listen: %s: %w", listener.Address(), err))
			return
		} // if

		sessionID := this.sessionmgr.add(session)
		coder, reactor := create()
		go session.Start(sessionID, coder, reactor)
	})
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
		session: this.sessionmgr.count(),
	}
}

// newListenmgr 建立監聽管理器
func newListenmgr() *listenmgr {
	return &listenmgr{}
}

// listenmgr 監聽管理器
type listenmgr struct {
	data []Listener // 監聽列表
	lock sync.Mutex // 執行緒鎖
}

// add 新增監聽
func (this *listenmgr) add(listener Listener) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data = append(this.data, listener)
}

// clear 清除監聽
func (this *listenmgr) clear() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, itor := range this.data {
		_ = itor.Stop()
	} // for
}

// address 取得監聽位址列表
func (this *listenmgr) address() []string {
	this.lock.Lock()
	defer this.lock.Unlock()

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
	lock      sync.Mutex              // 執行緒鎖
}

// add 新增會話
func (this *sessionmgr) add(sessioner Sessioner) SessionID {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.sessionID++
	this.data[this.sessionID] = sessioner
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
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.data[sessionID]
}

// count 取得會話數量
func (this *sessionmgr) count() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return len(this.data)
}
