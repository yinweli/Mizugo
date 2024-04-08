package nets

import (
	"fmt"
	"net"
	"time"

	"github.com/yinweli/Mizugo/mizugo/pools"
)

// NewTCPConnect 建立TCP連接器
func NewTCPConnect(ip, port string, timeout time.Duration) *TCPConnect {
	return &TCPConnect{
		address: net.JoinHostPort(ip, port),
		timeout: timeout,
	}
}

// TCPConnect TCP連接器, 負責用TCP協議建立連接以取得會話物件
type TCPConnect struct {
	address string        // 位址字串
	timeout time.Duration // 超時時間
}

// Connect 啟動連接
func (this *TCPConnect) Connect(bind Bind, unbind Unbind, wrong Wrong) {
	pools.DefaultPool.Submit(func() {
		// 由於連接完成/失敗後就直接結束, 所以不需要用context方式監控終止方式

		conn, err := net.DialTimeout("tcp", this.address, this.timeout)

		if err != nil {
			wrong.Do(fmt.Errorf("tcp connect: %v: %w", this.address, err))
			return
		} // if

		session := NewTCPSession(conn)
		session.Start(bind, unbind, wrong)
	})
}

// Address 取得位址
func (this *TCPConnect) Address() string {
	return this.address
}
