// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: stackmsg.proto

package msgs

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

// 堆棧訊息資料
type StackMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*StackUnit `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"` // 訊息列表
}

func (x *StackMsg) Reset() {
	*x = StackMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stackmsg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StackMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StackMsg) ProtoMessage() {}

func (x *StackMsg) ProtoReflect() protoreflect.Message {
	mi := &file_stackmsg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StackMsg.ProtoReflect.Descriptor instead.
func (*StackMsg) Descriptor() ([]byte, []int) {
	return file_stackmsg_proto_rawDescGZIP(), []int{0}
}

func (x *StackMsg) GetMessages() []*StackUnit {
	if x != nil {
		return x.Messages
	}
	return nil
}

// 堆棧訊息單元資料
type StackUnit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageID int32      `protobuf:"varint,1,opt,name=messageID,proto3" json:"messageID,omitempty"` // 訊息編號, 設置為int32以跟proto的列舉類型統一
	Message   *anypb.Any `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`      // 訊息資料
}

func (x *StackUnit) Reset() {
	*x = StackUnit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stackmsg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StackUnit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StackUnit) ProtoMessage() {}

func (x *StackUnit) ProtoReflect() protoreflect.Message {
	mi := &file_stackmsg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StackUnit.ProtoReflect.Descriptor instead.
func (*StackUnit) Descriptor() ([]byte, []int) {
	return file_stackmsg_proto_rawDescGZIP(), []int{1}
}

func (x *StackUnit) GetMessageID() int32 {
	if x != nil {
		return x.MessageID
	}
	return 0
}

func (x *StackUnit) GetMessage() *anypb.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

// 堆棧訊息測試資料
type StackTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"` // 測試字串
}

func (x *StackTest) Reset() {
	*x = StackTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stackmsg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StackTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StackTest) ProtoMessage() {}

func (x *StackTest) ProtoReflect() protoreflect.Message {
	mi := &file_stackmsg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StackTest.ProtoReflect.Descriptor instead.
func (*StackTest) Descriptor() ([]byte, []int) {
	return file_stackmsg_proto_rawDescGZIP(), []int{2}
}

func (x *StackTest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_stackmsg_proto protoreflect.FileDescriptor

var file_stackmsg_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x08, 0x53,
	0x74, 0x61, 0x63, 0x6b, 0x4d, 0x73, 0x67, 0x12, 0x26, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x53, 0x74, 0x61, 0x63,
	0x6b, 0x55, 0x6e, 0x69, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22,
	0x59, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x63, 0x6b, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1f, 0x0a, 0x09, 0x53, 0x74,
	0x61, 0x63, 0x6b, 0x54, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x42, 0x0c, 0x5a, 0x0a, 0x2f,
	0x6d, 0x73, 0x67, 0x73, 0x3b, 0x6d, 0x73, 0x67, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_stackmsg_proto_rawDescOnce sync.Once
	file_stackmsg_proto_rawDescData = file_stackmsg_proto_rawDesc
)

func file_stackmsg_proto_rawDescGZIP() []byte {
	file_stackmsg_proto_rawDescOnce.Do(func() {
		file_stackmsg_proto_rawDescData = protoimpl.X.CompressGZIP(file_stackmsg_proto_rawDescData)
	})
	return file_stackmsg_proto_rawDescData
}

var file_stackmsg_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_stackmsg_proto_goTypes = []interface{}{
	(*StackMsg)(nil),  // 0: StackMsg
	(*StackUnit)(nil), // 1: StackUnit
	(*StackTest)(nil), // 2: StackTest
	(*anypb.Any)(nil), // 3: google.protobuf.Any
}
var file_stackmsg_proto_depIdxs = []int32{
	1, // 0: StackMsg.messages:type_name -> StackUnit
	3, // 1: StackUnit.message:type_name -> google.protobuf.Any
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_stackmsg_proto_init() }
func file_stackmsg_proto_init() {
	if File_stackmsg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stackmsg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StackMsg); i {
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
		file_stackmsg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StackUnit); i {
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
		file_stackmsg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StackTest); i {
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
			RawDescriptor: file_stackmsg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_stackmsg_proto_goTypes,
		DependencyIndexes: file_stackmsg_proto_depIdxs,
		MessageInfos:      file_stackmsg_proto_msgTypes,
	}.Build()
	File_stackmsg_proto = out.File
	file_stackmsg_proto_rawDesc = nil
	file_stackmsg_proto_goTypes = nil
	file_stackmsg_proto_depIdxs = nil
}