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
func Start(name string, logger logs.Logger, initialize Initialize, finalize Finalize) {
	if serv.init.CompareAndSwap(false, true) == false {
		return
	} // if

	fmt.Printf("%v initialize\n", name)
	serv.configmgr = configs.NewConfigmgr()
	serv.netmgr = nets.NewNetmgr()
	serv.entitymgr = entitys.NewEntitymgr()
	serv.tagmgr = tags.NewTagmgr()
	serv.logger = logger

	if err := serv.logger.Initialize(); err != nil {
		fmt.Println(fmt.Errorf("logger: %w", err))
		goto Finish
	} // if

	if err := initialize(); err != nil {
		fmt.Println(fmt.Errorf("initialize: %w", err))
		goto Finish
	} // if

	fmt.Printf("%v start\n", name)
	serv.exec.Store(true)
	serv.close.Add(1)
	serv.close.Wait() // 進行等待, 直到關閉伺服器
	goto Shutdown

Shutdown: // 關閉伺服器處理
	finalize()
	serv.logger.Finalize()
	fmt.Printf("%v shutdown\n", name)

Finish: // 結束處理
	serv.netmgr = nil
	serv.entitymgr = nil
	serv.tagmgr = nil
	serv.logger = nil
	serv.init.Store(false)
	serv.exec.Store(false)
	fmt.Printf("%v finish\n", name)
}

// Close 關閉伺服器
func Close() {
	serv.close.Done()
}

// Configmgr 取得配置管理器
func Configmgr() *configs.Configmgr {
	if serv.exec.Load() {
		return serv.configmgr
	} // if

	return nil
}

// Netmgr 取得網路管理器
func Netmgr() *nets.Netmgr {
	if serv.exec.Load() {
		return serv.netmgr
	} // if

	return nil
}

// Entitymgr 實體管理器
func Entitymgr() *entitys.Entitymgr {
	if serv.exec.Load() {
		return serv.entitymgr
	} // if

	return nil
}

// Tagmgr 實體管理器
func Tagmgr() *tags.Tagmgr {
	if serv.exec.Load() {
		return serv.tagmgr
	} // if

	return nil
}

// Debug 記錄除錯訊息
func Debug(label string) logs.Stream {
	if serv.exec.Load() {
		return serv.logger.New(label, logs.LevelDebug)
	} // if

	return nil
}

// Info 記錄一般訊息
func Info(label string) logs.Stream {
	if serv.exec.Load() {
		return serv.logger.New(label, logs.LevelInfo)
	} // if

	return nil
}

// Warn 記錄警告訊息
func Warn(label string) logs.Stream {
	if serv.exec.Load() {
		return serv.logger.New(label, logs.LevelWarn)
	} // if

	return nil
}

// Error 記錄錯誤訊息
func Error(label string) logs.Stream {
	if serv.exec.Load() {
		return serv.logger.New(label, logs.LevelError)
	} // if

	return nil
}

// serv 伺服器資料
var serv struct {
	init      atomic.Bool        // 初始化旗標
	exec      atomic.Bool        // 執行旗標
	close     sync.WaitGroup     // 關閉旗標
	configmgr *configs.Configmgr // 配置管理器
	netmgr    *nets.Netmgr       // 網路管理器
	entitymgr *entitys.Entitymgr // 實體管理器
	tagmgr    *tags.Tagmgr       // 標籤管理器
	logger    logs.Logger        // 日誌管理器
}