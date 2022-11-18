package nets

import (
	"net"
)

// Connector 連接介面
type Connector interface {
	// Start 啟動連接, 若不是使用多執行緒啟動, 則可能被阻塞在這裡直到連接成功
	Start(complete Complete)

	// Address 取得位址
	Address() (addr net.Addr, err error)
}

// Listener 接聽介面
type Listener interface {
	// Start 啟動接聽, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止接聽
	Start(complete Complete)

	// Stop 停止接聽
	Stop() error

	// Address 取得位址
	Address() (addr net.Addr, err error)
}

// Sessioner 會話介面
type Sessioner interface {
	// Start 啟動會話
	Start(sessionID SessionID, receive Receive, inform Inform)

	// StopImmed 立即停止會話
	StopImmed()

	// StopWait 等待停止會話
	StopWait()

	// Send 傳送封包
	Send(packet []byte)

	// SessionID 取得會話編號
	SessionID() SessionID

	// RemoteAddr 取得遠端位址
	RemoteAddr() net.Addr

	// LocalAddr 取得本地位址
	LocalAddr() net.Addr
}

// Complete 完成函式類型
type Complete func(session Sessioner, err error)

// Receive 接收函式類型
type Receive func(packet []byte)

// Inform 通知函式類型
type Inform func(err error)

// SessionID 會話編號
type SessionID = int64
