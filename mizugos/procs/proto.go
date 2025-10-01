package procs

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
)

// NewProto 建立 Proto 處理器
func NewProto() *Proto {
	return &Proto{
		Procmgr: NewProcmgr(),
	}
}

// Proto 處理器
//
// 訊息結構基於 msgs.Proto
//   - messageID: 訊息編號
//   - message: 任意 Proto 訊息
//
// 對應的訊息定義檔:
//   - Go: support/proto-mizugo/msg-go/msgs/proto.pb.go
//   - C#: support/proto-mizugo/msg-cs/msgs/Proto.cs
//   - Proto 定義: support/proto-mizugo/proto.proto
//
// 職責:
//   - Encode: 將 *msgs.Proto 物件轉為 []byte 以利傳輸
//   - Decode: 將 []byte 解碼回 *msgs.Proto
//   - Process: 根據 messageID 呼叫對應的處理函式
//
// 另外提供 ProtoMarshal / ProtoUnmarshal 協助處理「payload ↔ msgs.Proto」的轉換
type Proto struct {
	*Procmgr // 管理器
}

// Encode 訊息編碼
//
// 輸入必須是 *msgs.Proto, 輸出為 []byte
func (this *Proto) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto encode: input nil")
	} // if

	temp, ok := input.(*msgs.Proto)

	if ok == false {
		return nil, fmt.Errorf("proto encode: input not *msgs.Proto")
	} // if

	output, err = proto.Marshal(temp)

	if err != nil {
		return nil, fmt.Errorf("proto encode: %w", err)
	} // if

	return output, nil
}

// Decode 訊息解碼
//
// 輸入必須是 []byte, 輸出為 *msgs.Proto
func (this *Proto) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto decode: input nil")
	} // if

	temp, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("proto decode: input not []byte")
	} // if

	message := &msgs.Proto{}

	if err = proto.Unmarshal(temp, message); err != nil {
		return nil, fmt.Errorf("proto decode: %w", err)
	} // if

	return message, nil
}

// Process 訊息處理
//
// 輸入必須是 *msgs.Proto, 會依據 messageID 找出並執行已註冊的處理函式, 若找不到對應的處理函式則回傳錯誤
func (this *Proto) Process(input any) error {
	if input == nil {
		return fmt.Errorf("proto process: input nil")
	} // if

	message, ok := input.(*msgs.Proto)

	if ok == false {
		return fmt.Errorf("proto process: input not *msgs.Proto")
	} // if

	process := this.Get(message.MessageID)

	if process == nil {
		return fmt.Errorf("proto process: not found: %v", message.MessageID)
	} // if

	process(message)
	return nil
}

// ProtoMarshal Proto 訊息序列化
//
// 將 messageID 與 payload 序列化為 *msgs.Proto
func ProtoMarshal(messageID int32, input proto.Message) (output *msgs.Proto, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto marshal: input nil")
	} // if

	message, err := anypb.New(input)

	if err != nil {
		return nil, fmt.Errorf("proto marshal: %w", err)
	} // if

	return &msgs.Proto{
		MessageID: messageID,
		Message:   message,
	}, nil
}

// ProtoUnmarshal Proto 訊息反序列化
//
// 將 *msgs.Proto 反序列化為 messageID 與 payload
func ProtoUnmarshal[T proto.Message](input any) (messageID int32, output T, err error) {
	if input == nil {
		return 0, output, fmt.Errorf("proto unmarshal: input nil")
	} // if

	message, ok := input.(*msgs.Proto)

	if ok == false {
		return 0, output, fmt.Errorf("proto unmarshal: input not *msgs.Proto")
	} // if

	if output, err = helps.FromProtoAny[T](message.Message); err != nil {
		return 0, output, fmt.Errorf("proto unmarshal: message: %w", err)
	} // if

	return message.MessageID, output, nil
}
