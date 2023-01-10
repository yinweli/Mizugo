package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// 簡單處理器, 封包結構使用SimpleMsg
// 訊息內容定義在mizugos/procs/simplemsg.proto
// 封包編碼通過json編碼成位元陣列, 再通過Base64簡單編碼
// 封包解碼通過Base64簡單解碼, 再通過json解碼成封包結構
// 由於沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用

// NewSimple 建立簡單處理器
func NewSimple() *Simple {
	return &Simple{
		Procmgr: NewProcmgr(),
	}
}

// Simple 簡單處理器
type Simple struct {
	*Procmgr // 處理管理器
}

// Encode 封包編碼
func (this *Simple) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[SimpleMsg](input)

	if err != nil {
		return nil, fmt.Errorf("simple encode: %w", err)
	} // if

	bytes, err := json.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("simple encode: %w", err)
	} // if

	output = utils.Base64Encode(bytes)
	return output, nil
}

// Decode 封包解碼
func (this *Simple) Decode(input []byte) (output any, err error) {
	bytes, err := utils.Base64Decode(input)

	if err != nil {
		return nil, fmt.Errorf("simple decode: %w", err)
	} // if

	message := &SimpleMsg{}

	if err := json.Unmarshal(bytes, message); err != nil {
		return nil, fmt.Errorf("simple decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Simple) Process(input any) error {
	message, err := utils.CastPointer[SimpleMsg](input)

	if err != nil {
		return fmt.Errorf("simple process: %w", err)
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("simple process: messageID not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// SimpleMarshal 序列化簡單訊息
func SimpleMarshal(messageID MessageID, input any) (output *SimpleMsg, err error) {
	message, err := json.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("simple marshal: %w", err)
	} // if

	return &SimpleMsg{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// SimpleUnmarshal 反序列化簡單訊息
func SimpleUnmarshal[T any](input any) (messageID MessageID, output *T, err error) {
	message, err := utils.CastPointer[SimpleMsg](input)

	if err != nil {
		return 0, output, fmt.Errorf("simple unmarshal: %w", err)
	} // if

	output = new(T)

	if err := json.Unmarshal(message.Message, output); err != nil {
		return 0, output, fmt.Errorf("simple unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}
