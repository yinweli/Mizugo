package trials

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// ProtoEqual 比對訊息是否一致
func ProtoEqual(expected, actual proto.Message) bool {
	if proto.Equal(expected, actual) == false {
		fmt.Printf("expected: %v\n", expected)
		fmt.Printf("actual: %v\n", actual)
		return false
	} // if

	return true
}
