package nets

import (
	"net"
)

// NewTCPSession 建立tcp會話器
func NewTCPSession(conn net.Conn) *TCPSession {
	return &TCPSession{
		conn: conn,
	}
}

// TCPSession tcp會話器
type TCPSession struct {
	conn      net.Conn  // 連線物件
	sessionID SessionID // 會話編號
	receive   Receive   // 接收函式
	inform    Inform    // 通知函式
}

// Start 啟動會話
func (this *TCPSession) Start(sessionID SessionID, receive Receive, inform Inform) {
	this.sessionID = sessionID
	this.receive = receive
	this.inform = inform
}

// StopImmed 立即停止會話
func (this *TCPSession) StopImmed() {

}

// StopWait 等待停止會話
func (this *TCPSession) StopWait() {

}

// Send 傳送封包
func (this *TCPSession) Send(packet []byte) {

}

// SessionID 取得會話編號
func (this *TCPSession) SessionID() SessionID {
	return this.sessionID
}

// RemoteAddr 取得遠端位址
func (this *TCPSession) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

// LocalAddr 取得本地位址
func (this *TCPSession) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}
