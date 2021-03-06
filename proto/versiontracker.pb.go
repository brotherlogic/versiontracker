// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        (unknown)
// source: versiontracker.proto

package versiontracker

import (
	context "context"
	proto1 "github.com/brotherlogic/buildserver/proto"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type NewVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version *proto1.Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *NewVersionRequest) Reset() {
	*x = NewVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_versiontracker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewVersionRequest) ProtoMessage() {}

func (x *NewVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_versiontracker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewVersionRequest.ProtoReflect.Descriptor instead.
func (*NewVersionRequest) Descriptor() ([]byte, []int) {
	return file_versiontracker_proto_rawDescGZIP(), []int{0}
}

func (x *NewVersionRequest) GetVersion() *proto1.Version {
	if x != nil {
		return x.Version
	}
	return nil
}

type NewVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NewVersionResponse) Reset() {
	*x = NewVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_versiontracker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewVersionResponse) ProtoMessage() {}

func (x *NewVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_versiontracker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewVersionResponse.ProtoReflect.Descriptor instead.
func (*NewVersionResponse) Descriptor() ([]byte, []int) {
	return file_versiontracker_proto_rawDescGZIP(), []int{1}
}

type NewJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version *proto1.Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *NewJobRequest) Reset() {
	*x = NewJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_versiontracker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewJobRequest) ProtoMessage() {}

func (x *NewJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_versiontracker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewJobRequest.ProtoReflect.Descriptor instead.
func (*NewJobRequest) Descriptor() ([]byte, []int) {
	return file_versiontracker_proto_rawDescGZIP(), []int{2}
}

func (x *NewJobRequest) GetVersion() *proto1.Version {
	if x != nil {
		return x.Version
	}
	return nil
}

type NewJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NewJobResponse) Reset() {
	*x = NewJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_versiontracker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewJobResponse) ProtoMessage() {}

