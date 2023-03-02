package mizugos

import (
	"fmt"
	"sync"

	"github.com/yinweli/Mizugo/mizugos/configs"
	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/depots"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/labels"
	"github.com/yinweli/Mizugo/mizugos/logs"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/pools"
)

// TODO: 遊戲伺服器間的網絡若出現問題 (DB設備斷線、應用程序斷線)，會自動嘗試重新連線嗎？
// TODO: 遊戲伺服器應用程式在發生問題或故障時有自動修復的功能嗎？

// Start 啟動伺服器, 用於啟動mizugo伺服器, 需要指定 Initialize 與 Finalize 來執行使用者的初始化與結束處理;
// 啟動伺服器執行的順序為
//   - 設置內部成員
//   - 建立各管理器, 注意! 這裡只建立而沒有對各管理器進行初始化
//   - 執行 Initialize
//   - 進入無限迴圈, 直到關閉伺服器(呼叫 Stop 後)
//   - 執行 Finalize
//   - 釋放內部成員
//   - 釋放各管理器
//   - 最後呼叫 ctxs.Root().Cancel() 讓由contexts.Ctx()衍生出來的執行緒最後都能被終止, 避免goroutine洩漏
//
// 為了讓程式持續執行, 此函式不能用執行緒執行, 也請不要執行此函式兩次
func Start(name string, initialize Initialize, finalize Finalize) {
	server.lock.Lock()
	server.name = name
	server.ctx = ctxs.Root().WithCancel()
	server.configmgr = configs.NewConfigmgr()
	server.metricsmgr = metrics.NewMetricsmgr()
	server.logmgr = logs.NewLogmgr()
	server.netmgr = nets.NewNetmgr()
	server.depotmgr = depots.NewDepotmgr()
	server.entitymgr = entitys.NewEntitymgr()
	server.labelmgr = labels.NewLabelmgr()
	server.poolmgr = pools.DefaultPool // 執行緒池管理器直接用預設的
	server.lock.Unlock()

	fmt.Printf("%v initialize\n", name)

	if err := initialize.Do(); err != nil {
		fmt.Println(fmt.Errorf("%v initialize: %w", name, err))
		goto Finalize
	} // if

	fmt.Printf("%v start\n", name)

	// 進行等待, 直到關閉伺服器
	for range server.ctx.Done() {
	} // for

	fmt.Printf("%v shutdown\n", name)
	finalize.Do()

Finalize: // 結束處理
	fmt.Printf("%v finalize\n", name)
	server.lock.Lock()
	server.name = ""
	server.ctx.Cancel()
	server.configmgr = nil
	server.metricsmgr = nil
	server.logmgr = nil
	server.netmgr = nil
	server.depotmgr = nil
	server.entitymgr = nil
	server.labelmgr = nil
	server.poolmgr = nil
	server.lock.Unlock()
	ctxs.Root().Cancel() // 關閉伺服器, 並且保證由contexts.Ctx()衍生出來的執行緒最後都能被終止, 避免goroutine洩漏
}

// Stop 關閉伺服器
func Stop() {
	server.lock.RLock()
	defer server.lock.RUnlock()

	server.ctx.Cancel()
}

// Name 取得伺服器名稱
func Name() string {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.name
}

// ===== 管理器功能 =====

// Configmgr 取得配置管理器
func Configmgr() *configs.Configmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.configmgr
}

// Metricsmgr 統計管理器
func Metricsmgr() *metrics.Metricsmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.metricsmgr
}

// Logmgr 日誌管理器
func Logmgr() *logs.Logmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr
}

// Netmgr 取得網路管理器
func Netmgr() *nets.Netmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.netmgr
}

// Depotmgr 取得資料庫管理器
func Depotmgr() *depots.Depotmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.depotmgr
}

// Entitymgr 實體管理器
func Entitymgr() *entitys.Entitymgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.entitymgr
}

// Labelmgr 標籤管理器
func Labelmgr() *labels.Labelmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.labelmgr
}

// Poolmgr 執行緒池管理器
func Poolmgr() *pools.Poolmgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.poolmgr
}

// ===== 日誌功能 =====

// Debug 記錄除錯訊息
func Debug(label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Debug(label)
}

// Info 記錄一般訊息
func Info(label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Info(label)
}

// Warn 記錄警告訊息
func Warn(label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Warn(label)
}

// Error 記錄錯誤訊息
func Error(label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Error(label)
}

// ===== 其他定義 =====

// Initialize 初始化處理函式類型
type Initialize func() error

// Do 執行處理
func (this Initialize) Do() error {
	if this != nil {
		return this()
	} // if

	return nil
}

// Finalize 結束處理函式類型
type Finalize func()

// Do 執行處理
func (this Finalize) Do() {
	if this != nil {
		this()
	} // if
}

// server 伺服器資料
var server struct {
	name       string              // 伺服器名稱
	ctx        ctxs.Ctx            // ctx物件
	configmgr  *configs.Configmgr  // 配置管理器
	metricsmgr *metrics.Metricsmgr // 統計管理器
	logmgr     *logs.Logmgr        // 日誌管理器
	netmgr     *nets.Netmgr        // 網路管理器
	depotmgr   *depots.Depotmgr    // 資料庫管理器
	entitymgr  *entitys.Entitymgr  // 實體管理器
	labelmgr   *labels.Labelmgr    // 標籤管理器
	poolmgr    *pools.Poolmgr      // 執行緒池管理器
	lock       sync.RWMutex        // 執行緒鎖
}
