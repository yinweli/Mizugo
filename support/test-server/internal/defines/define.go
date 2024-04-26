package defines

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
)

const ( // 常數定義
	EventCapacity = 1000 // 事件容量
)

const ( // 配置定義
	ConfigPath = "config" // 配置路徑
	ConfigType = "yaml"   // 配置類型
	ConfigEnv  = "test"   // 配置環境變數的前綴字
	ConfigFile = "config" // 配置名稱
)

const ( // 模組編號
	ModuleIDAuth entitys.ModuleID = iota + 1
	ModuleIDJson
	ModuleIDProto
	ModuleIDRaven
)
