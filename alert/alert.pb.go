// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.2
// source: alert/alert.proto

package alert

import (
	context "context"
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

type RequestAlert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host    string `protobuf:"bytes,1,opt,name=Host,proto3" json:"Host,omitempty"`
	Service string `protobuf:"bytes,2,opt,name=Service,proto3" json:"Service,omitempty"`
}

func (x *RequestAlert) Reset() {
	*x = RequestAlert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alert_alert_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAlert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAlert) ProtoMessage() {}

func (x *RequestAlert) ProtoReflect() protoreflect.Message {
	mi := &file_alert_alert_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestAlert.ProtoReflect.Descriptor instead.
func (*RequestAlert) Descriptor() ([]byte, []int) {
	return file_alert_alert_proto_rawDescGZIP(), []int{0}
}

func (x *RequestAlert) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *RequestAlert) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

type ResponseAlert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host     string `protobuf:"bytes,1,opt,name=Host,proto3" json:"Host,omitempty"`
	Service  string `protobuf:"bytes,2,opt,name=Service,proto3" json:"Service,omitempty"`
	NotesUrl string `protobuf:"bytes,3,opt,name=NotesUrl,proto3" json:"NotesUrl,omitempty"`
}

func (x *ResponseAlert) Reset() {
	*x = ResponseAlert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alert_alert_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseAlert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseAlert) ProtoMessage() {}

func (x *ResponseAlert) ProtoReflect() protoreflect.Message {
	mi := &file_alert_alert_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseAlert.ProtoReflect.Descriptor instead.
func (*ResponseAlert) Descriptor() ([]byte, []int) {
	return file_alert_alert_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseAlert) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *ResponseAlert) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *ResponseAlert) GetNotesUrl() string {
	if x != nil {
		return x.NotesUrl
	}
	return ""
}

type ResponseAlertsList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List string `protobuf:"bytes,1,opt,name=List,proto3" json:"List,omitempty"`
}

func (x *ResponseAlertsList) Reset() {
	*x = ResponseAlertsList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alert_alert_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseAlertsList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseAlertsList) ProtoMessage() {}

func (x *ResponseAlertsList) ProtoReflect() protoreflect.Message {
	mi := &file_alert_alert_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseAlertsList.ProtoReflect.Descriptor instead.
func (*ResponseAlertsList) Descriptor() ([]byte, []int) {
	return file_alert_alert_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseAlertsList) GetList() string {
	if x != nil {
		return x.List
	}
	return ""
}

var File_alert_alert_proto protoreflect.FileDescriptor

var file_alert_alert_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x2f, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x22, 0x3c, 0x0a, 0x0c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x6f,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x59, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x6f, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x65, 0x73,
	0x55, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4e, 0x6f, 0x74, 0x65, 0x73,
	0x55, 0x72, 0x6c, 0x22, 0x28, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41,
	0x6c, 0x65, 0x72, 0x74, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x4c, 0x69, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x32, 0x86, 0x01,
	0x0a, 0x05, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x3a, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4e, 0x6f,
	0x74, 0x65, 0x73, 0x55, 0x52, 0x4c, 0x12, 0x13, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x1a, 0x14, 0x2e, 0x61, 0x6c,
	0x65, 0x72, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x6c, 0x65, 0x72,
	0x74, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x13, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x6c, 0x65, 0x72,
	0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x61, 0x6c, 0x65, 0x72, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_alert_alert_proto_rawDescOnce sync.Once
	file_alert_alert_proto_rawDescData = file_alert_alert_proto_rawDesc
)

func file_alert_alert_proto_rawDescGZIP() []byte {
	file_alert_alert_proto_rawDescOnce.Do(func() {
		file_alert_alert_proto_rawDescData = protoimpl.X.CompressGZIP(file_alert_alert_proto_rawDescData)
	})
	return file_alert_alert_proto_rawDescData
}

