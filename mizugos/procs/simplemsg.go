package procs

// SimpleMsg 簡單訊息資料
type SimpleMsg struct {
	MessageID MessageID `json:"messageID"` // 訊息編號
	Message   []byte    `json:"message"`   // 訊息資料
}
