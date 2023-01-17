package defines

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
)

const ConfigPath = "config" // 配置路徑
const ConfigType = "yaml"   // 配置類型
const EventCapacity = 1000  // 事件容量

const ( // 模組編號
	ModuleIDPing entitys.ModuleID = iota + 1
)

const ( // 事件名稱
	EventPing = "ping"
)

const (
	PingWaitTime = time.Millisecond * 10 // Ping等待時間, 讓伺服器有時間設置新的密鑰
)
