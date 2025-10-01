package trials

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// ProtoEqual 比對訊息是否符合預期
func ProtoEqual(expected, actual any, option ...cmp.Option) bool {
	option = append([]cmp.Option{protocmp.Transform()}, option...) // 慣例是把 protocmp.Transform() 放在最前面

	if cmp.Equal(expected, actual, option...) == false {
		fmt.Println("proto not equal:")
		fmt.Println("  expected:")

		if m, ok := expected.(proto.Message); ok && m != nil {
			fmt.Printf("    %v\n", protoJSON.Format(m))
		} else {
			fmt.Printf("    not proto message\n")
		} // if

		fmt.Println("  actual:")

		if m, ok := actual.(proto.Message); ok && m != nil {
			fmt.Printf("    %v\n", protoJSON.Format(m))
		} else {
			fmt.Printf("    not proto message\n")
		} // if

		return false
	} // if

	return true
}

// ProtoListEqual 比對訊息列表是否符合預期
func ProtoListEqual[T any](expected, actual []T, option ...cmp.Option) bool {
	option = append([]cmp.Option{protocmp.Transform()}, option...) // 慣例是把 protocmp.Transform() 放在最前面

	if cmp.Equal(expected, actual, option...) == false {
		fmt.Println("proto not equal:")
		fmt.Println("  expected:")

		for _, itor := range expected {
			if m, ok := any(itor).(proto.Message); ok && m != nil {
				fmt.Printf("    %v\n", protoJSON.Format(m))
			} else {
				fmt.Printf("    not proto message\n")
			} // if
		} // for

		fmt.Println("  actual:")

		for _, itor := range actual {
			if m, ok := any(itor).(proto.Message); ok && m != nil {
				fmt.Printf("    %v\n", protoJSON.Format(m))
			} else {
				fmt.Printf("    not proto message\n")
			} // if
		} // for

		return false
	} // if

	return true
}

// ProtoListMatch 比對訊息列表是否符合預期, 無關順序
func ProtoListMatch[T any](expected, actual []T, option ...cmp.Option) bool {
	option = append([]cmp.Option{protocmp.Transform()}, option...) // 慣例是把 protocmp.Transform() 放在最前面
	match := map[int]bool{}

	for _, e := range expected {
		for i, a := range actual {
			if match[i] == false && cmp.Equal(e, a, option...) {
				match[i] = true
				break
			} // if
		} // for
	} // for

	if len(expected) != len(actual) || len(expected) != len(match) {
		fmt.Println("proto not equal:")
		fmt.Println("  expected:")

		for _, itor := range expected {
			if m, ok := any(itor).(proto.Message); ok && m != nil {
				fmt.Printf("    %v\n", protoJSON.Format(m))
			} else {
				fmt.Printf("    not proto message\n")
			} // if
		} // for

		fmt.Println("  actual:")

		for _, itor := range actual {
			if m, ok := any(itor).(proto.Message); ok && m != nil {
				fmt.Printf("    %v\n", protoJSON.Format(m))
			} else {
				fmt.Printf("    not proto message\n")
			} // if
		} // for

		return false
	} // if

	return true
}

// ProtoListHasData 檢查訊息列表是否有指定項目
func ProtoListHasData[T any](expected any, actual []T, option ...cmp.Option) bool {
	option = append([]cmp.Option{protocmp.Transform()}, option...) // 慣例是把 protocmp.Transform() 放在最前面

	for _, itor := range actual {
		if cmp.Equal(expected, itor, option...) {
			return true
		} // if
	} // for

	fmt.Println("proto not equal:")
	fmt.Println("  expected:")

	if m, ok := expected.(proto.Message); ok && m != nil {
		fmt.Printf("    %v\n", protoJSON.Format(m))
	} else {
		fmt.Printf("    not proto message\n")
	} // if

	fmt.Println("  actual:")

	for _, itor := range actual {
		if m, ok := any(itor).(proto.Message); ok && m != nil {
			fmt.Printf("    %v\n", protoJSON.Format(m))
		} else {
			fmt.Printf("    not proto message\n")
		} // if
	} // for

	return false
}

// ProtoListHasType 檢查訊息列表是否有指定類型
func ProtoListHasType[T any](expected any, actual []T) bool {
	for _, itor := range actual {
		if reflect.TypeOf(itor) == reflect.TypeOf(expected) {
			return true
		} // if
	} // for

	return false
}

var protoJSON = protojson.MarshalOptions{} // Proto 轉換 JSON 工具
