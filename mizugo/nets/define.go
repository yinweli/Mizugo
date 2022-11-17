package nets

import (
	"net"
)

// Connector 連接介面
type Connector interface {
	// Start 啟動連接
	Start(complete Complete)

	// Address 取得位址
	Address() net.Addr
}

// Listener 接聽介面
type Listener interface {
	// Start 啟動接聽
	Start(complete Complete)

	// Stop 停止接聽
	Stop()

	// Address 取得位址
	Address() net.Addr
}

// Sessioner 會話介面
type Sessioner interface {
	// Initialize 初始化處理
	Initialize(receive Receive, inform Inform)

	// CloseImmed 立即關閉
	CloseImmed()

	// CloseWait 等待關閉
	CloseWait()

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
