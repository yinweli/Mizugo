package procs

// 訊息處理, 核心由Processor介面與處理器組成
// * 簡介
//   負責封包編碼/封包解碼, 收到訊息時的處理流程, 管理訊息處理函式等
//   封包的加密與解密都會在封包編碼/封包解碼中實行
// * 處理機制
//   使用者可以選擇使用哪種處理機制
//   - Proto: 使用protobuf的簡單訊息處理機制
//   - Complex: 使用protobuf的複合訊息處理機制
// * 自訂處理機制
//   如果使用者想要自訂處理機制, 需要完成以下工作
//   - 建立訊息結構, 此結構最少必須包含MessageID
//     訊息結構如果要使用protobuf, 可以把定義檔放在support/proto/mizugo中
//   - 建立處理器, 此結構必須實現Processor介面
//     procs中有個procmgr, 可以用於負責Add/Del的處理, 可以在建立處理器時將procmgr組合進去

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
