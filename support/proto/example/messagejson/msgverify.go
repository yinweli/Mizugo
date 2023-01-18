package messages

// MPingJsonQ 要求PingJson
type MPingJsonQ struct {
	Time int64 // 傳送時間
}

// MPingJsonA 回應PingJson
type MPingJsonA struct {
	From  *MPingJsonQ // 來源訊息
	Count int64       // 封包計數
}
