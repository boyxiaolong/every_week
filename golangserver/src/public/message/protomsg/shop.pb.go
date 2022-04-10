// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shop.proto

package protomsg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MsgCL2GSShopBuyItemRequest struct {
	Type  uint32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	Id    uint32 `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	Count uint32 `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
}

func (m *MsgCL2GSShopBuyItemRequest) Reset()                    { *m = MsgCL2GSShopBuyItemRequest{} }
func (m *MsgCL2GSShopBuyItemRequest) String() string            { return proto.CompactTextString(m) }
func (*MsgCL2GSShopBuyItemRequest) ProtoMessage()               {}
func (*MsgCL2GSShopBuyItemRequest) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{0} }

func (m *MsgCL2GSShopBuyItemRequest) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *MsgCL2GSShopBuyItemRequest) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MsgCL2GSShopBuyItemRequest) GetCount() uint32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type MsgGS2CLShopBuyItemReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLShopBuyItemReply) Reset()                    { *m = MsgGS2CLShopBuyItemReply{} }
func (m *MsgGS2CLShopBuyItemReply) String() string            { return proto.CompactTextString(m) }
func (*MsgGS2CLShopBuyItemReply) ProtoMessage()               {}
func (*MsgGS2CLShopBuyItemReply) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{1} }

func (m *MsgGS2CLShopBuyItemReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

type MsgCL2GSShopRefreshRequest struct {
	Type uint32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
}

func (m *MsgCL2GSShopRefreshRequest) Reset()                    { *m = MsgCL2GSShopRefreshRequest{} }
func (m *MsgCL2GSShopRefreshRequest) String() string            { return proto.CompactTextString(m) }
func (*MsgCL2GSShopRefreshRequest) ProtoMessage()               {}
func (*MsgCL2GSShopRefreshRequest) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{2} }

func (m *MsgCL2GSShopRefreshRequest) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

type MsgGS2CLShopRefreshReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLShopRefreshReply) Reset()                    { *m = MsgGS2CLShopRefreshReply{} }
func (m *MsgGS2CLShopRefreshReply) String() string            { return proto.CompactTextString(m) }
func (*MsgGS2CLShopRefreshReply) ProtoMessage()               {}
func (*MsgGS2CLShopRefreshReply) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{3} }

func (m *MsgGS2CLShopRefreshReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

// 普通商店通知
type MsgGS2CLShopNotice struct {
	ShopData *ShopData `protobuf:"bytes,1,opt,name=shop_data,json=shopData" json:"shop_data,omitempty"`
}

func (m *MsgGS2CLShopNotice) Reset()                    { *m = MsgGS2CLShopNotice{} }
func (m *MsgGS2CLShopNotice) String() string            { return proto.CompactTextString(m) }
func (*MsgGS2CLShopNotice) ProtoMessage()               {}
func (*MsgGS2CLShopNotice) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{4} }

func (m *MsgGS2CLShopNotice) GetShopData() *ShopData {
	if m != nil {
		return m.ShopData
	}
	return nil
}

type MsgGS2CLShopManNotice struct {
	ShopData *MysteryShopManData `protobuf:"bytes,1,opt,name=shop_data,json=shopData" json:"shop_data,omitempty"`
}

func (m *MsgGS2CLShopManNotice) Reset()                    { *m = MsgGS2CLShopManNotice{} }
func (m *MsgGS2CLShopManNotice) String() string            { return proto.CompactTextString(m) }
func (*MsgGS2CLShopManNotice) ProtoMessage()               {}
func (*MsgGS2CLShopManNotice) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{5} }

func (m *MsgGS2CLShopManNotice) GetShopData() *MysteryShopManData {
	if m != nil {
		return m.ShopData
	}
	return nil
}

type MsgCL2GSShopManTalkRequest struct {
}

func (m *MsgCL2GSShopManTalkRequest) Reset()                    { *m = MsgCL2GSShopManTalkRequest{} }
func (m *MsgCL2GSShopManTalkRequest) String() string            { return proto.CompactTextString(m) }
func (*MsgCL2GSShopManTalkRequest) ProtoMessage()               {}
func (*MsgCL2GSShopManTalkRequest) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{6} }

type MsgGS2CLShopManTalkReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLShopManTalkReply) Reset()                    { *m = MsgGS2CLShopManTalkReply{} }
func (m *MsgGS2CLShopManTalkReply) String() string            { return proto.CompactTextString(m) }
func (*MsgGS2CLShopManTalkReply) ProtoMessage()               {}
func (*MsgGS2CLShopManTalkReply) Descriptor() ([]byte, []int) { return fileDescriptor73, []int{7} }

