package procs

// SimpleMsg 簡單訊息資料
type SimpleMsg struct {
	MessageID MessageID `json:"messageID"` // 訊息編號
	Message   string    `json:"message"`   // 訊息字串
	Sum       string    `json:"sum"`       // 驗證字串
}
