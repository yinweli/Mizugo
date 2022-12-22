package msgs

// Messenger 訊息介面
type Messenger interface {
	// Encode 封包編碼
	Encode(message any) (packet []byte, err error)

	// Decode 封包解碼
	Decode(packet []byte) (message any, err error)

	// Process 訊息處理
	Process(message any) error

	// Add 新增訊息處理
	Add(messageID MessageID, process Process)

	// Del 刪除訊息處理
	Del(messageID MessageID)
}

// Process 訊息處理函式類型
type Process func(messageID MessageID, message any)

// MessageID 訊息編號
type MessageID = int64
