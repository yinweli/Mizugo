package mizugos

import (
	"sync"

	"github.com/yinweli/Mizugo/mizugos/configs"
	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/labels"
	"github.com/yinweli/Mizugo/mizugos/logs"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/pools"
	"github.com/yinweli/Mizugo/mizugos/redmos"
)

// Start 啟動伺服器
//
// 啟動伺服器後按使用者的需要, 可以參考以下的寫法
//
//	defer func() {
//	    if cause := recover(); cause != nil {
//	        // 處理崩潰錯誤
//	    } // if
//	}()
//
//	ctx := ctxs.Root().WithCancel()
//	name := "伺服器名稱"
//	mizugos.Start() // 啟動伺服器
//
//	// 使用者自訂的初始化程序
//	// 如果有任何失敗, 執行 mizugos.Stop() 後退出
//
//	fmt.Printf("%v start\n", name)
//
//	for range ctx.Done() { // 進入無限迴圈直到執行 ctx.Cancel()
//	} // for
//
//	// 使用者自訂的結束程序
//	// 如果有任何失敗, 執行 mizugos.Stop() 後退出
//
//	mizugos.Stop() // 關閉伺服器
//	fmt.Printf("%v shutdown\n", name)
func Start() {
	server.lock.Lock()
	defer server.lock.Unlock()

	server.configmgr = configs.NewConfigmgr()
	server.metricsmgr = metrics.NewMetricsmgr()
	server.logmgr = logs.NewLogmgr()
	server.netmgr = nets.NewNetmgr()
	server.redmomgr = redmos.NewRedmomgr()
	server.entitymgr = entitys.NewEntitymgr()
	server.labelmgr = labels.NewLabelmgr()
	server.poolmgr = pools.DefaultPool // 執行緒池管理器直接用預設的
}

// Stop 關閉伺服器
func Stop() {
	server.lock.RLock()
	defer server.lock.RUnlock()

	server.configmgr = nil
	server.metricsmgr = nil
	server.logmgr = nil
	server.netmgr = nil
	server.redmomgr = nil
	server.entitymgr = nil
	server.labelmgr = nil
	server.poolmgr = nil
	ctxs.Root().Cancel() // 關閉由contexts.Ctx()衍生出來的執行緒, 避免goroutine洩漏
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

// Redmomgr 取得資料庫管理器
func Redmomgr() *redmos.Redmomgr {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.redmomgr
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
func Debug(name, label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Debug(name, label)
}

// Info 記錄一般訊息
func Info(name, label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Info(name, label)
}

// Warn 記錄警告訊息
func Warn(name, label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Warn(name, label)
}

// Error 記錄錯誤訊息
func Error(name, label string) logs.Stream {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.logmgr.Error(name, label)
}

// server 伺服器資料
var server struct {
	configmgr  *configs.Configmgr  // 配置管理器
	metricsmgr *metrics.Metricsmgr // 統計管理器
	logmgr     *logs.Logmgr        // 日誌管理器
	netmgr     *nets.Netmgr        // 網路管理器
	redmomgr   *redmos.Redmomgr    // 資料庫管理器
	entitymgr  *entitys.Entitymgr  // 實體管理器
	labelmgr   *labels.Labelmgr    // 標籤管理器
	poolmgr    *pools.Poolmgr      // 執行緒池管理器
	lock       sync.RWMutex        // 執行緒鎖
}
