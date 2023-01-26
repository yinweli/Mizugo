package msgs

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// TestMsg 僅用在測試的訊息資料
type TestMsg struct {
	MessageID int32         // 訊息編號
	Message   proto.Message // 訊息資料
}

// MarshalProtoMsg 序列化測試訊息到Proto訊息
func MarshalProtoMsg(input *TestMsg) *ProtoMsg {
	message, _ := anypb.New(input.Message)
	return &ProtoMsg{
		MessageID: input.MessageID,
		Message:   message,
	}
}

// MarshalPListMsg 序列化測試訊息到PList訊息
func MarshalPListMsg(input []TestMsg) *PListMsg {
	result := &PListMsg{}

	for _, itor := range input {
		message, _ := anypb.New(itor.Message)
		result.Messages = append(result.Messages, &PListUnit{
			MessageID: itor.MessageID,
			Message:   message,
		})
	} // for

	return result
}
