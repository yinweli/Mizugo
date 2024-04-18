package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
)

// NewProto 建立proto處理器
func NewProto() *Proto {
	return &Proto{
		Procmgr: NewProcmgr(),
	}
}

// Proto proto處理器, 封包結構使用ProtoMsg
//   - 訊息定義: support/proto/mizugo/protomsg.proto
//   - 封包編碼: protobuf編碼成位元陣列
//   - 封包解碼: protobuf解碼成訊息結構
type Proto struct {
	*Procmgr // 管理器
}

// Encode 封包編碼
func (this *Proto) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto encode: input nil")
	} // if

	temp, ok := input.(*msgs.ProtoMsg)

	if ok == false {
		return nil, fmt.Errorf("proto encode: input type")
	} // if

	output, err = proto.Marshal(temp)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Proto) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto decode: input nil")
	} // if

	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("proto decode: input type")
	} // if

	message := &msgs.ProtoMsg{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Proto) Process(input any) error {
	if input == nil {
		return fmt.Errorf("proto process: input nil")
	} // if

	message, ok := input.(*msgs.ProtoMsg)

	if ok == false {
		return fmt.Errorf("proto process: input type")
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
func ProtoUnmarshal[T any](input any) (messageID MessageID, output *T, err error) {
	if input == nil {
		return 0, nil, fmt.Errorf("proto unmarshal: input nil")
	} // if

	message, ok := input.(*msgs.ProtoMsg)

	if ok == false {
		return 0, nil, fmt.Errorf("proto unmarshal: input type")
	} // if

	temp, err := message.Message.UnmarshalNew()

	if err != nil {
		return 0, nil, fmt.Errorf("proto unmarshal: %w", err)
	} // if

	output, ok = temp.(any).(*T)

	if ok == false {
		return 0, nil, fmt.Errorf("proto unmarshal: cast failed")
	} // if

	return message.MessageID, output, nil
}
