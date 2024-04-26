package helps

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ProtoAny 將proto的any轉換為指定物件
func ProtoAny[T any](input *anypb.Any) (output *T, err error) {
	if input == nil {
		return nil, fmt.Errorf("proto any: input nil")
	} // if

	temp, err := input.UnmarshalNew()

	if err != nil {
		return nil, fmt.Errorf("proto any: %w", err)
	} // if

	output, ok := any(temp).(*T)

	if ok == false {
		return nil, fmt.Errorf("proto any: input type")
	} // if

	return output, nil
}

// ProtoJson 將proto轉為json字串
func ProtoJson(input proto.Message) string {
	result, _ := protojson.Marshal(input)
	return string(result)
}
