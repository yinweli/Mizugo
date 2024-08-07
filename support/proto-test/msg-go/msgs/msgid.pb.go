// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: msgid.proto

package msgs

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
	MsgID_Unknown MsgID = 0  // 不明/錯誤訊息編號, 此編號不可使用
	MsgID_JsonQ   MsgID = 1  // 要求Json
	MsgID_JsonA   MsgID = 2  // 回應Json
	MsgID_ProtoQ  MsgID = 3  // 要求Proto
	MsgID_ProtoA  MsgID = 4  // 回應Proto
	MsgID_RavenQ  MsgID = 5  // 要求Raven
	MsgID_RavenA  MsgID = 6  // 回應Raven
	MsgID_LoginQ  MsgID = 7  // 要求登入(使用Json處理器)
	MsgID_LoginA  MsgID = 8  // 回應登入(使用Json處理器)
	MsgID_UpdateQ MsgID = 9  // 要求更新(使用Json處理器)
	MsgID_UpdateA MsgID = 10 // 回應更新(使用Json處理器)
)

// Enum value maps for MsgID.
var (
	MsgID_name = map[int32]string{
		0:  "Unknown",
		1:  "JsonQ",
		2:  "JsonA",
		3:  "ProtoQ",
		4:  "ProtoA",
		5:  "RavenQ",
		6:  "RavenA",
		7:  "LoginQ",
		8:  "LoginA",
		9:  "UpdateQ",
		10: "UpdateA",
	}
	MsgID_value = map[string]int32{
		"Unknown": 0,
		"JsonQ":   1,
		"JsonA":   2,
		"ProtoQ":  3,
		"ProtoA":  4,
		"RavenQ":  5,
		"RavenA":  6,
		"LoginQ":  7,
		"LoginA":  8,
		"UpdateQ": 9,
		"UpdateA": 10,
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
	return file_msgid_proto_enumTypes[0].Descriptor()
}

func (MsgID) Type() protoreflect.EnumType {
	return &file_msgid_proto_enumTypes[0]
}

func (x MsgID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgID.Descriptor instead.
func (MsgID) EnumDescriptor() ([]byte, []int) {
	return file_msgid_proto_rawDescGZIP(), []int{0}
}

// 錯誤編號
type ErrID int32

const (
	ErrID_Success        ErrID = 0 // 成功
	ErrID_JsonUnmarshal  ErrID = 1 // Json反序列化失敗
	ErrID_ProtoUnmarshal ErrID = 2 // Proto反序列化失敗
	ErrID_RavenUnmarshal ErrID = 3 // Raven反序列化失敗
	ErrID_SubmitFailed   ErrID = 4 // 資料庫執行失敗
	ErrID_TokenNotMatch  ErrID = 5 // Token不匹配
)

// Enum value maps for ErrID.
var (
	ErrID_name = map[int32]string{
		0: "Success",
		1: "JsonUnmarshal",
		2: "ProtoUnmarshal",
		3: "RavenUnmarshal",
		4: "SubmitFailed",
		5: "TokenNotMatch",
	}
	ErrID_value = map[string]int32{
		"Success":        0,
		"JsonUnmarshal":  1,
		"ProtoUnmarshal": 2,
		"RavenUnmarshal": 3,
		"SubmitFailed":   4,
		"TokenNotMatch":  5,
	}
)

func (x ErrID) Enum() *ErrID {
	p := new(ErrID)
	*p = x
	return p
}

func (x ErrID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrID) Descriptor() protoreflect.EnumDescriptor {
	return file_msgid_proto_enumTypes[1].Descriptor()
}

func (ErrID) Type() protoreflect.EnumType {
	return &file_msgid_proto_enumTypes[1]
}

func (x ErrID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrID.Descriptor instead.
func (ErrID) EnumDescriptor() ([]byte, []int) {
	return file_msgid_proto_rawDescGZIP(), []int{1}
}

var File_msgid_proto protoreflect.FileDescriptor

var file_msgid_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x73, 0x67, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x8c, 0x01,
	0x0a, 0x05, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f,
	0x77, 0x6e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4a, 0x73, 0x6f, 0x6e, 0x51, 0x10, 0x01, 0x12,
	0x09, 0x0a, 0x05, 0x4a, 0x73, 0x6f, 0x6e, 0x41, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x51, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x41,
	0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x61, 0x76, 0x65, 0x6e, 0x51, 0x10, 0x05, 0x12, 0x0a,
	0x0a, 0x06, 0x52, 0x61, 0x76, 0x65, 0x6e, 0x41, 0x10, 0x06, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x51, 0x10, 0x07, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41,
	0x10, 0x08, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x51, 0x10, 0x09, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x10, 0x0a, 0x2a, 0x74, 0x0a, 0x05,
	0x45, 0x72, 0x72, 0x49, 0x44, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4a, 0x73, 0x6f, 0x6e, 0x55, 0x6e, 0x6d, 0x61, 0x72, 0x73,
	0x68, 0x61, 0x6c, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x55, 0x6e,
	0x6d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x61, 0x76,
	0x65, 0x6e, 0x55, 0x6e, 0x6d, 0x61, 0x72, 0x73, 0x68, 0x61, 0x6c, 0x10, 0x03, 0x12, 0x10, 0x0a,
	0x0c, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x04, 0x12,
	0x11, 0x0a, 0x0d, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x4e, 0x6f, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x10, 0x05, 0x42, 0x0c, 0x5a, 0x0a, 0x2f, 0x6d, 0x73, 0x67, 0x73, 0x3b, 0x6d, 0x73, 0x67, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msgid_proto_rawDescOnce sync.Once
	file_msgid_proto_rawDescData = file_msgid_proto_rawDesc
)

func file_msgid_proto_rawDescGZIP() []byte {
	file_msgid_proto_rawDescOnce.Do(func() {
		file_msgid_proto_rawDescData = protoimpl.X.CompressGZIP(file_msgid_proto_rawDescData)
	})
	return file_msgid_proto_rawDescData
}

var file_msgid_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_msgid_proto_goTypes = []interface{}{
	(MsgID)(0), // 0: MsgID
	(ErrID)(0), // 1: ErrID
}
var file_msgid_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_msgid_proto_init() }
func file_msgid_proto_init() {
	if File_msgid_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msgid_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msgid_proto_goTypes,
		DependencyIndexes: file_msgid_proto_depIdxs,
		EnumInfos:         file_msgid_proto_enumTypes,
	}.Build()
	File_msgid_proto = out.File
	file_msgid_proto_rawDesc = nil
	file_msgid_proto_goTypes = nil
	file_msgid_proto_depIdxs = nil
}
