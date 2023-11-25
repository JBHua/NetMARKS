// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.4
// source: water.proto

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

	Quantity     uint64 `protobuf:"varint,1,opt,name=quantity,proto3" json:"quantity,omitempty"`
	ResponseSize string `protobuf:"bytes,2,opt,name=response_size,json=responseSize,proto3" json:"response_size,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_water_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_water_proto_msgTypes[0]
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
	return file_water_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetQuantity() uint64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Request) GetResponseSize() string {
	if x != nil {
		return x.ResponseSize
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Quantity uint64    `protobuf:"varint,1,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Items    []*Single `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_water_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_water_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_water_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetQuantity() uint64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Response) GetItems() []*Single {
	if x != nil {
		return x.Items
	}
	return nil
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
		mi := &file_water_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Single) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Single) ProtoMessage() {}

func (x *Single) ProtoReflect() protoreflect.Message {
	mi := &file_water_proto_msgTypes[2]
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
	return file_water_proto_rawDescGZIP(), []int{2}
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

var File_water_proto protoreflect.FileDescriptor

var file_water_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x77, 0x61, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6e,
	0x65, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x5f, 0x67, 0x72, 0x61, 0x69, 0x6e, 0x22, 0x4a, 0x0a,
	0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x54, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x12, 0x2c, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x5f, 0x67, 0x72, 0x61, 0x69,
	0x6e, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22,
	0x41, 0x0a, 0x06, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x61, 0x6e,
	0x64, 0x6f, 0x6d, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x32, 0x47, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x07, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x12, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x72, 0x6b,
	0x73, 0x5f, 0x67, 0x72, 0x61, 0x69, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x6e, 0x65, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x5f, 0x67, 0x72, 0x61, 0x69, 0x6e,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1f, 0x5a, 0x1d, 0x4e,
	0x65, 0x74, 0x4d, 0x41, 0x52, 0x4b, 0x53, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x77, 0x61, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_water_proto_rawDescOnce sync.Once
	file_water_proto_rawDescData = file_water_proto_rawDesc
)

func file_water_proto_rawDescGZIP() []byte {
	file_water_proto_rawDescOnce.Do(func() {
		file_water_proto_rawDescData = protoimpl.X.CompressGZIP(file_water_proto_rawDescData)
	})
	return file_water_proto_rawDescData
}

var file_water_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_water_proto_goTypes = []interface{}{
	(*Request)(nil),  // 0: netmarks_grain.Request
	(*Response)(nil), // 1: netmarks_grain.Response
	(*Single)(nil),   // 2: netmarks_grain.Single
}
var file_water_proto_depIdxs = []int32{
	2, // 0: netmarks_grain.Response.items:type_name -> netmarks_grain.Single
	0, // 1: netmarks_grain.Water.Produce:input_type -> netmarks_grain.Request
	1, // 2: netmarks_grain.Water.Produce:output_type -> netmarks_grain.Response
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_water_proto_init() }
func file_water_proto_init() {
	if File_water_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_water_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_water_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_water_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_water_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_water_proto_goTypes,
		DependencyIndexes: file_water_proto_depIdxs,
		MessageInfos:      file_water_proto_msgTypes,
	}.Build()
	File_water_proto = out.File
	file_water_proto_rawDesc = nil
	file_water_proto_goTypes = nil
	file_water_proto_depIdxs = nil
}
