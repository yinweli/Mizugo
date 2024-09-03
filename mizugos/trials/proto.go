package trials

import (
	"fmt"

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
