package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
)

// NewJson 建立json處理器
func NewJson() *Json {
	return &Json{
		Procmgr: NewProcmgr(),
	}
}

// Json json處理器, 封包結構使用msgs.Json
//   - 訊息定義: support/proto/mizugo/msg-go/msgs-json/json.go
//   - 訊息定義: support/proto/mizugo/msg-cs/msgs-json/Json.cs
type Json struct {
	*Procmgr // 管理器
}

// Encode 封包編碼
func (this *Json) Encode(input any) (output any, err error) {
	temp, ok := input.(*msgs.Json)

	if ok == false {
		return nil, fmt.Errorf("json encode: input")
	} // if

	output, err = json.Marshal(temp)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Json) Decode(input any) (output any, err error) {
	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("json decode: input")
	} // if

	message := &msgs.Json{}

	if err = json.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Json) Process(input any) error {
	message, ok := input.(*msgs.Json)

	if ok == false {
		return fmt.Errorf("json process: input")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("json process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// JsonMarshal json訊息序列化
func JsonMarshal(messageID int32, input any) (output *msgs.Json, err error) {
	if input == nil {
		return nil, fmt.Errorf("json marshal: input")
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

// JsonUnmarshal json訊息反序列化
func JsonUnmarshal[T any](input any) (messageID int32, output *T, err error) {
	message, ok := input.(*msgs.Json)

	if ok == false {
		return 0, nil, fmt.Errorf("json unmarshal: input")
	} // if

	output = new(T)

	if err = json.Unmarshal(message.Message, output); err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}
