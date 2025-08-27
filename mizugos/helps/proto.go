package helps

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ToProtoAny 將一或多個 Proto 訊息轉為 []*anypb.Any
func ToProtoAny(input ...proto.Message) (output []*anypb.Any, err error) {
	output = []*anypb.Any{}

	for _, itor := range input {
		if itor == nil {
			return nil, fmt.Errorf("to proto any: input nil")
		} // if

		message, err := anypb.New(itor)

		if err != nil {
			return nil, fmt.Errorf("to proto any: %w", err)
		} // if

		output = append(output, message)
	} // for

	return output, nil
}

// FromProtoAny 將 any 反序列化為指定型別 T 的指標, T 應為 Proto 訊息
func FromProtoAny[T proto.Message](input *anypb.Any) (output T, err error) {
	if input == nil {
		return output, fmt.Errorf("from proto any: input nil")
	} // if

	message, err := input.UnmarshalNew()

	if err != nil {
		return output, fmt.Errorf("from proto any: %w", err)
	} // if

	output, ok := message.(T)

	if ok == false {
		return output, fmt.Errorf("from proto any: type mismatch")
	} // if

	return output, nil
}

// ProtoString 將 Proto 轉為字串
func ProtoString(input proto.Message) string {
	result, _ := protojson.Marshal(input)
	return string(result)
}
