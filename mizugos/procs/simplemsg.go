package procs

// SimpleMsg 簡單訊息資料
type SimpleMsg struct {
	MessageID MessageID `json:"messageID"` // 訊息編號
	Message   []byte    `json:"message"`   // 訊息資料
}

// SimpleMsgTest 簡單訊息測試用資料
type SimpleMsgTest struct {
	Message string // 訊息內容
}
