package nets

import (
	"fmt"
	"net"
	"time"
)

// NewTCPConnect 建立tcp連接器
func NewTCPConnect(ip, port string, timeout time.Duration) *TCPConnect {
	return &TCPConnect{
		address: net.JoinHostPort(ip, port),
		timeout: timeout,
	}
}

// TCPConnect tcp連接器
type TCPConnect struct {
	address string        // 位址字串
	timeout time.Duration // 逾時時間
}

// Connect 啟動連接
func (this *TCPConnect) Connect(completer Completer) {
	go func() {
		conn, err := net.DialTimeout("tcp", this.address, this.timeout)

		if err != nil {
			completer.Complete(nil, fmt.Errorf("tcp connect: %v: %w", this.address, err))
			return
		} // if

		completer.Complete(NewTCPSession(conn), nil)
	}()
}

// Address 取得位址
func (this *TCPConnect) Address() string {
	return this.address
}
