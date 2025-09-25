package nets

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/yinweli/Mizugo/v2/mizugos/pools"
)

// NewTCPListen 建立TCP接聽器
func NewTCPListen(ip, port string) *TCPListen {
	return &TCPListen{
		address: net.JoinHostPort(ip, port),
	}
}

// TCPListen TCP 接聽器
//
// 負責以 TCP 協議建立接聽並持續接受連線, 對每個成功的連線產生 TCPSession
//
// 特性:
//   - 僅負責「接聽與接受連線」; 連線建立後的讀寫生命週期交由 TCPSession 管理
//   - 支援 Stop: 會設置 closed 旗標並關閉底層 listener, 使 Accept 迴圈可正常退出
//   - 失敗或錯誤僅透過 wrong 回報, 不會以函式回傳值通知
type TCPListen struct {
	address string       // 接聽位址字串
	listen  net.Listener // 接聽物件
	closed  atomic.Bool  // 關閉旗標
}

// Listen 啟動接聽
//   - 會將任務提交給 pools.DefaultPool 執行, 因此 Listen 本身非阻塞
//   - 內部以 net.Listen 嘗試接聽, 失敗透過 wrong 回報
//   - 每次 Accept 成功會生成一個新的 TCPSession, 並呼叫 TCPSession.Start(bind, unbind) 啟動會話
//   - 若 Accept 發生錯誤: 若 closed=true 視為正常停止, 直接結束迴圈; 否則透過 wrong 回報並繼續嘗試接受下一個連線
//
// 回呼函式:
//   - bind: 會話綁定, 用於初始化
//   - unbind: 解綁回呼, 用於釋放資源
//   - wrong: 錯誤回呼, 用於集中處理接聽錯誤; 若為 nil 則錯誤會被忽略
func (this *TCPListen) Listen(bind Bind, unbind Unbind, wrong Wrong) {
	listen, err := net.Listen("tcp", this.address)

	if err != nil {
		wrong.Do(fmt.Errorf("tcp listen: %v: %w", this.address, err))
		return
	} // if

	this.listen = listen
	this.closed.Store(false)

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
			session.Start(bind, unbind)
		} // for
	})
}

// Stop 停止接聽
//   - 若 listener 存在: 先將 closed 設為 true, 再關閉 listener
//   - 關閉 listener 會使 Accept 立即返回錯誤, Accept 迴圈據此判斷 closed 後正常退出
func (this *TCPListen) Stop() error {
	if this.listen != nil {
		this.closed.Store(true)

		if err := this.listen.Close(); err != nil {
			return fmt.Errorf("tcp listen stop: %v, %w", this.address, err)
		} // if
	} // if

	return nil
}

// Address 取得接聽位址字串
func (this *TCPListen) Address() string {
	return this.address
}
