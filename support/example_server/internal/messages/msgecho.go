package messages

import (
	"time"
)

// EchoReq 要求回音
type EchoReq struct {
	Time time.Duration `json:"time"` // 傳送時間
	Echo string        `json:"echo"` // 回音字串
}

// EchoRes 回應回音
type EchoRes struct {
	From  EchoReq `json:"from"`  // 來源訊息
	Count int64   `json:"count"` // 封包計數
}
