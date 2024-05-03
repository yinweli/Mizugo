package helps

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ToProtoAny 將proto轉換為any
func ToProtoAny(input ...proto.Message) (output []*anypb.Any, err error) {
	output = []*anypb.Any{}
	opt := proto.MarshalOptions{}

	for _, itor := range input {
		temp := &anypb.Any{}

		if err = anypb.MarshalFrom(temp, itor, opt); err != nil {
			return nil, Err(err)
		} // if

		output = append(output, temp)
	} // for

	return output, nil
}

// FromProtoAny 將any轉換為指定物件
func FromProtoAny[T any](input *anypb.Any) (output *T, err error) {
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

// ProtoString 將proto轉為字串
func ProtoString(input proto.Message) string {
	result, _ := protojson.Marshal(input)
	return string(result)
}
