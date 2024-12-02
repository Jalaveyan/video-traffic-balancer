// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.29.0--rc3
// source: proto/balancer.proto

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

type RedirectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Video string `protobuf:"bytes,1,opt,name=video,proto3" json:"video,omitempty"`
}

func (x *RedirectRequest) Reset() {
	*x = RedirectRequest{}
	mi := &file_proto_balancer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedirectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectRequest) ProtoMessage() {}

func (x *RedirectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectRequest.ProtoReflect.Descriptor instead.
func (*RedirectRequest) Descriptor() ([]byte, []int) {
	return file_proto_balancer_proto_rawDescGZIP(), []int{0}
}

func (x *RedirectRequest) GetVideo() string {
	if x != nil {
		return x.Video
	}
	return ""
}

type RedirectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetUrl string `protobuf:"bytes,1,opt,name=targetUrl,proto3" json:"targetUrl,omitempty"`
}

func (x *RedirectResponse) Reset() {
	*x = RedirectResponse{}
	mi := &file_proto_balancer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RedirectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedirectResponse) ProtoMessage() {}

func (x *RedirectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_balancer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedirectResponse.ProtoReflect.Descriptor instead.
func (*RedirectResponse) Descriptor() ([]byte, []int) {
	return file_proto_balancer_proto_rawDescGZIP(), []int{1}
}

func (x *RedirectResponse) GetTargetUrl() string {
	if x != nil {
		return x.TargetUrl
	}
	return ""
}

var File_proto_balancer_proto protoreflect.FileDescriptor

var file_proto_balancer_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x62, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x22, 0x27, 0x0a, 0x0f, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x22, 0x30, 0x0a,
	0x10, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x72, 0x6c, 0x32,
	0x55, 0x0a, 0x08, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x12, 0x49, 0x0a, 0x08, 0x52,
	0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x1d, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x62, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_balancer_proto_rawDescOnce sync.Once
	file_proto_balancer_proto_rawDescData = file_proto_balancer_proto_rawDesc
)

func file_proto_balancer_proto_rawDescGZIP() []byte {
	file_proto_balancer_proto_rawDescOnce.Do(func() {
		file_proto_balancer_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_balancer_proto_rawDescData)
	})
	return file_proto_balancer_proto_rawDescData
}

var file_proto_balancer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_balancer_proto_goTypes = []any{
	(*RedirectRequest)(nil),  // 0: videobalance.RedirectRequest
	(*RedirectResponse)(nil), // 1: videobalance.RedirectResponse
}
var file_proto_balancer_proto_depIdxs = []int32{
	0, // 0: videobalance.Balancer.Redirect:input_type -> videobalance.RedirectRequest
	1, // 1: videobalance.Balancer.Redirect:output_type -> videobalance.RedirectResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_balancer_proto_init() }
func file_proto_balancer_proto_init() {
	if File_proto_balancer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_balancer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_balancer_proto_goTypes,
		DependencyIndexes: file_proto_balancer_proto_depIdxs,
		MessageInfos:      file_proto_balancer_proto_msgTypes,
	}.Build()
	File_proto_balancer_proto = out.File
	file_proto_balancer_proto_rawDesc = nil
	file_proto_balancer_proto_goTypes = nil
	file_proto_balancer_proto_depIdxs = nil
}
