// Code generated by protoc-gen-go. DO NOT EDIT.
// source: confirm_message.proto

package messages

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ConfirmMessage struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Hash                 string   `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmMessage) Reset()         { *m = ConfirmMessage{} }
func (m *ConfirmMessage) String() string { return proto.CompactTextString(m) }
func (*ConfirmMessage) ProtoMessage()    {}
func (*ConfirmMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f5e4c6ca99b9eef, []int{0}
}

func (m *ConfirmMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmMessage.Unmarshal(m, b)
}
func (m *ConfirmMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmMessage.Marshal(b, m, deterministic)
}
func (m *ConfirmMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmMessage.Merge(m, src)
}
func (m *ConfirmMessage) XXX_Size() int {
	return xxx_messageInfo_ConfirmMessage.Size(m)
}
func (m *ConfirmMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmMessage proto.InternalMessageInfo

func (m *ConfirmMessage) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ConfirmMessage) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func init() {
	proto.RegisterType((*ConfirmMessage)(nil), "messages.ConfirmMessage")
}

func init() { proto.RegisterFile("confirm_message.proto", fileDescriptor_0f5e4c6ca99b9eef) }

var fileDescriptor_0f5e4c6ca99b9eef = []byte{
	// 96 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4d, 0xce, 0xcf, 0x4b,
	0xcb, 0x2c, 0xca, 0x8d, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0x80, 0x72, 0x8b, 0x95, 0xcc, 0xb8, 0xf8, 0x9c, 0x21, 0x4a, 0x7c, 0x21, 0x42,
	0x42, 0x02, 0x5c, 0xcc, 0xd9, 0xa9, 0x95, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x20, 0xa6,
	0x90, 0x10, 0x17, 0x4b, 0x46, 0x62, 0x71, 0x86, 0x04, 0x13, 0x58, 0x08, 0xcc, 0x4e, 0x62, 0x03,
	0x1b, 0x64, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x95, 0x88, 0x4c, 0xa0, 0x61, 0x00, 0x00, 0x00,
}
