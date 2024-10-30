package trials

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// ProtoEqual 比對訊息是否一致
func ProtoEqual(expected, actual proto.Message, option ...cmp.Option) bool {
	if cmp.Equal(expected, actual, append(option, protocmp.Transform())...) == false {
		fmt.Printf("expected: %v\n", protojson.Format(expected))
		fmt.Printf("actual: %v\n", protojson.Format(actual))
		return false
	} // if

	return true
}

// ProtoContains 訊息列表中是否有指定類型
func ProtoContains(source []proto.Message, expected any) bool {
	for _, itor := range source {
		if reflect.TypeOf(itor) == reflect.TypeOf(expected) {
			return true
		} // if
	} // for

	return false
}
