package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
)

// NewJson 建立 JSON 處理器
func NewJson() *Json {
	return &Json{
		Procmgr: NewProcmgr(),
	}
}

// Json 處理器
//
// 訊息結構基於 msgs.Json
//   - messageID: 訊息編號
//   - message: 訊息內容
//
// 對應的訊息定義檔:
//   - Go: support/proto-mizugo/msg-go/msgs-json/json.go
//   - C#: support/proto-mizugo/msg-cs/msgs-json/Json.cs
//
// 職責:
//   - Encode: 將 *msgs.Json 物件轉為 []byte 以利傳輸
//   - Decode: 將 []byte 解碼回 *msgs.Json
//   - Process: 根據 messageID 呼叫對應的處理函式
//
// 另外提供 JsonMarshal / JsonUnmarshal 協助處理「payload ↔ msgs.Json」的轉換
type Json struct {
	*Procmgr // 管理器
}

// Encode 訊息編碼
//
// 輸入必須是 *msgs.Json, 輸出為 []byte
func (this *Json) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("json encode: input nil")
	} // if

	temp, ok := input.(*msgs.Json)

	if ok == false {
		return nil, fmt.Errorf("json encode: input not *msgs.Json")
	} // if

	output, err = json.Marshal(temp)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	return output, nil
}

// Decode 訊息解碼
//
// 輸入必須是 []byte, 輸出為 *msgs.Json
func (this *Json) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("json decode: input nil")
	} // if

	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("json decode: input not []byte")
	} // if

	message := &msgs.Json{}

	if err = json.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
//
// 輸入必須是 *msgs.Json, 會依據 messageID 找出並執行已註冊的處理函式, 若找不到對應的處理函式則回傳錯誤
func (this *Json) Process(input any) error {
	if input == nil {
		return fmt.Errorf("json process: input nil")
	} // if

	message, ok := input.(*msgs.Json)

	if ok == false {
		return fmt.Errorf("json process: input not *msgs.Json")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("json process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// JsonMarshal JSON 訊息序列化
//
// 將 messageID 與 payload 序列化為 *msgs.Json
func JsonMarshal(messageID int32, input any) (output *msgs.Json, err error) {
	if input == nil {
		return nil, fmt.Errorf("json marshal: input nil")
	} // if

	message, err := json.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	} // if

	return &msgs.Json{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// JsonUnmarshal JSON 訊息反序列化
//
// 將 *msgs.Json 反序列化為 messageID 與 payload
func JsonUnmarshal[T any](input any) (messageID int32, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("json unmarshal: input nil")
	} // if

	message, ok := input.(*msgs.Json)

	if ok == false {
		return 0, nil, fmt.Errorf("json unmarshal: input not *msgs.Json")
	} // if

	output = new(T)

	if err = json.Unmarshal(message.Message, output); err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}
