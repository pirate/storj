// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: piecestore.proto

package pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/duration"
	math "math"
	grpc "storj.io/fork/google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type BandwidthAction int32

const (
	BandwidthAction_PUT        BandwidthAction = 0
	BandwidthAction_GET        BandwidthAction = 1
	BandwidthAction_GET_AUDIT  BandwidthAction = 2
	BandwidthAction_GET_REPAIR BandwidthAction = 3
	BandwidthAction_PUT_REPAIR BandwidthAction = 4
)

var BandwidthAction_name = map[int32]string{
	0: "PUT",
	1: "GET",
	2: "GET_AUDIT",
	3: "GET_REPAIR",
	4: "PUT_REPAIR",
}

var BandwidthAction_value = map[string]int32{
	"PUT":        0,
	"GET":        1,
	"GET_AUDIT":  2,
	"GET_REPAIR": 3,
	"PUT_REPAIR": 4,
}

func (x BandwidthAction) String() string {
	return proto.EnumName(BandwidthAction_name, int32(x))
}

func (BandwidthAction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_569d535d76469daf, []int{0}
}

type PayerBandwidthAllocation struct {
	SatelliteId          NodeID          `protobuf:"bytes,1,opt,name=satellite_id,json=satelliteId,proto3,customtype=NodeID" json:"satellite_id"`
	UplinkId             NodeID          `protobuf:"bytes,2,opt,name=uplink_id,json=uplinkId,proto3,customtype=NodeID" json:"uplink_id"`
	MaxSize              int64           `protobuf:"varint,3,opt,name=max_size,json=maxSize,proto3" json:"max_size,omitempty"`
	ExpirationUnixSec    int64           `protobuf:"varint,4,opt,name=expiration_unix_sec,json=expirationUnixSec,proto3" json:"expiration_unix_sec,omitempty"`
	SerialNumber         string          `protobuf:"bytes,5,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	Action               BandwidthAction `protobuf:"varint,6,opt,name=action,proto3,enum=piecestoreroutes.BandwidthAction" json:"action,omitempty"`
	CreatedUnixSec       int64           `protobuf:"varint,7,opt,name=created_unix_sec,json=createdUnixSec,proto3" json:"created_unix_sec,omitempty"`
	Certs                [][]byte        `protobuf:"bytes,8,rep,name=certs,proto3" json:"certs,omitempty"`
	Signature            []byte          `protobuf:"bytes,9,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PayerBandwidthAllocation) Reset()         { *m = PayerBandwidthAllocation{} }
func (m *PayerBandwidthAllocation) String() string { return proto.CompactTextString(m) }
func (*PayerBandwidthAllocation) ProtoMessage()    {}
func (*PayerBandwidthAllocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_569d535d76469daf, []int{0}
}
func (m *PayerBandwidthAllocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayerBandwidthAllocation.Unmarshal(m, b)
}
func (m *PayerBandwidthAllocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayerBandwidthAllocation.Marshal(b, m, deterministic)
}
func (m *PayerBandwidthAllocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayerBandwidthAllocation.Merge(m, src)
}
func (m *PayerBandwidthAllocation) XXX_Size() int {
	return xxx_messageInfo_PayerBandwidthAllocation.Size(m)
}
func (m *PayerBandwidthAllocation) XXX_DiscardUnknown() {
	xxx_messageInfo_PayerBandwidthAllocation.DiscardUnknown(m)
}

var xxx_messageInfo_PayerBandwidthAllocation proto.InternalMessageInfo

func (m *PayerBandwidthAllocation) GetMaxSize() int64 {
	if m != nil {
		return m.MaxSize
	}
	return 0
}

func (m *PayerBandwidthAllocation) GetExpirationUnixSec() int64 {
	if m != nil {
		return m.ExpirationUnixSec
	}
	return 0
}

func (m *PayerBandwidthAllocation) GetSerialNumber() string {
	if m != nil {
		return m.SerialNumber
	}
	return ""
}

func (m *PayerBandwidthAllocation) GetAction() BandwidthAction {
	if m != nil {
		return m.Action
	}
	return BandwidthAction_PUT
}

func (m *PayerBandwidthAllocation) GetCreatedUnixSec() int64 {
	if m != nil {
		return m.CreatedUnixSec
	}
	return 0
}

func (m *PayerBandwidthAllocation) GetCerts() [][]byte {
	if m != nil {
		return m.Certs
	}
	return nil
}

func (m *PayerBandwidthAllocation) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type RenterBandwidthAllocation struct {
	PayerAllocation      PayerBandwidthAllocation `protobuf:"bytes,1,opt,name=payer_allocation,json=payerAllocation,proto3" json:"payer_allocation"`
	Total                int64                    `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	StorageNodeId        NodeID                   `protobuf:"bytes,3,opt,name=storage_node_id,json=storageNodeId,proto3,customtype=NodeID" json:"storage_node_id"`
	Certs                [][]byte                 `protobuf:"bytes,4,rep,name=certs,proto3" json:"certs,omitempty"`
	Signature            []byte                   `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *RenterBandwidthAllocation) Reset()         { *m = RenterBandwidthAllocation{} }
func (m *RenterBandwidthAllocation) String() string { return proto.CompactTextString(m) }
func (*RenterBandwidthAllocation) ProtoMessage()    {}
func (*RenterBandwidthAllocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_569d535d76469daf, []int{1}
}
func (m *RenterBandwidthAllocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RenterBandwidthAllocation.Unmarshal(m, b)
}
func (m *RenterBandwidthAllocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RenterBandwidthAllocation.Marshal(b, m, deterministic)
}
func (m *RenterBandwidthAllocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RenterBandwidthAllocation.Merge(m, src)
}
func (m *RenterBandwidthAllocation) XXX_Size() int {
	return xxx_messageInfo_RenterBandwidthAllocation.Size(m)
}
func (m *RenterBandwidthAllocation) XXX_DiscardUnknown() {
	xxx_messageInfo_RenterBandwidthAllocation.DiscardUnknown(m)
}

var xxx_messageInfo_RenterBandwidthAllocation proto.InternalMessageInfo

func (m *RenterBandwidthAllocation) GetPayerAllocation() PayerBandwidthAllocation {
	if m != nil {
		return m.PayerAllocation
	}
	return PayerBandwidthAllocation{}
}

func (m *RenterBandwidthAllocation) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *RenterBandwidthAllocation) GetCerts() [][]byte {
	if m != nil {
		return m.Certs
	}
	return nil
}

