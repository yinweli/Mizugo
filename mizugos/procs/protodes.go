package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// protoDes處理器, 封包結構使用ProtoDesMsg
// ProtoDesMsg由support/proto/protodesmsg/protodesmsg.bat產生
// ProtoDesMsg由support/proto/protodesmsg/protodesmsg.proto定義
// 封包編碼通過proto編碼成位元陣列, 再通過Des加密
// 封包解碼通過Des解密, 再通過proto解碼成封包結構
// 由於使用到Des加密, 所以需要在建立處理器時指定密鑰
// 安全性較高, 適合用來傳送一般封包

// NewProtoDes 建立protoDes處理器
func NewProtoDes(key []byte) *ProtoDes {
	return &ProtoDes{
		Procmgr: NewProcmgr(),
		key:     key,
	}
}

// ProtoDes protoDes處理器
type ProtoDes struct {
	*Procmgr        // 處理管理器
	key      []byte // 密鑰
}

// Encode 封包編碼
func (this *ProtoDes) Encode(message any) (packet []byte, err error) {
	msg, err := utils.CastPointer[ProtoDesMsg](message)

	if err != nil {
		return nil, fmt.Errorf("protodes encode: %w", err)
	} // if

	bytes, err := proto.Marshal(msg)

	if err != nil {
		return nil, fmt.Errorf("protodes encode: %w", err)
	} // if

	packet, err = utils.DesEncrypt(this.key, bytes)

	if err != nil {
		return nil, fmt.Errorf("protodes encode: %w", err)
	} // if

	return packet, nil
}

// Decode 封包解碼
func (this *ProtoDes) Decode(packet []byte) (message any, err error) {
	if packet == nil {
		return nil, fmt.Errorf("protodes decode: packet nil")
	} // if

	bytes, err := utils.DesDecrypt(this.key, packet)

	if err != nil {
		return nil, fmt.Errorf("protodes decode: %w", err)
	} // if

	msg := &ProtoDesMsg{}

	if err := proto.Unmarshal(bytes, msg); err != nil {
		return nil, fmt.Errorf("protodes decode: %w", err)
	} // if

	return msg, nil
}

// Process 訊息處理
func (this *ProtoDes) Process(message any) error {
	msg, err := utils.CastPointer[ProtoDesMsg](message)

	if err != nil {
		return fmt.Errorf("protodes process: %w", err)
	} // if

	process := this.Get(msg.MessageID)

	if process == nil {
		return fmt.Errorf("protodes process: messageID not found")
	} // if

	process(msg.MessageID, msg)
	return nil
}
