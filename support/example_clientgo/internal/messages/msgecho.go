package messages

// MsgEchoReq 要求回音
type MsgEchoReq struct {
	Time int64  `json:"time"` // 傳送時間
	Echo string `json:"echo"` // 回音字串
}

// MsgEchoRes 回應回音
type MsgEchoRes struct {
	From  MsgEchoReq `json:"from"`  // 來源訊息
	Count int64      `json:"count"` // 封包計數
}
