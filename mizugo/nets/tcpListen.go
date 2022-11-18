package nets

import (
	"fmt"
	"net"
	"strconv"
)

// NewTCPListen 建立tcp連接器
func NewTCPListen(ip string, port int) *TCPListen {
	return &TCPListen{
		ip:   ip,
		port: port,
	}
}

// TCPListen tcp接聽器
type TCPListen struct {
	ip     string       // 位址字串
	port   int          // 埠號
	lister net.Listener // 接聽物件
}

// Start 啟動接聽, 若不是使用多執行緒啟動, 則一定被阻塞在這裡直到停止接聽
func (this *TCPListen) Start(complete Complete) {
	addr, err := this.Address()

	if err != nil {
		complete(nil, fmt.Errorf("tcp listen start: %w", err))
		return
	} // if

	lister, err := net.Listen(addr.Network(), addr.String())

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
func (this *TCPListen) Address() (addr net.Addr, err error) {
	addr, err = net.ResolveTCPAddr("", this.ip+":"+strconv.Itoa(this.port)) // 第一個參數留空, 會自動幫我填寫正確的tcp網路名稱

	if err != nil {
		return nil, fmt.Errorf("tcp listen address: %w", err)
	} // if

	return addr, nil
}
