package trials

import (
	"time"
)

const Timeout = time.Millisecond * 200 // 預設等待時間(200 毫秒)

// WaitTimeout 進行等待(Sleep), 可自訂等待時間
//   - 若有傳入參數, 則使用第一個參數作為等待時間
//   - 若無參數 則使用預設的 Timeout
func WaitTimeout(duration ...time.Duration) {
	if len(duration) > 0 {
		time.Sleep(duration[0])
	} else {
		time.Sleep(Timeout)
	} // if
}
