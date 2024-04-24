package msgs

// JsonMsg json訊息資料
type JsonMsg struct {
	MessageID int32  `json:"messageID"` // 訊息編號
	Message   []byte `json:"message"`   // 訊息資料
}

// JsonTest json訊息測試用資料
type JsonTest struct {
	Data string // 測試字串
}
