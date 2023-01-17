package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/utils"
)

// proto處理器, 封包結構使用ProtoMsg
// 沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用
// 訊息內容: support/proto/mizugo/protomsg.proto
// 封包編碼: protobuf編碼成位元陣列, 再通過base64編碼
// 封包解碼: base64解碼, 再通過protobuf解碼成訊息結構

// NewProto 建立proto處理器
func NewProto() *Proto {
	return &Proto{
		procmgr: newProcmgr(),
	}
}

// Proto proto處理器
type Proto struct {
	*procmgr // 管理器
}

// Encode 封包編碼
func (this *Proto) Encode(input any) (output []byte, err error) {
	message, err := utils.CastPointer[msgs.ProtoMsg](input)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	bytes, err := proto.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	output = utils.Base64Encode(bytes)
	return output, nil
}

// Decode 封包解碼
func (this *Proto) Decode(input []byte) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto decode: input nil")
	} // if

	bytes, err := utils.Base64Decode(input)

	if err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	} // if

	message := &msgs.ProtoMsg{}

	if err = proto.Unmarshal(bytes, message); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Proto) Process(input any) error {
	message, err := utils.CastPointer[msgs.ProtoMsg](input)

	if err != nil {
		return fmt.Errorf("proto process: %w", err)
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("proto process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// ProtoMarshal 序列化proto訊息
func ProtoMarshal(messageID MessageID, input proto.Message) (output *msgs.ProtoMsg, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto marshal: input nil")
	} // if

	message, err := anypb.New(input)

	if err != nil {
		return nil, fmt.Errorf("proto marshal: %w", err)
	} // if

	return &msgs.ProtoMsg{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// ProtoUnmarshal 反序列化proto訊息
func ProtoUnmarshal(input any) (messageID MessageID, output proto.Message, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("proto unmarshal: input nil")
	} // if

	message, err := utils.CastPointer[msgs.ProtoMsg](input)

	if err != nil {
		return 0, nil, fmt.Errorf("proto unmarshal: %w", err)
	} // if

	if output, err = message.Message.UnmarshalNew(); err != nil {
		return 0, nil, fmt.Errorf("proto unmarshal: %w", err)
	} // if

	return message.MessageID, output, nil
}
