package mizugos

import (
	"github.com/yinweli/Mizugo/mizugos/configs"
	"github.com/yinweli/Mizugo/mizugos/entitys"
	"github.com/yinweli/Mizugo/mizugos/labels"
	"github.com/yinweli/Mizugo/mizugos/loggers"
	"github.com/yinweli/Mizugo/mizugos/metrics"
	"github.com/yinweli/Mizugo/mizugos/nets"
	"github.com/yinweli/Mizugo/mizugos/pools"
	"github.com/yinweli/Mizugo/mizugos/redmos"
	"github.com/yinweli/Mizugo/mizugos/triggers"
)

// Start 啟動伺服器
/*
範例:
	defer func() {
	    if cause := recover(); cause != nil {
	        // 處理崩潰錯誤
	    } // if
	}()
	mizugos.Start() // 啟動伺服器
	ctx, cancel := context.WithCancel(context.Background())
	// 使用者自訂的初始化程序
	// 如果有任何失敗, 執行 mizugos.Stop() 後退出
	for range ctx.Done() { // 進入無限迴圈直到執行 ctx.Cancel()
	} // for
	// 使用者自訂的結束程序
	// 如果有任何失敗, 執行 mizugos.Stop() 後退出
    cancel()
	mizugos.Stop() // 關閉伺服器
*/
func Start() {
	Config = configs.NewConfigmgr()
	Metrics = metrics.NewMetricsmgr()
	Logger = loggers.NewLogmgr()
	Network = nets.NewNetmgr()
	Redmo = redmos.NewRedmomgr()
	Entity = entitys.NewEntitymgr()
	Label = labels.NewLabelmgr()
	Pool = pools.DefaultPool // 執行緒池管理器直接用預設的
	Trigger = triggers.NewTriggermgr()
}

// Stop 關閉伺服器
func Stop() {
	if Config != nil {
		Config = nil
	} // if

	if Metrics != nil {
		Metrics.Finalize()
		Metrics = nil
	} // if

	if Logger != nil {
		Logger.Finalize()
		Logger = nil
	} // if

	if Network != nil {
		Network.Stop()
		Network = nil
	} // if

	if Redmo != nil {
		Redmo.Finalize()
		Redmo = nil
	} // if

	if Entity != nil {
		Entity.Clear()
		Entity = nil
	} // if

	if Label != nil {
		Label = nil
	} // if

	if Pool != nil {
		Pool.Finalize()
		Pool = nil
	} // if

	if Trigger != nil {
		Trigger.Finalize()
		Trigger = nil
	} // if
}

var Config *configs.Configmgr    // 配置管理器
var Metrics *metrics.Metricsmgr  // 度量管理器
var Logger *loggers.Logmgr       // 日誌管理器
var Network *nets.Netmgr         // 網路管理器
var Redmo *redmos.Redmomgr       // 資料庫管理器
var Entity *entitys.Entitymgr    // 實體管理器
var Label *labels.Labelmgr       // 標籤管理器
var Pool *pools.Poolmgr          // 執行緒池管理器
var Trigger *triggers.Triggermgr // 信號調度管理器
