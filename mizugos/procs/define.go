package procs

// Processor 處理介面
type Processor interface {
	// Encode 封包編碼
	Encode(input any) (output []byte, err error)

	// Decode 封包解碼
	Decode(input []byte) (output any, err error)

	// Process 訊息處理
	Process(input any) error

	// Add 新增訊息處理
	Add(messageID MessageID, process Process)

	// Del 刪除訊息處理
	Del(messageID MessageID)
}

// Process 訊息處理函式類型
type Process func(message any)

// MessageID 訊息編號, 設置為int32以跟proto的列舉類型統一
type MessageID = int32
