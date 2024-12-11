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
func ProtoEqual(expected, actual any, option ...cmp.Option) bool {
	if cmp.Equal(expected, actual, append(option, protocmp.Transform())...) == false {
		if message, ok := expected.(proto.Message); ok {
			fmt.Printf("expected: %v\n", protojson.Format(message))
		} else {
			fmt.Printf("expected: <not proto message>\n")
		} // if

		if message, ok := actual.(proto.Message); ok {
			fmt.Printf("actual: %v\n", protojson.Format(message))
		} else {
			fmt.Printf("actual: <not proto message>\n")
		} // if

		return false
	} // if

	return true
}

// ProtoListEqual 訊息列表是否符合
func ProtoListEqual[T any](expected, actual []T, option ...cmp.Option) bool {
	if cmp.Equal(expected, actual, append(option, protocmp.Transform())...) == false {
		builder := &strings.Builder{}

		for _, itor := range expected {
			if message, ok := any(itor).(proto.Message); ok {
				_, _ = fmt.Fprintf(builder, "\n%v,", indent(protojson.Format(message), "  "))
			} else {
				_, _ = fmt.Fprintf(builder, "\n  <not proto message>,")
			} // if
		} // for

		fmt.Printf("expected: %v\n", builder.String())
		builder.Reset()

		for _, itor := range actual {
			if message, ok := any(itor).(proto.Message); ok {
				_, _ = fmt.Fprintf(builder, "\n%v,", indent(protojson.Format(message), "  "))
			} else {
				_, _ = fmt.Fprintf(builder, "\n  <not proto message>,")
			} // if
		} // for

		fmt.Printf("actual: %v\n", builder.String())
		return false
	} // if

	return true
}

// ProtoListContain 訊息列表是否包含指定訊息
func ProtoListContain[T any](expected any, actual []T, option ...cmp.Option) bool {
	for _, itor := range actual {
		if cmp.Equal(expected, itor, append(option, protocmp.Transform())...) {
			return true
		} // if
	} // for

	if message, ok := expected.(proto.Message); ok {
		fmt.Printf("expected: %v\n", protojson.Format(message))
	} else {
		fmt.Printf("expected: <not proto message>\n")
	} // if

	builder := &strings.Builder{}

	for _, itor := range actual {
		if message, ok := any(itor).(proto.Message); ok {
			_, _ = fmt.Fprintf(builder, "\n%v,", indent(protojson.Format(message), "  "))
		} else {
			_, _ = fmt.Fprintf(builder, "\n  <not proto message>,")
		} // if
	} // for

	fmt.Printf("actual: %v\n", builder.String())
	return false
}

// ProtoListExist 訊息列表是否包含指定類型
func ProtoListExist(expected any, actual []proto.Message) bool {
	for _, itor := range actual {
		if reflect.TypeOf(itor) == reflect.TypeOf(expected) {
			return true
		} // if
	} // for

	return false
}

// indent 填加縮排到多行字串
func indent(input, indent string) string {
	line := strings.Split(input, "\n")

	for i, itor := range line {
		line[i] = indent + itor
	} // for

	return strings.Join(line, "\n")
}
