package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/msgs"
)

// NewProto 建立proto處理器
func NewProto() *Proto {
	return &Proto{
		Procmgr: NewProcmgr(),
	}
}

// Proto proto處理器, 封包結構使用msgs.Proto
//   - 訊息定義: support/proto/mizugo/proto.proto
type Proto struct {
	*Procmgr // 管理器
}

// Encode 封包編碼
func (this *Proto) Encode(input any) (output any, err error) {
	temp, ok := input.(*msgs.Proto)

	if ok == false {
		return nil, fmt.Errorf("proto encode: input")
	} // if

	output, err = proto.Marshal(temp)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	return output, nil
}

// Decode 封包解碼
func (this *Proto) Decode(input any) (output any, err error) {
	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("proto decode: input")
	} // if

	message := &msgs.Proto{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
func (this *Proto) Process(input any) error {
	message, ok := input.(*msgs.Proto)

	if ok == false {
		return fmt.Errorf("proto process: input")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("proto process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// ProtoMarshal 序列化proto訊息
func ProtoMarshal(messageID int32, input proto.Message) (output *msgs.Proto, err error) {
	message, err := anypb.New(input)

	if err != nil {
		return nil, fmt.Errorf("proto marshal: %w", err)
	} // if

	return &msgs.Proto{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// ProtoUnmarshal 反序列化proto訊息
func ProtoUnmarshal[T any](input any) (messageID int32, output *T, err error) {
	message, ok := input.(*msgs.Proto)

	if ok == false {
		return 0, nil, fmt.Errorf("proto unmarshal: input")
	} // if

	if output, err = helps.ProtoAny[T](message.Message); err != nil {
		return 0, nil, fmt.Errorf("proto unmarshal: message: %w", err)
	} // if

	return message.MessageID, output, nil
}
