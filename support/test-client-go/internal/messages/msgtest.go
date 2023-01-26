package messages

// MJsonQ 要求Json
type MJsonQ struct {
	Time int64 // 傳送時間
}

// MJsonA 回應Json
type MJsonA struct {
	From  *MJsonQ // 來源訊息
	Count int64   // 封包計數
}
