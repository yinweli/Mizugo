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

// Listen 啟動接聽
func (this *TCPListen) Listen(done Done) {
	listen, err := net.Listen("tcp", this.address)

	if err != nil {
		done(nil, fmt.Errorf("tcp listen: %v: %w", this.address, err))
		return
	} // if

	this.listen = listen

	go func() {
		for {
			conn, err := this.listen.Accept()

			if err != nil {
				if this.closed.Load() {
					return // 停止接聽, 這不算是錯誤, 但要結束接聽器了
				} else {
					done(nil, fmt.Errorf("tcp listen: %v: %w", this.address, err))
					continue // 這次連接出了問題, 但我們還是繼續接聽
				} // if
			} // if

			done(NewTCPSession(conn), nil)
		} // for
	}()
}

// Stop 停止接聽
func (this *TCPListen) Stop() error {
	if this.listen != nil {
		this.closed.Store(true)

		if err := this.listen.Close(); err != nil {
			return fmt.Errorf("tcp listen stop: %v, %w", this.address, err)
		} // if
	} // if

	return nil
}

// Address 取得位址
func (this *TCPListen) Address() string {
	return this.address
}
