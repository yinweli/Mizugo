package trials

import (
	"time"
)

const Timeout = time.Millisecond * 200 // 預設等待時間

// WaitTimeout 等待超時時間, 可以輸入等待時間, 預設等待200毫秒
func WaitTimeout(duration ...time.Duration) {
	if len(duration) > 0 {
		time.Sleep(duration[0])
	} else {
		time.Sleep(Timeout)
	} // if
}