func (m *RenterBandwidthAllocation) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type SignedMessage struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	PublicKey            []byte   `protobuf:"bytes,3,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedMessage) Reset()         { *m = SignedMessage{} }
func (m *SignedMessage) String() string { return proto.CompactTextString(m) }
func (*SignedMessage) ProtoMessage()    {}
func (*SignedMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_569d535d76469daf, []int{2}
}
func (m *SignedMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedMessage.Unmarshal(m, b)
}
func (m *SignedMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedMessage.Marshal(b, m, deterministic)
}
func (m *SignedMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedMessage.Merge(m, src)
}
func (m *SignedMessage) XXX_Size() int {
	return xxx_messageInfo_SignedMessage.Size(m)
}
func (m *SignedMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SignedMessage proto.InternalMessageInfo

func (m *SignedMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SignedMessage) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *SignedMessage) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

type SignedHash struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Certs                [][]byte `protobuf:"bytes,2,rep,name=certs,proto3" json:"certs,omitempty"`
	Signature            []byte   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedHash) Reset()         { *m = SignedHash{} }
func (m *SignedHash) String() string { return proto.CompactTextString(m) }
func (*SignedHash) ProtoMessage()    {}
func (*SignedHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_569d535d76469daf, []int{3}
}
func (m *SignedHash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedHash.Unmarshal(m, b)
}
func (m *SignedHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedHash.Marshal(b, m, deterministic)
}
func (m *SignedHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedHash.Merge(m, src)
}
func (m *SignedHash) XXX_Size() int {
	return xxx_messageInfo_SignedHash.Size(m)
}
func (m *SignedHash) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedHash.DiscardUnknown(m)
}

var xxx_messageInfo_SignedHash proto.InternalMessageInfo

func (m *SignedHash) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *SignedHash) GetCerts() [][]byte {
	if m != nil {
		return m.Certs
	}
	return nil
}

func (m *SignedHash) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterEnum("piecestoreroutes.BandwidthAction", BandwidthAction_name, BandwidthAction_value)
	proto.RegisterType((*PayerBandwidthAllocation)(nil), "piecestoreroutes.PayerBandwidthAllocation")
	proto.RegisterType((*RenterBandwidthAllocation)(nil), "piecestoreroutes.RenterBandwidthAllocation")
	proto.RegisterType((*SignedMessage)(nil), "piecestoreroutes.SignedMessage")
	proto.RegisterType((*SignedHash)(nil), "piecestoreroutes.SignedHash")
}

func init() { proto.RegisterFile("piecestore.proto", fileDescriptor_569d535d76469daf) }

var fileDescriptor_569d535d76469daf = []byte{
	// 548 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x41, 0x4f, 0xdb, 0x3e,
	0x1c, 0x25, 0x4d, 0x5a, 0xc8, 0x8f, 0xb6, 0xe4, 0xef, 0x3f, 0x87, 0x80, 0xb6, 0xd1, 0xb1, 0x4b,
	0xc4, 0xa4, 0xa2, 0x31, 0x69, 0xd2, 0x8e, 0x54, 0x20, 0x56, 0x4d, 0x43, 0x95, 0x69, 0x2f, 0xdb,
	0x21, 0x73, 0xe3, 0xdf, 0x52, 0x8b, 0x10, 0x47, 0xb1, 0xa3, 0x15, 0xae, 0xfb, 0x72, 0xfb, 0x0c,
	0x3b, 0xf0, 0x41, 0x76, 0x9a, 0xe2, 0x40, 0xc3, 0xba, 0xf6, 0xe6, 0xf7, 0xfc, 0xf2, 0xfc, 0x7e,
	0x7e, 0x0e, 0x78, 0x99, 0xc0, 0x08, 0x95, 0x96, 0x39, 0xf6, 0xb3, 0x5c, 0x6a, 0x49, 0x9e, 0x30,
	0xb9, 0x2c, 0x34, 0xaa, 0x7d, 0x88, 0x65, 0x2c, 0xab, 0xdd, 0xfd, 0x17, 0xb1, 0x94, 0x71, 0x82,
	0xc7, 0x06, 0x4d, 0x8b, 0x6f, 0xc7, 0xbc, 0xc8, 0x99, 0x16, 0x32, 0xad, 0xf6, 0x0f, 0x7f, 0xd8,
	0xe0, 0x8f, 0xd8, 0x2d, 0xe6, 0x03, 0x96, 0xf2, 0xef, 0x82, 0xeb, 0xd9, 0x69, 0x92, 0xc8, 0xc8,
	0x48, 0xc8, 0x1b, 0x68, 0x2b, 0xa6, 0x31, 0x49, 0x84, 0xc6, 0x50, 0x70, 0xdf, 0xea, 0x59, 0x41,
	0x7b, 0xd0, 0xfd, 0x79, 0x7f, 0xb0, 0xf1, 0xeb, 0xfe, 0xa0, 0x75, 0x29, 0x39, 0x0e, 0xcf, 0xe8,
	0xf6, 0x42, 0x33, 0xe4, 0xe4, 0x35, 0xb8, 0x45, 0x96, 0x88, 0xf4, 0xba, 0xd4, 0x37, 0x56, 0xea,
	0xb7, 0x2a, 0xc1, 0x90, 0x93, 0x3d, 0xd8, 0xba, 0x61, 0xf3, 0x50, 0x89, 0x3b, 0xf4, 0xed, 0x9e,
	0x15, 0xd8, 0x74, 0xf3, 0x86, 0xcd, 0xaf, 0xc4, 0x1d, 0x92, 0x3e, 0xfc, 0x8f, 0xf3, 0x4c, 0x54,
	0x59, 0xc3, 0x22, 0x15, 0xf3, 0x50, 0x61, 0xe4, 0x3b, 0x46, 0xf5, 0x5f, 0xbd, 0x35, 0x49, 0xc5,
	0xfc, 0x0a, 0x23, 0xf2, 0x0a, 0x3a, 0x0a, 0x73, 0xc1, 0x92, 0x30, 0x2d, 0x6e, 0xa6, 0x98, 0xfb,
	0xcd, 0x9e, 0x15, 0xb8, 0xb4, 0x5d, 0x91, 0x97, 0x86, 0x23, 0xef, 0xa1, 0xc5, 0xa2, 0xf2, 0x2b,
	0xbf, 0xd5, 0xb3, 0x82, 0xee, 0xc9, 0xcb, 0xfe, 0xf2, 0xdd, 0xf5, 0xeb, 0x6b, 0x30, 0x42, 0xfa,
	0xf0, 0x01, 0x09, 0xc0, 0x8b, 0x72, 0x64, 0x1a, 0x79, 0x1d, 0x66, 0xd3, 0x84, 0xe9, 0x3e, 0xf0,
	0x8f, 0x49, 0x76, 0xa1, 0x19, 0x61, 0xae, 0x95, 0xbf, 0xd5, 0xb3, 0x83, 0x36, 0xad, 0x00, 0x79,
	0x06, 0xae, 0x12, 0x71, 0xca, 0x74, 0x91, 0xa3, 0xef, 0x96, 0xf7, 0x42, 0x6b, 0xe2, 0xf0, 0xb7,
	0x05, 0x7b, 0x14, 0x53, 0xbd, 0xba, 0x86, 0x2f, 0xe0, 0x65, 0x65, 0x45, 0x21, 0x5b, 0x70, 0xa6,
	0x8a, 0xed, 0x93, 0xa3, 0x7f, 0x07, 0x58, 0x57, 0xe6, 0xc0, 0x29, 0x6b, 0xa0, 0x3b, 0xc6, 0xe9,
	0x89, 0xf9, 0x2e, 0x34, 0xb5, 0xd4, 0x2c, 0x31, 0x65, 0xd9, 0xb4, 0x02, 0xe4, 0x1d, 0xec, 0x94,
	0xa6, 0x2c, 0xc6, 0x30, 0x95, 0xdc, 0x94, 0x6f, 0xaf, 0x2c, 0xb3, 0xf3, 0x20, 0x33, 0x90, 0xd7,
	0xc3, 0x3b, 0x6b, 0x87, 0x6f, 0x2e, 0x0f, 0xff, 0x15, 0x3a, 0x57, 0x22, 0x4e, 0x91, 0x7f, 0x42,
	0xa5, 0x58, 0x8c, 0x84, 0x80, 0xc3, 0x99, 0x66, 0xd5, 0x73, 0xa3, 0x66, 0xfd, 0xb7, 0x45, 0x63,
	0xc9, 0x82, 0x3c, 0x07, 0xc8, 0x8a, 0x69, 0x22, 0xa2, 0xf0, 0x1a, 0x6f, 0xab, 0xa4, 0xd4, 0xad,
	0x98, 0x8f, 0x78, 0x7b, 0x38, 0x06, 0xa8, 0x4e, 0xf8, 0xc0, 0xd4, 0xac, 0xb4, 0x9f, 0x31, 0x35,
	0x7b, 0xb4, 0x2f, 0xd7, 0x75, 0xee, 0xc6, 0xda, 0xdc, 0xf6, 0xd2, 0xa1, 0x47, 0x14, 0x76, 0x96,
	0x5e, 0x0b, 0xd9, 0x04, 0x7b, 0x34, 0x19, 0x7b, 0x1b, 0xe5, 0xe2, 0xe2, 0x7c, 0xec, 0x59, 0xa4,
	0x03, 0xee, 0xc5, 0xf9, 0x38, 0x3c, 0x9d, 0x9c, 0x0d, 0xc7, 0x5e, 0x83, 0x74, 0x01, 0x4a, 0x48,
	0xcf, 0x47, 0xa7, 0x43, 0xea, 0xd9, 0x25, 0x1e, 0x4d, 0x16, 0xd8, 0x19, 0x38, 0x9f, 0x1b, 0xd9,
	0x74, 0xda, 0x32, 0xff, 0xe6, 0xdb, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6a, 0x31, 0xcd, 0xe9,
	0xed, 0x03, 0x00, 0x00,
}
