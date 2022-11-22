package nets

import (
	"fmt"
	"net"
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
	lister  net.Listener // 接聽物件
}

// Start 啟動接聽, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止接聽
func (this *TCPListen) Start(complete Complete) {
	lister, err := net.Listen("tcp", this.address)

	if err != nil {
		complete(nil, fmt.Errorf("tcp listen start: %w", err))
		return
	} // if

	this.lister = lister

	for {
		conn, err := this.lister.Accept()

		if err != nil {
			complete(nil, fmt.Errorf("tcp listen start: %w", err))
			return
		} // if

		complete(NewTCPSession(conn), nil)
	} // for
}

// Stop 停止接聽
func (this *TCPListen) Stop() error {
	if this.lister != nil {
		if err := this.lister.Close(); err != nil {
			return fmt.Errorf("tcp listen stop: %w", err)
		} // if
	} // if

	return nil
}

// Address 取得位址
func (this *TCPListen) Address() string {
	return this.address
}
