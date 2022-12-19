package mizugos

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/yinweli/Mizugo/cores/configs"
	"github.com/yinweli/Mizugo/cores/entitys"
	"github.com/yinweli/Mizugo/cores/logs"
	"github.com/yinweli/Mizugo/cores/nets"
	"github.com/yinweli/Mizugo/cores/tags"
)

// Initialize 初始化處理函式類型
type Initialize func() error

// Finalize 結束處理函式類型
type Finalize func()

// Start 啟動伺服器, 為了讓程式持續執行, 此函式不能用執行緒執行
func Start(name string, initialize Initialize, finalize Finalize) {
	if serv.start.CompareAndSwap(false, true) == false {
		return
	} // if

	serv.lock.Lock()
	serv.name = name
	serv.configmgr = configs.NewConfigmgr()
	serv.netmgr = nets.NewNetmgr()
	serv.entitymgr = entitys.NewEntitymgr()
	serv.tagmgr = tags.NewTagmgr()
	serv.logmgr = logs.NewLogmgr()
	serv.lock.Unlock()

	fmt.Printf("%v initialize\n", name)

	if err := initialize(); err != nil {
		fmt.Println(fmt.Errorf("initialize: %w", err))
		goto Finalize
	} // if

	fmt.Printf("%v start\n", name)
	serv.close.Add(1)
	serv.close.Wait() // 進行等待, 直到關閉伺服器
	goto Shutdown

Shutdown: // 關閉伺服器處理
	fmt.Printf("%v shutdown\n", name)
	finalize()

Finalize: // 結束處理
	fmt.Printf("%v finalize\n", name)
	serv.lock.Lock()
	serv.name = ""
	serv.configmgr = nil
	serv.netmgr = nil
	serv.entitymgr = nil
	serv.tagmgr = nil
	serv.logmgr = nil
	serv.lock.Unlock()
	serv.start.Store(false)
}

// Close 關閉伺服器
func Close() {
	serv.close.Done()
}

// Name 取得伺服器名稱
func Name() string {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.name
}

// Configmgr 取得配置管理器
func Configmgr() *configs.Configmgr {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.configmgr
}

// Netmgr 取得網路管理器
func Netmgr() *nets.Netmgr {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.netmgr
}

// Entitymgr 實體管理器
func Entitymgr() *entitys.Entitymgr {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.entitymgr
}

// Tagmgr 實體管理器
func Tagmgr() *tags.Tagmgr {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.tagmgr
}

// Logmgr 日誌管理器
func Logmgr() *logs.Logmgr {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.logmgr
}

// Debug 記錄除錯訊息
func Debug(label string) logs.Stream {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.logmgr.Debug(label)
}

// Info 記錄一般訊息
func Info(label string) logs.Stream {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.logmgr.Info(label)
}

// Warn 記錄警告訊息
func Warn(label string) logs.Stream {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.logmgr.Warn(label)
}

// Error 記錄錯誤訊息
func Error(label string) logs.Stream {
	serv.lock.RLock()
	defer serv.lock.RUnlock()

	return serv.logmgr.Error(label)
}

// serv 伺服器資料
var serv struct {
	name      string             // 伺服器名稱
	configmgr *configs.Configmgr // 配置管理器
	netmgr    *nets.Netmgr       // 網路管理器
	entitymgr *entitys.Entitymgr // 實體管理器
	tagmgr    *tags.Tagmgr       // 標籤管理器
	logmgr    *logs.Logmgr       // 日誌管理器
	lock      sync.RWMutex       // 執行緒鎖
	start     atomic.Bool        // 啟動旗標
	close     sync.WaitGroup     // 關閉旗標
}
