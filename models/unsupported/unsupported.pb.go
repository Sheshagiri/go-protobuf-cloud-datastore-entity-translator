// Code generated by protoc-gen-go. DO NOT EDIT.
// source: unsupported.proto

package unsupported

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

type Model struct {
	Uint32Key            uint32   `protobuf:"varint,1,opt,name=uint32_key,json=uint32Key,proto3" json:"uint32_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Model) Reset()         { *m = Model{} }
func (m *Model) String() string { return proto.CompactTextString(m) }
func (*Model) ProtoMessage()    {}
func (*Model) Descriptor() ([]byte, []int) {
	return fileDescriptor_86386df630750539, []int{0}
}

func (m *Model) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Model.Unmarshal(m, b)
}
func (m *Model) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Model.Marshal(b, m, deterministic)
}
func (m *Model) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Model.Merge(m, src)
}
func (m *Model) XXX_Size() int {
	return xxx_messageInfo_Model.Size(m)
}
func (m *Model) XXX_DiscardUnknown() {
	xxx_messageInfo_Model.DiscardUnknown(m)
}

var xxx_messageInfo_Model proto.InternalMessageInfo

func (m *Model) GetUint32Key() uint32 {
	if m != nil {
		return m.Uint32Key
	}
	return 0
}

func init() {
	proto.RegisterType((*Model)(nil), "Model")
}

func init() { proto.RegisterFile("unsupported.proto", fileDescriptor_86386df630750539) }

var fileDescriptor_86386df630750539 = []byte{
	// 84 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0xcd, 0x2b, 0x2e,
	0x2d, 0x28, 0xc8, 0x2f, 0x2a, 0x49, 0x4d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe3,
	0x62, 0xf5, 0xcd, 0x4f, 0x49, 0xcd, 0x11, 0x92, 0xe5, 0xe2, 0x2a, 0xcd, 0xcc, 0x2b, 0x31, 0x36,
	0x8a, 0xcf, 0x4e, 0xad, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0xe2, 0x84, 0x88, 0x78, 0xa7,
	0x56, 0x26, 0xb1, 0x81, 0x95, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x9a, 0xf2, 0x8b, 0x02,
	0x43, 0x00, 0x00, 0x00,
}
