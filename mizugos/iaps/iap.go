package iaps

import (
	"time"
)

// IAPResult IAP 驗證結果資料
type IAPResult struct {
	Err  error     // 驗證結果, 若為 nil 表示驗證成功, 否則失敗
	Time time.Time // 購買時間
}

const (
	capacity   = 10000                                         // 驗證通道容量
	retry      = 3                                             // 驗證重試次數
	interval   = time.Millisecond * 2                          // 驗證間隔時間
	timeout    = time.Second * 5                               // 驗證逾時時間
	timeoutMax = (timeout + interval) * time.Duration(retry+1) // 驗證最大時間
)

// channelTry 對通道做非阻塞送出, 通道若已經關閉就丟掉結果避免卡死
func channelTry[T any](c chan T, t T) {
	select {
	case c <- t:
	default:
	} // select
}
