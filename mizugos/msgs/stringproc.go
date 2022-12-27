package msgs

import (
	"encoding/json"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// NewStringProc 建立字串處理器
func NewStringProc() *StringProc {
	return &StringProc{
		Procmgr: NewProcmgr(),
	}
}

// StringProc 字串處理器
type StringProc struct {
	*Procmgr // 處理管理器
}

// Encode 封包編碼
func (this *StringProc) Encode(message any) (packet []byte, err error) {
	msg, err := Cast[StringMsg](message)

	if err != nil {
		return nil, fmt.Errorf("stringmsg encode: %w", err)
	} // if

	msg.Sum = utils.MD5String([]byte(msg.Message))
	bytes, err := json.Marshal(msg)

	if err != nil {
		return nil, fmt.Errorf("stringmsg encode: %w", err)
	} // if

	packet = utils.Base64Encode(bytes)
	return packet, nil
}

// Decode 封包解碼
func (this *StringProc) Decode(packet []byte) (message any, err error) {
	bytes, err := utils.Base64Decode(packet)

	if err != nil {
		return nil, fmt.Errorf("stringmsg decode: %w", err)
	} // if

	msg := &StringMsg{}

	if err := json.Unmarshal(bytes, msg); err != nil {
		return nil, fmt.Errorf("stringmsg decode: %w", err)
	} // if

	sum := utils.MD5String([]byte(msg.Message))

	if msg.Sum != sum {
		return nil, fmt.Errorf("stringmsg decode: sum failed")
	} // if

	return msg, nil
}

// Process 訊息處理
func (this *StringProc) Process(message any) error {
	msg, err := Cast[StringMsg](message)

	if err != nil {
		return fmt.Errorf("stringmsg process: %w", err)
	} // if

	process := this.Get(msg.MessageID)

	if process == nil {
		return fmt.Errorf("stringmsg process: messageID not found")
	} // if

	process(msg.MessageID, msg)
	return nil
}

// StringMsg 字串訊息資料
type StringMsg struct {
	MessageID MessageID `json:"messageID"` // 訊息編號
	Message   string    `json:"message"`   // 訊息字串
	Sum       string    `json:"sum"`       // 驗證字串
}
