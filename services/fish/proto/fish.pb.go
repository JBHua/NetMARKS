// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.4
// source: fish.proto

package proto

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fish_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_fish_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_fish_proto_rawDescGZIP(), []int{0}
}

type Single struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RandomMetadata string `protobuf:"bytes,2,opt,name=random_metadata,json=randomMetadata,proto3" json:"random_metadata,omitempty"`
}

func (x *Single) Reset() {
	*x = Single{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fish_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Single) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Single) ProtoMessage() {}

func (x *Single) ProtoReflect() protoreflect.Message {
	mi := &file_fish_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Single.ProtoReflect.Descriptor instead.
func (*Single) Descriptor() ([]byte, []int) {
	return file_fish_proto_rawDescGZIP(), []int{1}
}

func (x *Single) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Single) GetRandomMetadata() string {
	if x != nil {
		return x.RandomMetadata
	}
	return ""
}

var File_fish_proto protoreflect.FileDescriptor

var file_fish_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x66, 0x69, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x6e, 0x65,
	0x74, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x5f, 0x66, 0x69, 0x73, 0x68, 0x22, 0x09, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x41, 0x0a, 0x06, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x27, 0x0a, 0x0f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x61, 0x6e, 0x64, 0x6f,
	0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x32, 0x42, 0x0a, 0x04, 0x46, 0x69, 0x73,
	0x68, 0x12, 0x3a, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x6e,
	0x65, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x5f, 0x66, 0x69, 0x73, 0x68, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x5f,
	0x66, 0x69, 0x73, 0x68, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x22, 0x00, 0x42, 0x2f, 0x5a,
	0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4a, 0x42, 0x48, 0x75,
	0x61, 0x2f, 0x4e, 0x65, 0x74, 0x4d, 0x41, 0x52, 0x4b, 0x53, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x66, 0x69, 0x73, 0x68, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fish_proto_rawDescOnce sync.Once
	file_fish_proto_rawDescData = file_fish_proto_rawDesc
)

func file_fish_proto_rawDescGZIP() []byte {
	file_fish_proto_rawDescOnce.Do(func() {
		file_fish_proto_rawDescData = protoimpl.X.CompressGZIP(file_fish_proto_rawDescData)
	})
	return file_fish_proto_rawDescData
}

var file_fish_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_fish_proto_goTypes = []interface{}{
	(*Request)(nil), // 0: netmarks_fish.Request
	(*Single)(nil),  // 1: netmarks_fish.Single
}
var file_fish_proto_depIdxs = []int32{
	0, // 0: netmarks_fish.Fish.Produce:input_type -> netmarks_fish.Request
	1, // 1: netmarks_fish.Fish.Produce:output_type -> netmarks_fish.Single
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fish_proto_init() }
func file_fish_proto_init() {
	if File_fish_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fish_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_fish_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Single); i {
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
			RawDescriptor: file_fish_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fish_proto_goTypes,
		DependencyIndexes: file_fish_proto_depIdxs,
		MessageInfos:      file_fish_proto_msgTypes,
	}.Build()
	File_fish_proto = out.File
	file_fish_proto_rawDesc = nil
	file_fish_proto_goTypes = nil
	file_fish_proto_depIdxs = nil
}
