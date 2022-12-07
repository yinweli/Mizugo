package nets

import (
	"fmt"
	"net"
	"sync/atomic"
)

// NewTCPListen 建立tcp連接器
func NewTCPListen(ip, port string) *TCPListen {
	return &TCPListen{
		address: net.JoinHostPort(ip, port),
	}
}

// TCPListen tcp接聽器
type TCPListen struct {
	address string       // 位址字串
	listen  net.Listener // 接聽物件
	closed  atomic.Bool  // 關閉旗標
}

// Start 啟動接聽, 若不是使用多執行緒啟動, 則會被阻塞在這裡直到停止接聽
func (this *TCPListen) Start(complete Complete) {
	listen, err := net.Listen("tcp", this.address)

	if err != nil {
		complete(nil, fmt.Errorf("tcp listen start: %s: %w", this.address, err))
		return
	} // if

	this.listen = listen

	for {
		conn, err := this.listen.Accept()

		if err != nil {
			if this.closed.Load() == false { // 停止接聽前的錯誤才算是錯誤
				complete(nil, fmt.Errorf("tcp listen start: %s: %w", this.address, err))
			} // if

			return
		} // if

		complete(NewTCPSession(conn), nil)
	} // for
}

// Stop 停止接聽
func (this *TCPListen) Stop() error {
	if this.listen != nil {
		this.closed.Store(true)

		if err := this.listen.Close(); err != nil {
			return fmt.Errorf("tcp listen stop: %s, %w", this.address, err)
		} // if
	} // if

	return nil
}

// Address 取得位址
func (this *TCPListen) Address() string {
	return this.address
}
