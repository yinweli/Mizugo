package metrics

import (
	"time"
)

const (
	serverTimeout = time.Second * 5 // http伺服器逾時時間
	interval1     = 60              // 間隔時間: 1分鐘
	interval5     = 300             // 間隔時間: 5分鐘
	interval10    = 600             // 間隔時間: 10分鐘
	interval60    = 3600            // 間隔時間: 60分鐘
)
