package metrics

import (
	"expvar"
	"time"
)

const ( // 時間定義
	serverTimeout = time.Second * 5 // http伺服器超時時間
	interval1     = 60              // 間隔時間: 1分鐘
	interval5     = 300             // 間隔時間: 5分鐘
	interval10    = 600             // 間隔時間: 10分鐘
	interval60    = 3600            // 間隔時間: 60分鐘
)

// Int 整數統計
type Int = expvar.Int

// Float 浮點數統計
type Float = expvar.Float

// String 字串統計
type String = expvar.String

// Map 映射統計
type Map = expvar.Map
