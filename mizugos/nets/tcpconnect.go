package nets

import (
	"fmt"
	"net"
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos/pools"
)

// NewTCPConnect 建立TCP連接器
func NewTCPConnect(ip, port string, timeout time.Duration) *TCPConnect {
	return &TCPConnect{
		address: net.JoinHostPort(ip, port),
		timeout: timeout,
	}
}

// TCPConnect TCP 連接器
//
// 負責以 TCP 協議建立連線並產生 TCPSession
//
// 特性：
//   - 僅負責「發起連線」; 連線成功後的讀寫生命週期交由 TCPSession 管理
//   - 本物件本身沒有 Stop / Close, 呼叫端不需額外回收
//   - 失敗或錯誤僅透過 wrong 回報, 不會以函式回傳值通知
type TCPConnect struct {
	address string        // 連接位址字串
	timeout time.Duration // 逾時時間
}

// Connect 啟動連線
//   - 會將任務提交給 pools.DefaultPool 執行, 因此 Connect 本身非阻塞
//   - 內部以 net.DialTimeout 嘗試連線, 失敗透過 wrong 回報
//   - 成功建立連線後, 會生成一個新的 TCPSession, 並呼叫 TCPSession.Start(bind, unbind) 啟動會話
//   - 連線完成或失敗後任務即結束
//
// 回呼函式:
//   - bind: 會話綁定, 用於初始化
//   - unbind: 解綁回呼, 用於釋放資源
//   - wrong: 錯誤回呼, 用於集中處理連線錯誤; 若為 nil 則錯誤會被忽略
func (this *TCPConnect) Connect(bind Bind, unbind Unbind, wrong Wrong) {
	pools.DefaultPool.Submit(func() {
		// 由於連接完成/失敗後就直接結束, 所以不需要用context方式監控終止方式

		conn, err := net.DialTimeout("tcp", this.address, this.timeout)

		if err != nil {
			wrong.Do(fmt.Errorf("tcp connect: %v: %w", this.address, err))
			return
		} // if

		session := NewTCPSession(conn)
		session.Start(bind, unbind)
	})
}

// Address 取得連接位址字串
func (this *TCPConnect) Address() string {
	return this.address
}
