// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/manager/manager.proto

package manager

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

type Protocol int32

const (
	Protocol_WSS  Protocol = 0
	Protocol_GRPC Protocol = 1
)

// Enum value maps for Protocol.
var (
	Protocol_name = map[int32]string{
		0: "WSS",
		1: "GRPC",
	}
	Protocol_value = map[string]int32{
		"WSS":  0,
		"GRPC": 1,
	}
)

func (x Protocol) Enum() *Protocol {
	p := new(Protocol)
	*p = x
	return p
}

func (x Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_manager_manager_proto_enumTypes[0].Descriptor()
}

func (Protocol) Type() protoreflect.EnumType {
	return &file_proto_manager_manager_proto_enumTypes[0]
}

func (x Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Protocol.Descriptor instead.
func (Protocol) EnumDescriptor() ([]byte, []int) {
	return file_proto_manager_manager_proto_rawDescGZIP(), []int{0}
}

type NodeRegistrationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip              string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	InternetAddress string   `protobuf:"bytes,2,opt,name=internetAddress,proto3" json:"internetAddress,omitempty"`
	NodeId          string   `protobuf:"bytes,3,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Protocol        Protocol `protobuf:"varint,4,opt,name=protocol,proto3,enum=proto.Protocol" json:"protocol,omitempty"`
	WssPath         string   `protobuf:"bytes,5,opt,name=wssPath,proto3" json:"wssPath,omitempty"`
}

func (x *NodeRegistrationRequest) Reset() {
	*x = NodeRegistrationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_manager_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeRegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeRegistrationRequest) ProtoMessage() {}

func (x *NodeRegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_manager_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeRegistrationRequest.ProtoReflect.Descriptor instead.
func (*NodeRegistrationRequest) Descriptor() ([]byte, []int) {
	return file_proto_manager_manager_proto_rawDescGZIP(), []int{0}
}

func (x *NodeRegistrationRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *NodeRegistrationRequest) GetInternetAddress() string {
	if x != nil {
		return x.InternetAddress
	}
	return ""
}

func (x *NodeRegistrationRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *NodeRegistrationRequest) GetProtocol() Protocol {
	if x != nil {
		return x.Protocol
	}
	return Protocol_WSS
}

func (x *NodeRegistrationRequest) GetWssPath() string {
	if x != nil {
		return x.WssPath
	}
	return ""
}

type NodeRegistrationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId string `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	AesKey string `protobuf:"bytes,2,opt,name=aesKey,proto3" json:"aesKey,omitempty"` // 用于解析用户 token
}

func (x *NodeRegistrationResponse) Reset() {
	*x = NodeRegistrationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_manager_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeRegistrationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeRegistrationResponse) ProtoMessage() {}

func (x *NodeRegistrationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_manager_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeRegistrationResponse.ProtoReflect.Descriptor instead.
func (*NodeRegistrationResponse) Descriptor() ([]byte, []int) {
	return file_proto_manager_manager_proto_rawDescGZIP(), []int{1}
}

func (x *NodeRegistrationResponse) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *NodeRegistrationResponse) GetAesKey() string {
	if x != nil {
		return x.AesKey
	}
	return ""
}

type TrafficReportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jwt     string `protobuf:"bytes,1,opt,name=jwt,proto3" json:"jwt,omitempty"`          // 用户jwt
	Traffic int64  `protobuf:"varint,2,opt,name=traffic,proto3" json:"traffic,omitempty"` // 用户每消耗10m 上报一次
}

func (x *TrafficReportRequest) Reset() {
	*x = TrafficReportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_manager_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrafficReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrafficReportRequest) ProtoMessage() {}

func (x *TrafficReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_manager_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrafficReportRequest.ProtoReflect.Descriptor instead.
func (*TrafficReportRequest) Descriptor() ([]byte, []int) {
	return file_proto_manager_manager_proto_rawDescGZIP(), []int{2}
}

func (x *TrafficReportRequest) GetJwt() string {
	if x != nil {
		return x.Jwt
	}
	return ""
}

func (x *TrafficReportRequest) GetTraffic() int64 {
	if x != nil {
		return x.Traffic
	}
	return 0
}

type TrafficReportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TrafficReportResponse) Reset() {
	*x = TrafficReportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_manager_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrafficReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrafficReportResponse) ProtoMessage() {}

