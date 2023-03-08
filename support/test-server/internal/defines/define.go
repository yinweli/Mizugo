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
	ModuleIDAuth entitys.ModuleID = iota + 1
	ModuleIDJson
	ModuleIDProto
)

const ( // 資料庫名稱
	MajorName = "major" // 主要資料庫名稱
	MinorName = "minor" // 次要資料庫名稱
	MixedName = "mixed" // 混合資料庫名稱
)

const ( // mongo資料庫名稱
	MongoDB    = "auth" // mongo資料庫名稱
	MongoTable = "auth" // mongo資料表名稱
)