func (x *NewJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_versiontracker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewJobResponse.ProtoReflect.Descriptor instead.
func (*NewJobResponse) Descriptor() ([]byte, []int) {
	return file_versiontracker_proto_rawDescGZIP(), []int{3}
}

var File_versiontracker_proto protoreflect.FileDescriptor

var file_versiontracker_proto_rawDesc = []byte{
	0x0a, 0x14, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x1a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x72, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2f,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x11, 0x4e, 0x65, 0x77, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x14, 0x0a, 0x12, 0x4e, 0x65, 0x77, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3f,
	0x0a, 0x0d, 0x4e, 0x65, 0x77, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2e, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22,
	0x10, 0x0a, 0x0e, 0x4e, 0x65, 0x77, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0xb9, 0x01, 0x0a, 0x15, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x72, 0x61,
	0x63, 0x6b, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x4e,
	0x65, 0x77, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x4e, 0x65, 0x77, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x4e, 0x65,
	0x77, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x4e, 0x65, 0x77, 0x4a, 0x6f, 0x62, 0x12, 0x1d, 0x2e, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x4e, 0x65,
	0x77, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x4e, 0x65, 0x77,
	0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_versiontracker_proto_rawDescOnce sync.Once
	file_versiontracker_proto_rawDescData = file_versiontracker_proto_rawDesc
)

func file_versiontracker_proto_rawDescGZIP() []byte {
	file_versiontracker_proto_rawDescOnce.Do(func() {
		file_versiontracker_proto_rawDescData = protoimpl.X.CompressGZIP(file_versiontracker_proto_rawDescData)
	})
	return file_versiontracker_proto_rawDescData
}

var file_versiontracker_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_versiontracker_proto_goTypes = []interface{}{
	(*NewVersionRequest)(nil),  // 0: versiontracker.NewVersionRequest
	(*NewVersionResponse)(nil), // 1: versiontracker.NewVersionResponse
	(*NewJobRequest)(nil),      // 2: versiontracker.NewJobRequest
	(*NewJobResponse)(nil),     // 3: versiontracker.NewJobResponse
	(*proto1.Version)(nil),     // 4: buildserver.Version
}
var file_versiontracker_proto_depIdxs = []int32{
	4, // 0: versiontracker.NewVersionRequest.version:type_name -> buildserver.Version
	4, // 1: versiontracker.NewJobRequest.version:type_name -> buildserver.Version
	0, // 2: versiontracker.VersionTrackerService.NewVersion:input_type -> versiontracker.NewVersionRequest
	2, // 3: versiontracker.VersionTrackerService.NewJob:input_type -> versiontracker.NewJobRequest
	1, // 4: versiontracker.VersionTrackerService.NewVersion:output_type -> versiontracker.NewVersionResponse
	3, // 5: versiontracker.VersionTrackerService.NewJob:output_type -> versiontracker.NewJobResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_versiontracker_proto_init() }
func file_versiontracker_proto_init() {
	if File_versiontracker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_versiontracker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewVersionRequest); i {
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
		file_versiontracker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewVersionResponse); i {
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
		file_versiontracker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewJobRequest); i {
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
		file_versiontracker_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewJobResponse); i {
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
			RawDescriptor: file_versiontracker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_versiontracker_proto_goTypes,
		DependencyIndexes: file_versiontracker_proto_depIdxs,
		MessageInfos:      file_versiontracker_proto_msgTypes,
	}.Build()
	File_versiontracker_proto = out.File
	file_versiontracker_proto_rawDesc = nil
	file_versiontracker_proto_goTypes = nil
	file_versiontracker_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// VersionTrackerServiceClient is the client API for VersionTrackerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VersionTrackerServiceClient interface {
	NewVersion(ctx context.Context, in *NewVersionRequest, opts ...grpc.CallOption) (*NewVersionResponse, error)
	NewJob(ctx context.Context, in *NewJobRequest, opts ...grpc.CallOption) (*NewJobResponse, error)
}

type versionTrackerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVersionTrackerServiceClient(cc grpc.ClientConnInterface) VersionTrackerServiceClient {
	return &versionTrackerServiceClient{cc}
}

func (c *versionTrackerServiceClient) NewVersion(ctx context.Context, in *NewVersionRequest, opts ...grpc.CallOption) (*NewVersionResponse, error) {
	out := new(NewVersionResponse)
	err := c.cc.Invoke(ctx, "/versiontracker.VersionTrackerService/NewVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *versionTrackerServiceClient) NewJob(ctx context.Context, in *NewJobRequest, opts ...grpc.CallOption) (*NewJobResponse, error) {
	out := new(NewJobResponse)
	err := c.cc.Invoke(ctx, "/versiontracker.VersionTrackerService/NewJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VersionTrackerServiceServer is the server API for VersionTrackerService service.
type VersionTrackerServiceServer interface {
	NewVersion(context.Context, *NewVersionRequest) (*NewVersionResponse, error)
	NewJob(context.Context, *NewJobRequest) (*NewJobResponse, error)
}

// UnimplementedVersionTrackerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedVersionTrackerServiceServer struct {
}

func (*UnimplementedVersionTrackerServiceServer) NewVersion(context.Context, *NewVersionRequest) (*NewVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewVersion not implemented")
}
func (*UnimplementedVersionTrackerServiceServer) NewJob(context.Context, *NewJobRequest) (*NewJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewJob not implemented")
}

func RegisterVersionTrackerServiceServer(s *grpc.Server, srv VersionTrackerServiceServer) {
	s.RegisterService(&_VersionTrackerService_serviceDesc, srv)
}

func _VersionTrackerService_NewVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionTrackerServiceServer).NewVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/versiontracker.VersionTrackerService/NewVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionTrackerServiceServer).NewVersion(ctx, req.(*NewVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VersionTrackerService_NewJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionTrackerServiceServer).NewJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/versiontracker.VersionTrackerService/NewJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionTrackerServiceServer).NewJob(ctx, req.(*NewJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _VersionTrackerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "versiontracker.VersionTrackerService",
	HandlerType: (*VersionTrackerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewVersion",
			Handler:    _VersionTrackerService_NewVersion_Handler,
		},
		{
			MethodName: "NewJob",
			Handler:    _VersionTrackerService_NewJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "versiontracker.proto",
}
