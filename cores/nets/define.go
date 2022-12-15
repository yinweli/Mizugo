package nets

import (
	"net"
)

// Connecter 連接介面
type Connecter interface {
	// Connect 啟動連接
	Connect(completer Completer)

	// Address 取得位址
	Address() string
}

// Listener 接聽介面
type Listener interface {
	// Listen 啟動接聽
	Listen(completer Completer)

	// Stop 停止接聽
	Stop() error

	// Address 取得位址
	Address() string
}

// Completer 完成會話介面
type Completer interface {
	// Complete 完成會話
	Complete(session Sessioner, err error)
}

// Complete 完成會話函式類型
type Complete func(session Sessioner, err error)

// Complete 完成會話
func (this Complete) Complete(session Sessioner, err error) {
	this(session, err)
}

// Sessioner 會話介面
type Sessioner interface {
	// Start 啟動會話, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止會話; 當由連接器/接聽器獲得會話器之後, 需要啟動會話才可以傳送或接收封包
	Start(sessionID SessionID, binder Binder)

	// Stop 停止會話, 不會等待會話內部循環結束
	Stop()

	// StopWait 停止會話, 會等待會話內部循環結束
	StopWait()

	// Send 傳送封包
	Send(message any)

	// SessionID 取得會話編號
	SessionID() SessionID

	// RemoteAddr 取得遠端位址
	RemoteAddr() net.Addr

	// LocalAddr 取得本地位址
	LocalAddr() net.Addr
}

// Binder 綁定介面
type Binder interface {
	// Bind 綁定處理
	Bind(session Sessioner) (reactor Reactor, unbinder Unbinder)

	// Error 錯誤處理
	Error(err error)
}

// Reactor 反應介面
type Reactor interface {
	// Encode 封包編碼, 用在傳送封包時
	Encode(message any) (packet []byte, err error)

	// Decode 封包解碼, 用在接收封包時
	Decode(packet []byte) (message any, err error)

	// Receive 接收處理
	Receive(message any) error

	// Error 錯誤處理
	Error(err error)
}

// Unbinder 解綁介面
type Unbinder interface {
	// Unbind 解綁處理
	Unbind()
}

// Unbind 解綁處理函式類型
type Unbind func()

// Unbind 解綁處理
func (this Unbind) Unbind() {
	this()
}

// SessionID 會話編號
type SessionID = int64