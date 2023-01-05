package metrics

import (
	"time"
)

const (
	pattern       = "/debug/vars"   // 監控路由字串
	serverTimeout = time.Second * 5 // 監控伺服器逾時時間
	interval1     = 60              // 間隔時間: 1分鐘
	interval5     = 300             // 間隔時間: 5分鐘
	interval10    = 600             // 間隔時間: 10分鐘
	interval60    = 3600            // 間隔時間: 60分鐘
)

// 如果需要查看統計數據, 可以通過以下工具
// https://github.com/divan/expvarmon
// 如果想查看記憶體使用情況, 可使用以下參數
// -ports="http://帳號:密碼@網址:埠號"
// -i 間隔時間
