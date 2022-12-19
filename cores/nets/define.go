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
	Bind(session Sessioner) (unbinder Unbinder, encoder Encoder, receiver Receiver)

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

// TODO: 是否要把Completer + Binder + Unbinder + Encode + Decode + Receive + Error改回函式指標?
// TODO: Encode + Decode + Receive可能會是一組的封包前置處理程序, 後面應當要接到包括封包編號與封包物件的封包處理器
// TODO: 封包處理器可能可以固定介面了(封包編號+封包物件)
// TODO: Encode + Decode實作時要考慮可用開關切換編解碼流程, 例如首次封包不須加密解密之類(因為密鑰還沒送到/收到)

// Encoder 編碼介面
type Encoder interface {
	// Encode 封包編碼, 用在傳送封包時
	Encode(message any) (packet []byte, err error)

	// Decode 封包解碼, 用在接收封包時
	Decode(packet []byte) (message any, err error)
}

// Receiver 接收介面
type Receiver interface {
	// Receive 接收處理
	Receive(message any) error

	// Error 錯誤處理
	Error(err error)
}

// SessionID 會話編號
type SessionID = int64
