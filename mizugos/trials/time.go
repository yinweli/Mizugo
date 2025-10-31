package trials

import (
	"time"
)

const Timeout = time.Millisecond * 200 // 預設等待時間(200 毫秒)

// WaitTimeout 進行等待時間
//   - 若有傳入參數, 則使用第一個參數作為等待時間
//   - 若無參數 則使用預設的 Timeout
func WaitTimeout(timeout ...time.Duration) {
	if len(timeout) > 0 {
		time.Sleep(timeout[0])
	} else {
		time.Sleep(Timeout)
	} // if
}

// WaitFor 進行等待條件, 或是直到逾時
func WaitFor(timeout time.Duration, cond func() bool) {
	deadline := time.Now().Add(timeout)

	for {
		if cond() {
			return
		} // if

		if time.Now().After(deadline) {
			return
		} // if

		time.Sleep(Timeout)
	} // for
}
