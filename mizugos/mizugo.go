package mizugos

import (
	"fmt"
	"sync/atomic"

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
	if serv.exec.CompareAndSwap(false, true) == false {
		return
	} // if

	fmt.Printf("%v initialize\n", name)

	serv.close = make(chan bool, 1)
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

	Info(name).Message("server start").End()

	// 進行無限迴圈, 直到關閉伺服器
	for {
		select {
		case <-serv.close:
			goto Shutdown
		} // select
	} // for

Shutdown:
	finalize()
	serv.logger.Finalize()
	fmt.Printf("%v shutdown\n", name)

Finish:
	close(serv.close)
	serv.netmgr = nil
	serv.entitymgr = nil
	serv.tagmgr = nil
	serv.logger = nil
	serv.exec.Store(false)
	fmt.Printf("%v finish\n", name)
}

// Close 關閉伺服器
func Close() {
	if serv.exec.Load() {
		serv.close <- true
	} // if
}

// Netmgr 取得網路管理器
func Netmgr() *nets.Netmgr {
	return serv.netmgr
}

// Entitymgr 實體管理器
func Entitymgr() *entitys.Entitymgr {
	return serv.entitymgr
}

// Tagmgr 實體管理器
func Tagmgr() *tags.Tagmgr {
	return serv.tagmgr
}

// Debug 記錄除錯訊息
func Debug(label string) logs.Stream {
	return serv.logger.New(label, logs.LevelDebug)
}

// Info 記錄一般訊息
func Info(label string) logs.Stream {
	return serv.logger.New(label, logs.LevelInfo)
}

// Warn 記錄警告訊息
func Warn(label string) logs.Stream {
	return serv.logger.New(label, logs.LevelWarn)
}

// Error 記錄錯誤訊息
func Error(label string) logs.Stream {
	return serv.logger.New(label, logs.LevelError)
}

// serv 伺服器資料
var serv struct {
	exec      atomic.Bool        // 執行旗標
	close     chan bool          // 關閉通道
	netmgr    *nets.Netmgr       // 網路管理器
	entitymgr *entitys.Entitymgr // 實體管理器
	tagmgr    *tags.Tagmgr       // 標籤管理器
	logger    logs.Logger        // 日誌管理器
}