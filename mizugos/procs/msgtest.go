package procs

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/msgs"
)

// TestMsg 僅用在測試的訊息資料
type TestMsg struct {
	MessageID int32         // 訊息編號
	Message   proto.Message // 訊息資料
}

// MarshalProtoMsg 序列化測試訊息到Proto訊息
func MarshalProtoMsg(input *TestMsg) *msgs.ProtoMsg {
	message, _ := anypb.New(input.Message)
	return &msgs.ProtoMsg{
		MessageID: input.MessageID,
		Message:   message,
	}
}

// MarshalPListMsg 序列化測試訊息到PList訊息
func MarshalPListMsg(input []TestMsg) *msgs.PListMsg {
	result := &msgs.PListMsg{}

	for _, itor := range input {
		message, _ := anypb.New(itor.Message)
		result.Messages = append(result.Messages, &msgs.PListUnit{
			MessageID: itor.MessageID,
			Message:   message,
		})
	} // for

	return result
}
