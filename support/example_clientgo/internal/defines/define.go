package defines

import (
	"github.com/yinweli/Mizugo/mizugos/procs"
)

const ConfigPath = "config" // 配置路徑
const ConfigType = "yaml"   // 配置類型

const ( // 入口名稱
	EntryEchoSingle = "echosingle"
	EntryEchoCycle  = "echocycle"
)

const ( // 標籤名稱
	LabelEchoSingle = "echosingle"
	LabelEchoCycle  = "echocycle"
)

const ( // 訊息編號
	MessageIDEcho = procs.MessageID(1)
)