package mizugo

import (
	"sync"

	"github.com/yinweli/Mizugo/core/entitys"
	"github.com/yinweli/Mizugo/core/logs"
	"github.com/yinweli/Mizugo/core/nets"
	"github.com/yinweli/Mizugo/core/tags"
)

// Initialize 初始化核心資料
func Initialize(logger logs.Logger) {
	once.Do(func() {
		inst = &instance{
			logger:    logger,
			netmgr:    nets.NewNetmgr(),
			entitymgr: entitys.NewEntitymgr(),
			tagmgr:    tags.NewTagmgr(),
		}
	})
}

// Logger 取得日誌管理器
func Logger() logs.Logger {
	return inst.logger
}

// Netmgr 取得網路管理器
func Netmgr() *nets.Netmgr {
	return inst.netmgr
}

// Entitymgr 實體管理器
func Entitymgr() *entitys.Entitymgr {
	return inst.entitymgr
}

// Tagmgr 實體管理器
func Tagmgr() *tags.Tagmgr {
	return inst.tagmgr
}

// instance 實體資料
type instance struct {
	logger    logs.Logger        // 日誌管理器
	netmgr    *nets.Netmgr       // 網路管理器
	entitymgr *entitys.Entitymgr // 實體管理器
	tagmgr    *tags.Tagmgr       // 標籤管理器
}

var inst *instance // 實體資料
var once sync.Once // 單次執行緒鎖
