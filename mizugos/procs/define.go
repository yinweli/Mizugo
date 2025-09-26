package procs

// Processor 處理器介面
//
// 負責將網路傳輸的封包與程式內部的訊息物件之間做轉換, 並在收到訊息後呼叫對應的處理函式(Process)
//
// 職責:
//   - Encode: 將內部訊息物件轉換成封包格式(序列化)
//   - Decode: 將封包格式轉換回內部訊息物件(反序列化)
//   - Process: 接收封包後, 呼叫對應的訊息處理函式
//   - Add / Del / Get: 管理訊息處理函式的註冊與查詢
//
// 建立自訂處理器的流程:
//  1. 定義訊息結構, 必須包含 messageID 欄位 (int32), 作為訊息識別
//  2. 若使用 Protobuf, 可將定義檔放在 support/proto-mizugo 目錄下, 由編譯器自動生成程式碼
//  3. 定義一個結構實作 Processor 介面, 並內含 Procmgr, 只需自行實作 Encode / Decode / Process 三個函式, 其他訊息處理管理功能由 Procmgr 提供
//  4. 將自訂 Processor 加入到框架中使用
//
// mizugo 內建的處理器:
//   - Json: 使用 JSON 格式序列化/反序列化訊息
//   - Proto: 使用 Protobuf 格式序列化/反序列化訊息
//   - Raven: 使用 Raven 格式序列化/反序列化訊息
type Processor interface {
	// Encode 封包編碼, 將內部訊息物件轉換成封包資料, 用於傳輸
	Encode(input any) (output any, err error)

	// Decode 封包解碼, 將封包資料轉換回內部訊息物件, 用於後續處理
	Decode(input any) (output any, err error)

	// Process 訊息處理, 當解碼成功後, 會呼叫對應的訊息處理函式
	Process(input any) error

	// Add 新增訊息處理
	Add(messageID int32, process Process)

	// Del 刪除訊息處理
	Del(messageID int32)

	// Get 取得訊息處理
	Get(messageID int32) Process
}

// Process 訊息處理函式類型
//   - 每一個 messageID 對應一個 Process 函式
//   - 當訊息解碼成功後, 會呼叫對應的 Process 並傳入訊息物件
//   - 使用者需自行將 message 轉換成正確型別 (例如 *msgs.Json / *msgs.Proto)
type Process func(message any)
