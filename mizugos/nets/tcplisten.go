package nets

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/yinweli/Mizugo/mizugos/pools"
)

// NewTCPListen 建立TCP接聽器
func NewTCPListen(ip, port string) *TCPListen {
	return &TCPListen{
		address: net.JoinHostPort(ip, port),
	}
}

// TCPListen TCP接聽器, 負責用TCP協議建立接聽, 並等待客戶端連接以取得會話物件
type TCPListen struct {
	address string       // 位址字串
	listen  net.Listener // 接聽物件
	closed  atomic.Bool  // 關閉旗標
}

// Listen 啟動接聽
func (this *TCPListen) Listen(bind Bind, unbind Unbind, wrong Wrong) {
	listen, err := net.Listen("tcp", this.address)

	if err != nil {
		wrong.Do(fmt.Errorf("tcp listen: %v: %w", this.address, err))
		return
	} // if

	this.listen = listen

	pools.DefaultPool.Submit(func() {
		// 由於listen.Accept的執行方式, 所以不需要用context方式監控終止方式

		for {
			conn, err := this.listen.Accept()

			if err != nil {
				if this.closed.Load() {
					return // 停止接聽, 這不算是錯誤, 但要結束接聽器了
				} else {
					wrong.Do(fmt.Errorf("tcp listen: %v: %w", this.address, err))
					continue // 這次連接出了問題, 但我們還是繼續接聽
				} // if
			} // if

			session := NewTCPSession(conn)
			session.Start(bind, unbind, wrong)
		} // for
	})
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
