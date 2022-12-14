// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: protodesmsg.proto

package procs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ProtoDes訊息資料
type ProtoDesMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageID int32      `protobuf:"varint,1,opt,name=messageID,proto3" json:"messageID,omitempty"` // 訊息編號, 設置為int32以跟proto的列舉類型統一
	Message   *anypb.Any `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`      // 訊息資料
}

func (x *ProtoDesMsg) Reset() {
	*x = ProtoDesMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protodesmsg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoDesMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoDesMsg) ProtoMessage() {}

func (x *ProtoDesMsg) ProtoReflect() protoreflect.Message {
	mi := &file_protodesmsg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoDesMsg.ProtoReflect.Descriptor instead.
func (*ProtoDesMsg) Descriptor() ([]byte, []int) {
	return file_protodesmsg_proto_rawDescGZIP(), []int{0}
}

func (x *ProtoDesMsg) GetMessageID() int32 {
	if x != nil {
		return x.MessageID
	}
	return 0
}

func (x *ProtoDesMsg) GetMessage() *anypb.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

// ProtoDes訊息測試用資料
type ProtoDesMsgTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"` // 訊息內容
}

func (x *ProtoDesMsgTest) Reset() {
	*x = ProtoDesMsgTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protodesmsg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtoDesMsgTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoDesMsgTest) ProtoMessage() {}

func (x *ProtoDesMsgTest) ProtoReflect() protoreflect.Message {
	mi := &file_protodesmsg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoDesMsgTest.ProtoReflect.Descriptor instead.
func (*ProtoDesMsgTest) Descriptor() ([]byte, []int) {
	return file_protodesmsg_proto_rawDescGZIP(), []int{1}
}

func (x *ProtoDesMsgTest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_protodesmsg_proto protoreflect.FileDescriptor

var file_protodesmsg_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x64, 0x65, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5b,
	0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x44, 0x65, 0x73, 0x4d, 0x73, 0x67, 0x12, 0x1c, 0x0a,
	0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x2e, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2b, 0x0a, 0x0f, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x44, 0x65, 0x73, 0x4d, 0x73, 0x67, 0x54, 0x65, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x14, 0x5a, 0x12, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x64, 0x65, 0x73, 0x6d, 0x73, 0x67, 0x3b, 0x70, 0x72, 0x6f, 0x63, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protodesmsg_proto_rawDescOnce sync.Once
	file_protodesmsg_proto_rawDescData = file_protodesmsg_proto_rawDesc
)

func file_protodesmsg_proto_rawDescGZIP() []byte {
	file_protodesmsg_proto_rawDescOnce.Do(func() {
		file_protodesmsg_proto_rawDescData = protoimpl.X.CompressGZIP(file_protodesmsg_proto_rawDescData)
	})
	return file_protodesmsg_proto_rawDescData
}

var file_protodesmsg_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protodesmsg_proto_goTypes = []interface{}{
	(*ProtoDesMsg)(nil),     // 0: ProtoDesMsg
	(*ProtoDesMsgTest)(nil), // 1: ProtoDesMsgTest
	(*anypb.Any)(nil),       // 2: google.protobuf.Any
}
var file_protodesmsg_proto_depIdxs = []int32{
	2, // 0: ProtoDesMsg.message:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protodesmsg_proto_init() }
func file_protodesmsg_proto_init() {
	if File_protodesmsg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protodesmsg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtoDesMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protodesmsg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtoDesMsgTest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protodesmsg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protodesmsg_proto_goTypes,
		DependencyIndexes: file_protodesmsg_proto_depIdxs,
		MessageInfos:      file_protodesmsg_proto_msgTypes,
	}.Build()
	File_protodesmsg_proto = out.File
	file_protodesmsg_proto_rawDesc = nil
	file_protodesmsg_proto_goTypes = nil
	file_protodesmsg_proto_depIdxs = nil
}
