package defines

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
)

const EventCapacity = 1000 // 事件容量

const ( // 配置定義
	ConfigPath = "config" // 配置路徑
	ConfigFile = "config" // 配置名稱
	ConfigType = "yaml"   // 配置類型
)

const ( // 模組編號
	ModuleIDJson entitys.ModuleID = iota + 1
	ModuleIDProto
	ModuleIDPList
)

const ( // 事件名稱
	EventKey   = "key"
	EventBegin = "begin"
)
