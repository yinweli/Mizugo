package packets

// Packeter 封包介面
type Packeter interface {
	// Encode 封包編碼
	Encode(message any) (packet []byte, err error)

	// Decode 封包解碼
	Decode(packet []byte) (message any, err error)

	// Process 訊息處理
	Process(message any)

	// Add 新增訊息處理
	Add(messageID MessageID, messenger Messenger)

	// Del 刪除訊息處理
	Del(messageID MessageID)
}

// Messenger 訊息處理函式類型
type Messenger func(messageID MessageID, message any)

// MessageID 訊息編號
type MessageID = int64
