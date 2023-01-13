package defines

import (
	"time"

	"github.com/yinweli/Mizugo/mizugos/entitys"
)

const ConfigPath = "config"                // 配置路徑
const ConfigType = "yaml"                  // 配置類型
const EventCapacity = 1000                 // 事件容量
const EchoCount = 16                       // 回音字串長度
const PingWaitTime = time.Millisecond * 10 // Ping等待時間, 讓伺服器有時間設置新的密鑰

const ( // 模組編號
	ModuleIDEcho entitys.ModuleID = iota + 1 // TODO: 考慮一下到底是人工產生模組編號, 還是用hash產生?
	ModuleIDKey
	ModuleIDPing
)

const ( // 事件名稱
	EventCompleteKey = "completeKey"
)
