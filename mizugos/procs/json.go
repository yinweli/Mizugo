package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/msgs"
)

// NewJson 建立json處理器
func NewJson() *Json {
	return &Json{
		Procmgr: NewProcmgr(),
	}
}

// Json json處理器, 封包結構使用JsonMsg
//   - 訊息定義: support/proto/mizugo/msg-go/msgs-json/jsonmsg.go
//   - 訊息定義: support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
//   - 封包編碼: json編碼成位元陣列
//   - 封包解碼: json解碼成訊息結構
type Json struct {
	*Procmgr // 管理器
}

// Encode 封包編碼
func (this *Json) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("json encode: input nil")
	} // if

	temp, ok := input.(*msgs.JsonMsg)

	if ok == false {
		return nil, fmt.Errorf("json encode: input type")
	} // if

	output, err = json.Marshal(temp)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Json) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("json decode: input nil")
	} // if

	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("json decode: input type")
	} // if

	message := &msgs.JsonMsg{}

	if err = json.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Json) Process(input any) error {
	if input == nil {
		return fmt.Errorf("json process: input nil")
	} // if

	message, ok := input.(*msgs.JsonMsg)

	if ok == false {
		return fmt.Errorf("json process: input type")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("json process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// JsonMarshal json訊息序列化
func JsonMarshal(messageID MessageID, input any) (output *msgs.JsonMsg, err error) {
	if input == nil {
		return nil, fmt.Errorf("json marshal: input nil")
	} // if

	message, err := json.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	} // if

	return &msgs.JsonMsg{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// JsonUnmarshal json訊息反序列化
func JsonUnmarshal[T any](input any) (messageID MessageID, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("json unmarshal: input nil")
	} // if

	message, ok := input.(*msgs.JsonMsg)

	if ok == false {
		return 0, nil, fmt.Errorf("json unmarshal: input type")
	} // if

	output = new(T)

	if err = json.Unmarshal(message.Message, output); err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}
