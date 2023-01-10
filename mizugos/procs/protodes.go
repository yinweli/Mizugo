package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// protoDes處理器, 封包結構使用ProtoDesMsg
// 訊息內容定義在support/proto/protodesmsg/protodesmsg.proto
// 訊息產生則執行support/proto/protodesmsg/protodesmsg.bat
// 封包編碼通過proto編碼成位元陣列, 再通過Des加密
// 封包解碼通過Des解密, 再通過proto解碼成封包結構
// 由於使用到Des加密, 所以需要在建立處理器時指定密鑰
// 安全性較高, 適合用來傳送一般封包

// NewProtoDes 建立protoDes處理器
func NewProtoDes() *ProtoDes {
	return &ProtoDes{
		Procmgr: NewProcmgr(),
	}
}

// ProtoDes protoDes處理器
type ProtoDes struct {
	*Procmgr                        // 處理管理器
	key      utils.SyncAttr[[]byte] // 密鑰
}

// Key 設定密鑰
func (this *ProtoDes) Key(key []byte) *ProtoDes {
	this.key.Set(key)
	return this
}

// Encode 封包編碼
func (this *ProtoDes) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[ProtoDesMsg](input)

	if err != nil {
		return nil, fmt.Errorf("protodes encode: %w", err)
	} // if

	bytes, err := proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("protodes encode: %w", err)
	} // if

	output, err = utils.DesEncrypt(this.key.Get(), bytes)

	if err != nil {
		return nil, fmt.Errorf("protodes encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *ProtoDes) Decode(input []byte) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("protodes decode: packet nil")
	} // if

	bytes, err := utils.DesDecrypt(this.key.Get(), input)

	if err != nil {
		return nil, fmt.Errorf("protodes decode: %w", err)
	} // if

	message := &ProtoDesMsg{}

	if err := proto.Unmarshal(bytes, message); err != nil {
		return nil, fmt.Errorf("protodes decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *ProtoDes) Process(input any) error {
	message, err := utils.CastPointer[ProtoDesMsg](input)

	if err != nil {
		return fmt.Errorf("protodes process: %w", err)
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("protodes process: messageID not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// ProtoDesMarshal 序列化protoDes訊息
func ProtoDesMarshal(messageID MessageID, input proto.Message) (output *ProtoDesMsg, err error) {
	message, err := anypb.New(input)

	if err != nil {
		return nil, fmt.Errorf("protodes marshal: %w", err)
	} // if

	return &ProtoDesMsg{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// ProtoDesUnmarshal 反序列化protoDes訊息
func ProtoDesUnmarshal[T proto.Message](input any) (messageID MessageID, output T, err error) {
	message, err := utils.CastPointer[ProtoDesMsg](input)

	if err != nil {
		return 0, output, fmt.Errorf("protodes unmarshal: %w", err)
	} // if

	temp, err := message.Message.UnmarshalNew()

	if err != nil {
		return 0, output, fmt.Errorf("protodes unmarshal: %w", err)
	} // if

	output, ok := temp.(T)

	if ok == false {
		return 0, output, fmt.Errorf("protodes unmarshal: cast failed")
	} // if

	return message.MessageID, output, nil
}
