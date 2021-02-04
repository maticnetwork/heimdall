// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: clerk/v1beta/msg.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgEventRecordRequest struct {
	From            string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	TxHash          string `protobuf:"bytes,2,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty" yaml:"tx_hash"`
	LogIndex        uint64 `protobuf:"varint,3,opt,name=log_index,json=logIndex,proto3" json:"log_index,omitempty" yaml:"log_index"`
	BlockNumber     uint64 `protobuf:"varint,4,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty" yaml:"block_number"`
	ContractAddress string `protobuf:"bytes,5,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty" yaml:"contract_address"`
	Data            []byte `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	Id              uint64 `protobuf:"varint,7,opt,name=id,proto3" json:"id,omitempty"`
	ChainId         string `protobuf:"bytes,8,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty" yaml:"chain_id"`
}

func (m *MsgEventRecordRequest) Reset()         { *m = MsgEventRecordRequest{} }
func (m *MsgEventRecordRequest) String() string { return proto.CompactTextString(m) }
func (*MsgEventRecordRequest) ProtoMessage()    {}
func (*MsgEventRecordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d767b6ebeb4c1f4, []int{0}
}
func (m *MsgEventRecordRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgEventRecordRequest.Unmarshal(m, b)
}
func (m *MsgEventRecordRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgEventRecordRequest.Marshal(b, m, deterministic)
}
func (m *MsgEventRecordRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEventRecordRequest.Merge(m, src)
}
func (m *MsgEventRecordRequest) XXX_Size() int {
	return xxx_messageInfo_MsgEventRecordRequest.Size(m)
}
func (m *MsgEventRecordRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEventRecordRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEventRecordRequest proto.InternalMessageInfo

// MsgEventRecordResponse defines MsgEventRecord response type.
type MsgEventRecordResponse struct {
}

func (m *MsgEventRecordResponse) Reset()         { *m = MsgEventRecordResponse{} }
func (m *MsgEventRecordResponse) String() string { return proto.CompactTextString(m) }
func (*MsgEventRecordResponse) ProtoMessage()    {}
func (*MsgEventRecordResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d767b6ebeb4c1f4, []int{1}
}
func (m *MsgEventRecordResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgEventRecordResponse.Unmarshal(m, b)
}
func (m *MsgEventRecordResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgEventRecordResponse.Marshal(b, m, deterministic)
}
func (m *MsgEventRecordResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEventRecordResponse.Merge(m, src)
}
func (m *MsgEventRecordResponse) XXX_Size() int {
	return xxx_messageInfo_MsgEventRecordResponse.Size(m)
}
func (m *MsgEventRecordResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEventRecordResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEventRecordResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgEventRecordRequest)(nil), "heimdall.clerk.v1beta1.MsgEventRecordRequest")
	proto.RegisterType((*MsgEventRecordResponse)(nil), "heimdall.clerk.v1beta1.MsgEventRecordResponse")
}

func init() { proto.RegisterFile("clerk/v1beta/msg.proto", fileDescriptor_2d767b6ebeb4c1f4) }

var fileDescriptor_2d767b6ebeb4c1f4 = []byte{
	// 436 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xc7, 0x9d, 0x34, 0x24, 0xe9, 0x51, 0xa5, 0xd5, 0xb5, 0xa4, 0xa7, 0x22, 0xd9, 0x91, 0xa7,
	0x48, 0x88, 0xb3, 0x02, 0x5b, 0x36, 0x22, 0x81, 0xda, 0xa1, 0x0c, 0x37, 0xb2, 0x58, 0x67, 0xdf,
	0x61, 0x5b, 0xb1, 0x7d, 0xe1, 0xee, 0x12, 0xdc, 0x6f, 0xc0, 0xc8, 0x47, 0xe0, 0xe3, 0x30, 0x76,
	0x64, 0xb2, 0x50, 0x32, 0xb1, 0xfa, 0x13, 0x20, 0xdf, 0xa5, 0x08, 0xaa, 0x0c, 0x6c, 0xff, 0x7b,
	0xff, 0xdf, 0xd3, 0xff, 0xd9, 0xef, 0x81, 0x71, 0x9c, 0x73, 0xb9, 0x0c, 0x36, 0xb3, 0x88, 0x6b,
	0x1a, 0x14, 0x2a, 0xc1, 0x2b, 0x29, 0xb4, 0x80, 0xe3, 0x94, 0x67, 0x05, 0xa3, 0x79, 0x8e, 0x0d,
	0x80, 0x2d, 0x30, 0xbb, 0xba, 0x48, 0x44, 0x22, 0x0c, 0x12, 0xb4, 0xca, 0xd2, 0xfe, 0xaf, 0x2e,
	0x78, 0x76, 0xab, 0x92, 0xb7, 0x1b, 0x5e, 0x6a, 0xc2, 0x63, 0x21, 0x19, 0xe1, 0x9f, 0xd6, 0x5c,
	0x69, 0x08, 0x41, 0xef, 0xa3, 0x14, 0x05, 0xea, 0x4c, 0x3a, 0xd3, 0x63, 0x62, 0x34, 0x7c, 0x01,
	0x06, 0xba, 0x0a, 0x53, 0xaa, 0x52, 0xd4, 0x6d, 0xcb, 0x0b, 0xd8, 0xd4, 0xde, 0xe8, 0x8e, 0x16,
	0xf9, 0xdc, 0xdf, 0x1b, 0x3e, 0xe9, 0xeb, 0xea, 0x9a, 0xaa, 0x14, 0xce, 0xc0, 0x71, 0x2e, 0x92,
	0x30, 0x2b, 0x19, 0xaf, 0xd0, 0xd1, 0xa4, 0x33, 0xed, 0x2d, 0x2e, 0x9a, 0xda, 0x3b, 0xb3, 0xf8,
	0x1f, 0xcb, 0x27, 0xc3, 0x5c, 0x24, 0x37, 0xad, 0x84, 0x73, 0x70, 0x12, 0xe5, 0x22, 0x5e, 0x86,
	0xe5, 0xba, 0x88, 0xb8, 0x44, 0x3d, 0xd3, 0x75, 0xd9, 0xd4, 0xde, 0xb9, 0xed, 0xfa, 0xdb, 0xf5,
	0xc9, 0x53, 0xf3, 0x7c, 0x6f, 0x5e, 0xf0, 0x1d, 0x38, 0x8b, 0x45, 0xa9, 0x25, 0x8d, 0x75, 0x48,
	0x19, 0x93, 0x5c, 0x29, 0xf4, 0xc4, 0x0c, 0xf9, 0xbc, 0xa9, 0xbd, 0x4b, 0xdb, 0xff, 0x98, 0xf0,
	0xc9, 0xe9, 0x43, 0xe9, 0x8d, 0xad, 0xb4, 0xdf, 0xcd, 0xa8, 0xa6, 0xa8, 0x3f, 0xe9, 0x4c, 0x4f,
	0x88, 0xd1, 0x70, 0x04, 0xba, 0x19, 0x43, 0x83, 0x76, 0x1a, 0xd2, 0xcd, 0x18, 0xc4, 0x60, 0x18,
	0xa7, 0x34, 0x2b, 0xc3, 0x8c, 0xa1, 0xa1, 0xc9, 0x38, 0x6f, 0x6a, 0xef, 0x74, 0x9f, 0xb1, 0x77,
	0x7c, 0x32, 0x30, 0xf2, 0x86, 0xcd, 0x7b, 0x5f, 0xbe, 0x79, 0x8e, 0x8f, 0xc0, 0xf8, 0xf1, 0xaf,
	0x56, 0x2b, 0x51, 0x2a, 0xfe, 0x6a, 0x03, 0x8e, 0x6e, 0x55, 0x02, 0x05, 0x18, 0xfd, 0x0b, 0xc0,
	0x97, 0xf8, 0xf0, 0x36, 0xf1, 0xc1, 0x9d, 0x5d, 0xe1, 0xff, 0xc5, 0x6d, 0xee, 0xe2, 0xfa, 0xfb,
	0xd6, 0x75, 0xee, 0xb7, 0xae, 0xf3, 0x73, 0xeb, 0x3a, 0x5f, 0x77, 0xae, 0x73, 0xbf, 0x73, 0x9d,
	0x1f, 0x3b, 0xd7, 0xf9, 0x80, 0x93, 0x4c, 0xa7, 0xeb, 0x08, 0xc7, 0xa2, 0x08, 0x0a, 0xaa, 0xb3,
	0xb8, 0xe4, 0xfa, 0xb3, 0x90, 0xcb, 0xe0, 0x21, 0x20, 0xa8, 0x02, 0x7b, 0x80, 0xfa, 0x6e, 0xc5,
	0x55, 0xd4, 0x37, 0xe7, 0xf4, 0xfa, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5a, 0x00, 0xbb, 0x07,
	0x96, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// MsgEventRecord defines a method to join a new event record.
	MsgEventRecord(ctx context.Context, in *MsgEventRecordRequest, opts ...grpc.CallOption) (*MsgEventRecordResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) MsgEventRecord(ctx context.Context, in *MsgEventRecordRequest, opts ...grpc.CallOption) (*MsgEventRecordResponse, error) {
	out := new(MsgEventRecordResponse)
	err := c.cc.Invoke(ctx, "/heimdall.clerk.v1beta1.Msg/MsgEventRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// MsgEventRecord defines a method to join a new event record.
	MsgEventRecord(context.Context, *MsgEventRecordRequest) (*MsgEventRecordResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) MsgEventRecord(ctx context.Context, req *MsgEventRecordRequest) (*MsgEventRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MsgEventRecord not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_MsgEventRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEventRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).MsgEventRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/heimdall.clerk.v1beta1.Msg/MsgEventRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).MsgEventRecord(ctx, req.(*MsgEventRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "heimdall.clerk.v1beta1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MsgEventRecord",
			Handler:    _Msg_MsgEventRecord_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "clerk/v1beta/msg.proto",
}