package procs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

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
		return nil, fmt.Errorf("simplemsg encode: %w", err)
	} // if

	msg.Sum = utils.MD5String([]byte(msg.Message))
	bytes, err := json.Marshal(msg)

	if err != nil {
		return nil, fmt.Errorf("simplemsg encode: %w", err)
	} // if

	packet = utils.Base64Encode(bytes)
	return packet, nil
}

// Decode 封包解碼
func (this *Simple) Decode(packet []byte) (message any, err error) {
	bytes, err := utils.Base64Decode(packet)

	if err != nil {
		return nil, fmt.Errorf("simplemsg decode: %w", err)
	} // if

	msg := &SimpleMsg{}

	if err := json.Unmarshal(bytes, msg); err != nil {
		return nil, fmt.Errorf("simplemsg decode: %w", err)
	} // if

	sum := utils.MD5String([]byte(msg.Message))

	if msg.Sum != sum {
		return nil, fmt.Errorf("simplemsg decode: sum failed")
	} // if

	return msg, nil
}

// Process 訊息處理
func (this *Simple) Process(message any) error {
	msg, err := utils.CastPointer[SimpleMsg](message)

	if err != nil {
		return fmt.Errorf("simplemsg process: %w", err)
	} // if

	process := this.Get(msg.MessageID)

	if process == nil {
		return fmt.Errorf("simplemsg process: messageID not found")
	} // if

	process(msg.MessageID, msg)
	return nil
}
