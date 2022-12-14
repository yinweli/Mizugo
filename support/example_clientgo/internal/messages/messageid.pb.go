// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: messageid.proto

package messages

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 訊息編號
type MsgID int32

const (
	MsgID_Unknown MsgID = 0 // 不明/錯誤訊息編號, 此編號不可使用
	MsgID_EchoReq MsgID = 1 // 要求回音(用簡單封包做)
	MsgID_EchoRes MsgID = 2 // 回應回音(用簡單封包做)
	MsgID_KeyReq  MsgID = 3 // 要求密鑰
	MsgID_KeyRes  MsgID = 4 // 回應密鑰
	MsgID_PingReq MsgID = 5 // 要求Ping
	MsgID_PingRes MsgID = 6 // 回應Ping
)

// Enum value maps for MsgID.
var (
	MsgID_name = map[int32]string{
		0: "Unknown",
		1: "EchoReq",
		2: "EchoRes",
		3: "KeyReq",
		4: "KeyRes",
		5: "PingReq",
		6: "PingRes",
	}
	MsgID_value = map[string]int32{
		"Unknown": 0,
		"EchoReq": 1,
		"EchoRes": 2,
		"KeyReq":  3,
		"KeyRes":  4,
		"PingReq": 5,
		"PingRes": 6,
	}
)

func (x MsgID) Enum() *MsgID {
	p := new(MsgID)
	*p = x
	return p
}

func (x MsgID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgID) Descriptor() protoreflect.EnumDescriptor {
	return file_messageid_proto_enumTypes[0].Descriptor()
}

func (MsgID) Type() protoreflect.EnumType {
	return &file_messageid_proto_enumTypes[0]
}

func (x MsgID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgID.Descriptor instead.
func (MsgID) EnumDescriptor() ([]byte, []int) {
	return file_messageid_proto_rawDescGZIP(), []int{0}
}

var File_messageid_proto protoreflect.FileDescriptor

var file_messageid_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2a, 0x60, 0x0a, 0x05, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e,
	0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x63, 0x68, 0x6f, 0x52,
	0x65, 0x71, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x10,
	0x02, 0x12, 0x0a, 0x0a, 0x06, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x10, 0x03, 0x12, 0x0a, 0x0a,
	0x06, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x10, 0x06, 0x42, 0x14, 0x5a, 0x12, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_messageid_proto_rawDescOnce sync.Once
	file_messageid_proto_rawDescData = file_messageid_proto_rawDesc
)

func file_messageid_proto_rawDescGZIP() []byte {
	file_messageid_proto_rawDescOnce.Do(func() {
		file_messageid_proto_rawDescData = protoimpl.X.CompressGZIP(file_messageid_proto_rawDescData)
	})
	return file_messageid_proto_rawDescData
}

var file_messageid_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_messageid_proto_goTypes = []interface{}{
	(MsgID)(0), // 0: MsgID
}
var file_messageid_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_messageid_proto_init() }
func file_messageid_proto_init() {
	if File_messageid_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messageid_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messageid_proto_goTypes,
		DependencyIndexes: file_messageid_proto_depIdxs,
		EnumInfos:         file_messageid_proto_enumTypes,
	}.Build()
	File_messageid_proto = out.File
	file_messageid_proto_rawDesc = nil
	file_messageid_proto_goTypes = nil
	file_messageid_proto_depIdxs = nil
}
