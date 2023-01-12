package mizugos

import (
	"context"
	"fmt"
	"sync"

	"github.com/yinweli/Mizugo/mizugos/configs"
	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/labels"
	"github.com/yinweli/Mizugo/mizugos/logs"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/mizugos/nets"
)

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

// Start 啟動伺服器, 為了讓程式持續執行, 此函式不能用執行緒執行; 也請不要執行此函式兩次
func Start(name string, initialize Initialize, finalize Finalize) {
	server.lock.Lock()
	server.name = name
	server.ctx, server.cancel = context.WithCancel(contexts.Ctx())
	server.configmgr = configs.NewConfigmgr()
	server.netmgr = nets.NewNetmgr()
	server.entitymgr = entitys.NewEntitymgr()
	server.labelmgr = labels.NewLabelmgr()
	server.logmgr = logs.NewLogmgr()
	server.metricsmgr = metrics.NewMetricsmgr()
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
	server.cancel()
	server.configmgr = nil
	server.netmgr = nil
	server.entitymgr = nil
	server.labelmgr = nil
	server.logmgr = nil
	server.metricsmgr = nil
	server.lock.Unlock()
	contexts.Cancel() // 用來保證由contexts.Ctx()衍生出來的執行緒最後都能被終止, 避免goroutine洩漏
}

// Stop 關閉伺服器
func Stop() {
	server.lock.RLock()
	defer server.lock.RUnlock()

	if server.cancel != nil {
		server.cancel()
	} // if
}

// Name 取得伺服器名稱
func Name() string {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return server.name
}

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

// server 伺服器資料
var server struct {
	name       string              // 伺服器名稱
	ctx        context.Context     // ctx物件
	cancel     context.CancelFunc  // 取消物件
	configmgr  *configs.Configmgr  // 配置管理器
	metricsmgr *metrics.Metricsmgr // 統計管理器
	logmgr     *logs.Logmgr        // 日誌管理器
	netmgr     *nets.Netmgr        // 網路管理器
	entitymgr  *entitys.Entitymgr  // 實體管理器
	labelmgr   *labels.Labelmgr    // 標籤管理器
	lock       sync.RWMutex        // 執行緒鎖
}
