package trials

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// ProtoEqual 訊息是否符合
func ProtoEqual(expected, actual proto.Message, option ...cmp.Option) bool {
	if cmp.Equal(expected, actual, append(option, protocmp.Transform())...) == false {
		fmt.Printf("expected: %v\n", protojson.Format(expected))
		fmt.Printf("actual: %v\n", protojson.Format(actual))
		return false
	} // if

	return true
}

// ProtoContains 訊息列表是否有符合訊息
func ProtoContains[T any](expected proto.Message, actual []T, option ...cmp.Option) bool {
	builder := &strings.Builder{}

	for _, itor := range actual {
		message, ok := any(itor).(proto.Message)

		if ok == false {
			_, _ = fmt.Fprintf(builder, "\n<not proto message>,")
			continue
		} // if

		if cmp.Equal(expected, message, append(option, protocmp.Transform())...) == false {
			_, _ = fmt.Fprintf(builder, "\n%v,", protojson.Format(message))
			continue
		} // if

		return true
	} // for

	fmt.Printf("expected: %v\n", protojson.Format(expected))
	fmt.Printf("actual: %v\n", builder.String())
	return false
}

// ProtoTypeExist 訊息列表是否有指定類型
func ProtoTypeExist(expected any, actual []proto.Message) bool {
	for _, itor := range actual {
		if reflect.TypeOf(itor) == reflect.TypeOf(expected) {
			return true
		} // if
	} // for

	return false
}
