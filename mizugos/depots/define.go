package depots

import (
	"time"
)

const ( // redis定義
	Timeout = time.Second * 30 // redis超時時間
	Nil     = ""               // redis回應空字串, 通常在GET命令找不到索引時, 會以此字串回報給使用者
	Ok      = "OK"             // redis回應完成, 通常在SET命令順利完成後, 會以此字串回報給使用者
)

const ( // redis索引前綴定義, 前綴最後須以`:`結束
	prefixLock = "lock:" // 鎖定前綴詞
)
