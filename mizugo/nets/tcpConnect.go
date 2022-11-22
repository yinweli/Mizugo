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

// Start 啟動連接, 若不是使用多執行緒啟動, 則可能被阻塞在這裡直到連接成功
func (this *TCPConnect) Start(complete Complete) {
	conn, err := net.DialTimeout("tcp", this.address, this.timeout)

	if err != nil {
		complete(nil, fmt.Errorf("tcp connect start: %w", err))
		return
	} // if

	complete(NewTCPSession(conn), nil)
}

// Address 取得位址
func (this *TCPConnect) Address() string {
	return this.address
}