func (x *TrafficReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_manager_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrafficReportResponse.ProtoReflect.Descriptor instead.
func (*TrafficReportResponse) Descriptor() ([]byte, []int) {
	return file_proto_manager_manager_proto_rawDescGZIP(), []int{3}
}

var File_proto_manager_manager_proto protoreflect.FileDescriptor

var file_proto_manager_manager_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb2, 0x01, 0x0a, 0x17, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70,
	0x12, 0x28, 0x0a, 0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f,
	0x64, 0x65, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x12, 0x2b, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12,
	0x18, 0x0a, 0x07, 0x77, 0x73, 0x73, 0x50, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x77, 0x73, 0x73, 0x50, 0x61, 0x74, 0x68, 0x22, 0x4a, 0x0a, 0x18, 0x4e, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x65, 0x73, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x65, 0x73, 0x4b, 0x65, 0x79, 0x22, 0x42, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x6a, 0x77, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6a, 0x77, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x22, 0x17, 0x0a, 0x15, 0x54, 0x72, 0x61,
	0x66, 0x66, 0x69, 0x63, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2a, 0x1d, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x07,
	0x0a, 0x03, 0x57, 0x53, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x52, 0x50, 0x43, 0x10,
	0x01, 0x32, 0xb7, 0x01, 0x0a, 0x10, 0x47, 0x75, 0x61, 0x72, 0x64, 0x4c, 0x69, 0x6e, 0x6b, 0x4d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x55, 0x0a, 0x10, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a,
	0x0d, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x47, 0x75, 0x61, 0x72, 0x64, 0x4c, 0x69, 0x6e, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_manager_manager_proto_rawDescOnce sync.Once
	file_proto_manager_manager_proto_rawDescData = file_proto_manager_manager_proto_rawDesc
)

func file_proto_manager_manager_proto_rawDescGZIP() []byte {
	file_proto_manager_manager_proto_rawDescOnce.Do(func() {
		file_proto_manager_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_manager_manager_proto_rawDescData)
	})
	return file_proto_manager_manager_proto_rawDescData
}

var file_proto_manager_manager_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_manager_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_manager_manager_proto_goTypes = []interface{}{
	(Protocol)(0),                    // 0: proto.protocol
	(*NodeRegistrationRequest)(nil),  // 1: proto.NodeRegistrationRequest
	(*NodeRegistrationResponse)(nil), // 2: proto.NodeRegistrationResponse
	(*TrafficReportRequest)(nil),     // 3: proto.TrafficReportRequest
	(*TrafficReportResponse)(nil),    // 4: proto.TrafficReportResponse
}
var file_proto_manager_manager_proto_depIdxs = []int32{
	0, // 0: proto.NodeRegistrationRequest.protocol:type_name -> proto.protocol
	1, // 1: proto.GuardLinkManager.NodeRegistration:input_type -> proto.NodeRegistrationRequest
	3, // 2: proto.GuardLinkManager.TrafficReport:input_type -> proto.TrafficReportRequest
	2, // 3: proto.GuardLinkManager.NodeRegistration:output_type -> proto.NodeRegistrationResponse
	4, // 4: proto.GuardLinkManager.TrafficReport:output_type -> proto.TrafficReportResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_manager_manager_proto_init() }
func file_proto_manager_manager_proto_init() {
	if File_proto_manager_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_manager_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeRegistrationRequest); i {
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
		file_proto_manager_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeRegistrationResponse); i {
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
		file_proto_manager_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrafficReportRequest); i {
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
		file_proto_manager_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrafficReportResponse); i {
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
			RawDescriptor: file_proto_manager_manager_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_manager_manager_proto_goTypes,
		DependencyIndexes: file_proto_manager_manager_proto_depIdxs,
		EnumInfos:         file_proto_manager_manager_proto_enumTypes,
		MessageInfos:      file_proto_manager_manager_proto_msgTypes,
	}.Build()
	File_proto_manager_manager_proto = out.File
	file_proto_manager_manager_proto_rawDesc = nil
	file_proto_manager_manager_proto_goTypes = nil
	file_proto_manager_manager_proto_depIdxs = nil
}