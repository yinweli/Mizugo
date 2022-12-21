package nets

import (
	"net"
)

// Connecter 連接介面
type Connecter interface {
	// Connect 啟動連接
	Connect(done Done)

	// Address 取得位址
	Address() string
}

// Listener 接聽介面
type Listener interface {
	// Listen 啟動接聽
	Listen(done Done)

	// Stop 停止接聽
	Stop() error

	// Address 取得位址
	Address() string
}

// Sessioner 會話介面
type Sessioner interface {
	// Start 啟動會話, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止會話; 當由連接器/接聽器獲得會話器之後, 需要啟動會話才可以傳送或接收封包
	Start(sessionID SessionID, bind Bind)

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

// Done 完成會話函式類型
type Done func(session Sessioner, err error)

// Bind 綁定處理函式類型
type Bind func(session Sessioner) *Notice

// Unbind 解綁處理函式類型
type Unbind func()

// Encode 封包編碼處理函式類型, 用在傳送封包時
type Encode func(message any) (packet []byte, err error)

// Decode 封包解碼處理函式類型, 用在接收封包時
type Decode func(packet []byte) (message any, err error)

// Receive 接收封包處理函式類型
type Receive func(message any) error

// Wrong 錯誤處理函式類型
type Wrong func(err error)

// SessionID 會話編號
type SessionID = int64

// Notice 通知資料
type Notice struct {
	Unbind  // 解綁處理函式
	Encode  // 封包編碼處理函式
	Decode  // 封包解碼處理函式
	Receive // 接收封包處理函式
	Wrong   // 錯誤處理函式
}
