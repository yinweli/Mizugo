package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
)

// json處理器, 封包結構使用msgs.JsonMsg
// 沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用
// 訊息內容: mizugos/msgs/jsonmsg.proto
// 封包編碼: json編碼成位元陣列, 再通過base64編碼
// 封包解碼: base64解碼, 再通過json解碼成訊息結構

// NewJson 建立json處理器
func NewJson() *Json {
	return &Json{
		procmgr: newProcmgr(),
	}
}

// Json json處理器
type Json struct {
	*procmgr // 管理器
}

// Encode 封包編碼
func (this *Json) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[msgs.JsonMsg](input)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	bytes, err := json.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	} // if

	output = utils.Base64Encode(bytes)
	return output, nil
}

// Decode 封包解碼
func (this *Json) Decode(input []byte) (output any, err error) {
	bytes, err := utils.Base64Decode(input)

	if err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	} // if

	message := &msgs.JsonMsg{}

	if err = json.Unmarshal(bytes, message); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Json) Process(input any) error {
	message, err := utils.CastPointer[msgs.JsonMsg](input)

	if err != nil {
		return fmt.Errorf("json process: %w", err)
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

	message, err := utils.CastPointer[msgs.JsonMsg](input)

	if err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	output = new(T)

	if err = json.Unmarshal(message.Message, output); err != nil {
		return 0, nil, fmt.Errorf("json unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}
