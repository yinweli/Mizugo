package metrics

import (
	"expvar"
	"time"
)

const ( // 時間定義
	timeout    = time.Second * 5 // http伺服器超時時間
	interval1  = 60              // 間隔時間: 1分鐘
	interval5  = 300             // 間隔時間: 5分鐘
	interval10 = 600             // 間隔時間: 10分鐘
	interval60 = 3600            // 間隔時間: 60分鐘
)

// Int 整數度量
type Int = expvar.Int

// Float 浮點數度量
type Float = expvar.Float

// String 字串度量
type String = expvar.String

// Map 映射度量
type Map = expvar.Map
