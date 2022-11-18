package nets

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// NewTCPConnect 建立tcp連接器
func NewTCPConnect(ip string, port int, timeout time.Duration) *TCPConnect {
	return &TCPConnect{
		ip:      ip,
		port:    port,
		timeout: timeout,
	}
}

// TCPConnect tcp連接器
type TCPConnect struct {
	ip      string        // 位址字串
	port    int           // 埠號
	timeout time.Duration // 逾時時間
}

// Start 啟動連接, 若不是使用多執行緒啟動, 則可能被阻塞在這裡直到連接成功
func (this *TCPConnect) Start(complete Complete) {
	addr, err := this.Address()

	if err != nil {
		complete(nil, fmt.Errorf("tcp connect start: %w", err))
		return
	} // if

	conn, err := net.DialTimeout(addr.Network(), addr.String(), this.timeout)

	if err != nil {
		complete(nil, fmt.Errorf("tcp connect start: %w", err))
		return
	} // if

	complete(NewTCPSession(conn), nil)
}

// Address 取得位址
func (this *TCPConnect) Address() (addr net.Addr, err error) {
	addr, err = net.ResolveTCPAddr("", this.ip+":"+strconv.Itoa(this.port)) // 第一個參數留空, 會自動幫我填寫正確的tcp網路名稱

	if err != nil {
		return nil, fmt.Errorf("tcp connect address: %w", err)
	} // if

	return addr, nil
}
