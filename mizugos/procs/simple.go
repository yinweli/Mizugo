package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// 簡單處理器, 封包結構使用SimpleMsg
// 封包編碼通過json編碼成位元陣列, 再通過Base64簡單編碼
// 封包解碼通過Base64簡單解碼, 再通過json解碼成封包結構
// 由於沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用

// NewSimpleMsg 建立簡單訊息
func NewSimpleMsg(messageID MessageID, message []byte) *SimpleMsg {
	return &SimpleMsg{
		MessageID: messageID,
		Message:   message,
	}
}

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
func (this *Simple) Encode(message any) (packet []byte, err error) {
	msg, err := utils.CastPointer[SimpleMsg](message)

	if err != nil {
		return nil, fmt.Errorf("simple encode: %w", err)
	} // if

	bytes, err := json.Marshal(msg)

	if err != nil {
		return nil, fmt.Errorf("simple encode: %w", err)
	} // if

	packet = utils.Base64Encode(bytes)
	return packet, nil
}

// Decode 封包解碼
func (this *Simple) Decode(packet []byte) (message any, err error) {
	bytes, err := utils.Base64Decode(packet)

	if err != nil {
		return nil, fmt.Errorf("simple decode: %w", err)
	} // if

	msg := &SimpleMsg{}

	if err := json.Unmarshal(bytes, msg); err != nil {
		return nil, fmt.Errorf("simple decode: %w", err)
	} // if

	return msg, nil
}

// Process 訊息處理
func (this *Simple) Process(message any) error {
	msg, err := utils.CastPointer[SimpleMsg](message)

	if err != nil {
		return fmt.Errorf("simple process: %w", err)
	} // if

	process := this.Get(msg.MessageID)

	if process == nil {
		return fmt.Errorf("simple process: messageID not found")
	} // if

	process(msg.MessageID, msg)
	return nil
}
