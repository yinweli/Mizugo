package procs

// Processor 處理介面, 負責以下功能
//   - 封包編碼: 在 Encode 中實現
//   - 封包解碼: 在 Decode 中實現
//   - 收到訊息時的處理: 在 Process 中實現
//   - 管理訊息處理函式: 在 Add, Del 中實現
//
// 如果想要建立新的處理結構, 需要遵循以下流程
//   - 定義訊息結構, 訊息結構必須包含 MessageID
//   - 訊息結構如果要使用protobuf, 可以把定義檔放在support/proto/mizugo中
//   - 定義處理結構, 處理結構需要繼承 Processor 介面, 並實現所有函式;
//     在處理結構中包含 Procmgr 結構來實現訊息處理功能, 這樣只要實作 Encode, Decode, Process 三個函式就可以了
//
// mizugo提供的預設處理器有 Json, Proto
type Processor interface {
	// Encode 封包編碼
	Encode(input any) (output any, err error)

	// Decode 封包解碼
	Decode(input any) (output any, err error)

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