func (m *MsgGS2CLShopManTalkReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

func init() {
	proto.RegisterType((*MsgCL2GSShopBuyItemRequest)(nil), "protomsg.MsgCL2GSShopBuyItemRequest")
	proto.RegisterType((*MsgGS2CLShopBuyItemReply)(nil), "protomsg.MsgGS2CLShopBuyItemReply")
	proto.RegisterType((*MsgCL2GSShopRefreshRequest)(nil), "protomsg.MsgCL2GSShopRefreshRequest")
	proto.RegisterType((*MsgGS2CLShopRefreshReply)(nil), "protomsg.MsgGS2CLShopRefreshReply")
	proto.RegisterType((*MsgGS2CLShopNotice)(nil), "protomsg.MsgGS2CLShopNotice")
	proto.RegisterType((*MsgGS2CLShopManNotice)(nil), "protomsg.MsgGS2CLShopManNotice")
	proto.RegisterType((*MsgCL2GSShopManTalkRequest)(nil), "protomsg.MsgCL2GSShopManTalkRequest")
	proto.RegisterType((*MsgGS2CLShopManTalkReply)(nil), "protomsg.MsgGS2CLShopManTalkReply")
}

func init() { proto.RegisterFile("shop.proto", fileDescriptor73) }

var fileDescriptor73 = []byte{
	// 286 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xc1, 0x4f, 0xb3, 0x40,
	0x10, 0xc5, 0x03, 0xdf, 0x57, 0xd3, 0x8e, 0xd1, 0xc3, 0xc6, 0x26, 0xa4, 0xa9, 0x49, 0xc3, 0xa9,
	0x27, 0x34, 0x78, 0xea, 0x55, 0x34, 0x8d, 0x49, 0xf1, 0xb0, 0x18, 0xaf, 0x64, 0x85, 0x11, 0x88,
	0x85, 0xc1, 0xdd, 0xe5, 0xc0, 0x7f, 0x6f, 0xba, 0x62, 0x5b, 0x30, 0x51, 0x4f, 0x0c, 0xb3, 0xf3,
	0x7e, 0x6f, 0xe6, 0x01, 0xa8, 0x9c, 0x6a, 0xaf, 0x96, 0xa4, 0x89, 0x8d, 0xcd, 0xa7, 0x54, 0xd9,
	0x6c, 0x5a, 0xaa, 0x2c, 0x4e, 0xa8, 0x2c, 0xa9, 0x8a, 0x0f, 0x03, 0xee, 0x33, 0xcc, 0x42, 0x95,
	0x05, 0x1b, 0x7f, 0x1d, 0x45, 0x39, 0xd5, 0xb7, 0x4d, 0xfb, 0xa0, 0xb1, 0xe4, 0xf8, 0xde, 0xa0,
	0xd2, 0x8c, 0xc1, 0x7f, 0xdd, 0xd6, 0xe8, 0x58, 0x0b, 0x6b, 0x79, 0xc6, 0x4d, 0xcd, 0xce, 0xc1,
	0x2e, 0x52, 0xc7, 0x36, 0x1d, 0xbb, 0x48, 0xd9, 0x05, 0x8c, 0x12, 0x6a, 0x2a, 0xed, 0xfc, 0x33,
	0xad, 0xcf, 0x1f, 0x77, 0x05, 0x4e, 0xa8, 0xb2, 0x75, 0xe4, 0x07, 0x9b, 0x1e, 0xb7, 0xde, 0xb6,
	0xec, 0x12, 0x00, 0xa5, 0x24, 0x19, 0x27, 0x94, 0xa2, 0x93, 0x2c, 0xac, 0xe5, 0x88, 0x4f, 0x4c,
	0x27, 0xa0, 0x14, 0xdd, 0xeb, 0xfe, 0x4a, 0x1c, 0x5f, 0x25, 0xaa, 0xfc, 0x87, 0x95, 0x86, 0x66,
	0x7b, 0xc5, 0x1f, 0xcc, 0xee, 0x81, 0x1d, 0x4b, 0x1f, 0x49, 0x17, 0x09, 0xb2, 0x2b, 0x98, 0xec,
	0x32, 0x8a, 0x53, 0xa1, 0x85, 0x71, 0x3a, 0xf5, 0x99, 0xf7, 0x15, 0xa5, 0xb7, 0x1b, 0xbc, 0x13,
	0x5a, 0xf0, 0xb1, 0xea, 0x2a, 0x97, 0xc3, 0xf4, 0x18, 0x13, 0x8a, 0xaa, 0x23, 0xad, 0xbe, 0x93,
	0xe6, 0x07, 0x52, 0xd8, 0x2a, 0x8d, 0xb2, 0xed, 0x24, 0x03, 0xe6, 0xbc, 0x9f, 0x43, 0x28, 0xaa,
	0x27, 0xb1, 0x7d, 0xeb, 0x72, 0x18, 0xde, 0xbc, 0x7f, 0xfd, 0xfd, 0xe6, 0x97, 0x13, 0xe3, 0x7f,
	0xf3, 0x11, 0x00, 0x00, 0xff, 0xff, 0x3e, 0x0d, 0x4b, 0x3e, 0x29, 0x02, 0x00, 0x00,
}