var file_alert_alert_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_alert_alert_proto_goTypes = []interface{}{
	(*RequestAlert)(nil),       // 0: alert.RequestAlert
	(*ResponseAlert)(nil),      // 1: alert.ResponseAlert
	(*ResponseAlertsList)(nil), // 2: alert.ResponseAlertsList
}
var file_alert_alert_proto_depIdxs = []int32{
	0, // 0: alert.Alert.GetNotesURL:input_type -> alert.RequestAlert
	0, // 1: alert.Alert.GetAlertsList:input_type -> alert.RequestAlert
	1, // 2: alert.Alert.GetNotesURL:output_type -> alert.ResponseAlert
	2, // 3: alert.Alert.GetAlertsList:output_type -> alert.ResponseAlertsList
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_alert_alert_proto_init() }
func file_alert_alert_proto_init() {
	if File_alert_alert_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_alert_alert_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestAlert); i {
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
		file_alert_alert_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseAlert); i {
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
		file_alert_alert_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseAlertsList); i {
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
			RawDescriptor: file_alert_alert_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_alert_alert_proto_goTypes,
		DependencyIndexes: file_alert_alert_proto_depIdxs,
		MessageInfos:      file_alert_alert_proto_msgTypes,
	}.Build()
	File_alert_alert_proto = out.File
	file_alert_alert_proto_rawDesc = nil
	file_alert_alert_proto_goTypes = nil
	file_alert_alert_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AlertClient is the client API for Alert service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AlertClient interface {
	GetNotesURL(ctx context.Context, in *RequestAlert, opts ...grpc.CallOption) (*ResponseAlert, error)
	GetAlertsList(ctx context.Context, in *RequestAlert, opts ...grpc.CallOption) (*ResponseAlertsList, error)
}

type alertClient struct {
	cc grpc.ClientConnInterface
}

func NewAlertClient(cc grpc.ClientConnInterface) AlertClient {
	return &alertClient{cc}
}

func (c *alertClient) GetNotesURL(ctx context.Context, in *RequestAlert, opts ...grpc.CallOption) (*ResponseAlert, error) {
	out := new(ResponseAlert)
	err := c.cc.Invoke(ctx, "/alert.Alert/GetNotesURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertClient) GetAlertsList(ctx context.Context, in *RequestAlert, opts ...grpc.CallOption) (*ResponseAlertsList, error) {
	out := new(ResponseAlertsList)
	err := c.cc.Invoke(ctx, "/alert.Alert/GetAlertsList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlertServer is the server API for Alert service.
type AlertServer interface {
	GetNotesURL(context.Context, *RequestAlert) (*ResponseAlert, error)
	GetAlertsList(context.Context, *RequestAlert) (*ResponseAlertsList, error)
}

// UnimplementedAlertServer can be embedded to have forward compatible implementations.
type UnimplementedAlertServer struct {
}

func (*UnimplementedAlertServer) GetNotesURL(context.Context, *RequestAlert) (*ResponseAlert, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotesURL not implemented")
}
func (*UnimplementedAlertServer) GetAlertsList(context.Context, *RequestAlert) (*ResponseAlertsList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlertsList not implemented")
}

func RegisterAlertServer(s *grpc.Server, srv AlertServer) {
	s.RegisterService(&_Alert_serviceDesc, srv)
}

func _Alert_GetNotesURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAlert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertServer).GetNotesURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alert.Alert/GetNotesURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertServer).GetNotesURL(ctx, req.(*RequestAlert))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alert_GetAlertsList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAlert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertServer).GetAlertsList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/alert.Alert/GetAlertsList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertServer).GetAlertsList(ctx, req.(*RequestAlert))
	}
	return interceptor(ctx, in, info, handler)
}

var _Alert_serviceDesc = grpc.ServiceDesc{
	ServiceName: "alert.Alert",
	HandlerType: (*AlertServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNotesURL",
			Handler:    _Alert_GetNotesURL_Handler,
		},
		{
			MethodName: "GetAlertsList",
			Handler:    _Alert_GetAlertsList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alert/alert.proto",
}
