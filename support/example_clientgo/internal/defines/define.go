package defines

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
)

const ConfigPath = "config" // 配置路徑
const ConfigType = "yaml"   // 配置類型
const EventCapacity = 1000  // 事件容量

const ( // 模組編號
	ModuleIDPingJson entitys.ModuleID = iota + 1
	ModuleIDPingProto
	ModuleIDPingStack
)

const ( // 事件名稱
	EventKey  = "key"
	EventPing = "ping"
)
