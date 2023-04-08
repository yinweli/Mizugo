package defines

import (
	"github.com/yinweli/Mizugo/mizugos/entitys"
)

const ( // 程式名稱
	CmdServer = "test-server" // 測試伺服器名稱
)

const ( // 常數定義
	EventCapacity = 1000 // 事件容量
)

const ( // 配置定義
	ConfigPath = "config" // 配置路徑
	ConfigFile = "config" // 配置名稱
	ConfigType = "yaml"   // 配置類型
)

const ( // 模組編號
	ModuleIDAuth entitys.ModuleID = iota + 1
	ModuleIDJson
	ModuleIDProto
)
