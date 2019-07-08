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

type Parent struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Parent) Reset()         { *m = Parent{} }
func (m *Parent) String() string { return proto.CompactTextString(m) }
func (*Parent) ProtoMessage()    {}
func (*Parent) Descriptor() ([]byte, []int) {
	return fileDescriptor_86386df630750539, []int{1}
}

func (m *Parent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Parent.Unmarshal(m, b)
}
func (m *Parent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Parent.Marshal(b, m, deterministic)
}
func (m *Parent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Parent.Merge(m, src)
}
func (m *Parent) XXX_Size() int {
	return xxx_messageInfo_Parent.Size(m)
}
func (m *Parent) XXX_DiscardUnknown() {
	xxx_messageInfo_Parent.DiscardUnknown(m)
}

var xxx_messageInfo_Parent proto.InternalMessageInfo

func (m *Parent) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Child struct {
	Parent               *Parent  `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Child) Reset()         { *m = Child{} }
func (m *Child) String() string { return proto.CompactTextString(m) }
func (*Child) ProtoMessage()    {}
func (*Child) Descriptor() ([]byte, []int) {
	return fileDescriptor_86386df630750539, []int{2}
}

func (m *Child) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Child.Unmarshal(m, b)
}
func (m *Child) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Child.Marshal(b, m, deterministic)
}
func (m *Child) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Child.Merge(m, src)
}
func (m *Child) XXX_Size() int {
	return xxx_messageInfo_Child.Size(m)
}
func (m *Child) XXX_DiscardUnknown() {
	xxx_messageInfo_Child.DiscardUnknown(m)
}

var xxx_messageInfo_Child proto.InternalMessageInfo

func (m *Child) GetParent() *Parent {
	if m != nil {
		return m.Parent
	}
	return nil
}

func (m *Child) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Model)(nil), "Model")
	proto.RegisterType((*Parent)(nil), "Parent")
	proto.RegisterType((*Child)(nil), "Child")
}

func init() { proto.RegisterFile("unsupported.proto", fileDescriptor_86386df630750539) }

var fileDescriptor_86386df630750539 = []byte{
	// 140 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0xcd, 0x2b, 0x2e,
	0x2d, 0x28, 0xc8, 0x2f, 0x2a, 0x49, 0x4d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52, 0xe3,
	0x62, 0xf5, 0xcd, 0x4f, 0x49, 0xcd, 0x11, 0x92, 0xe5, 0xe2, 0x2a, 0xcd, 0xcc, 0x2b, 0x31, 0x36,
	0x8a, 0xcf, 0x4e, 0xad, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0xe2, 0x84, 0x88, 0x78, 0xa7,
	0x56, 0x2a, 0xc9, 0x70, 0xb1, 0x05, 0x24, 0x16, 0xa5, 0xe6, 0x95, 0x08, 0x09, 0x71, 0xb1, 0xe4,
	0x25, 0xe6, 0xa6, 0x82, 0x95, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0x36, 0x5c, 0xac, 0xce, 0x19, 0x99,
	0x39, 0x29, 0x42, 0xf2, 0x5c, 0x6c, 0x05, 0x60, 0x65, 0x60, 0x69, 0x6e, 0x23, 0x76, 0x3d, 0x88,
	0xae, 0x20, 0xa8, 0x30, 0x5c, 0x37, 0x13, 0x42, 0x77, 0x12, 0x1b, 0xd8, 0x29, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x87, 0xa3, 0x21, 0x7d, 0x9f, 0x00, 0x00, 0x00,
}
