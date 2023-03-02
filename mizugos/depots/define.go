package depots

import (
	"time"
)

const ( // redis定義
	Timeout  = time.Second * 30 // redis超時時間
	RedisNil = ""               // redis回應空字串, 通常在GET命令找不到索引時, 會以此字串回報給使用者
	RedisOk  = "OK"             // redis回應完成, 通常在SET命令順利完成後, 會以此字串回報給使用者
)
