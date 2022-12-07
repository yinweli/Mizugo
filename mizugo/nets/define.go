package nets

import (
	"net"
)

// Connecter 連接介面
type Connecter interface {
	// Start 啟動連接, 若不是使用多執行緒啟動, 則可能被阻塞在這裡直到連接完成
	Start(complete Complete)

	// Address 取得位址
	Address() string
}

// Listener 接聽介面
type Listener interface {
	// Start 啟動接聽, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止接聽
	Start(complete Complete)

	// Stop 停止接聽
	Stop() error

	// Address 取得位址
	Address() string
}

// Sessioner 會話介面
type Sessioner interface {
	// Start 啟動會話, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止會話; 當由連接器/監聽器獲得會話器之後, 需要啟動會話才可以傳送或接收封包
	Start(sessionID SessionID, coder Coder, reactor Reactor)

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

// Coder 編碼介面
type Coder interface {
	// Encode 封包編碼, 用在傳送封包時
	Encode(message any) (packet []byte, err error)

	// Decode 封包解碼, 用在接收封包時
	Decode(packet []byte) (message any, err error)
}

// Reactor 處理介面
type Reactor interface {
	// Active 啟動處理
	Active()

	// Inactive 結束處理
	Inactive()

	// Error 錯誤處理
	Error(err error)

	// Receive 接收處理
	Receive(message any) error
}

// Complete 完成函式類型
type Complete func(session Sessioner, err error)

// SessionID 會話編號
type SessionID = int64
