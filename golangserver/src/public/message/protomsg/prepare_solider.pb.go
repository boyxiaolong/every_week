// Code generated by protoc-gen-go. DO NOT EDIT.
// source: prepare_solider.proto

package protomsg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// ---------------------------------客户端协议--------------------------------- //
type MsgGS2CLPrepareSoldierNotice struct {
	BaseInfo *PrepareSoldierBaseInfo      `protobuf:"bytes,1,opt,name=base_info,json=baseInfo" json:"base_info,omitempty"`
	WorkInfo *PrepareSoldierTrainWorkInfo `protobuf:"bytes,2,opt,name=work_info,json=workInfo" json:"work_info,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierNotice) Reset()                    { *m = MsgGS2CLPrepareSoldierNotice{} }
func (m *MsgGS2CLPrepareSoldierNotice) String() string            { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierNotice) ProtoMessage()               {}
func (*MsgGS2CLPrepareSoldierNotice) Descriptor() ([]byte, []int) { return fileDescriptor64, []int{0} }

func (m *MsgGS2CLPrepareSoldierNotice) GetBaseInfo() *PrepareSoldierBaseInfo {
	if m != nil {
		return m.BaseInfo
	}
	return nil
}

func (m *MsgGS2CLPrepareSoldierNotice) GetWorkInfo() *PrepareSoldierTrainWorkInfo {
	if m != nil {
		return m.WorkInfo
	}
	return nil
}

// 预备兵信息更新
type MsgGS2CLPrepareSoldierBaseNotice struct {
	BaseInfo *PrepareSoldierBaseInfo `protobuf:"bytes,1,opt,name=base_info,json=baseInfo" json:"base_info,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierBaseNotice) Reset()         { *m = MsgGS2CLPrepareSoldierBaseNotice{} }
func (m *MsgGS2CLPrepareSoldierBaseNotice) String() string { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierBaseNotice) ProtoMessage()    {}
func (*MsgGS2CLPrepareSoldierBaseNotice) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{1}
}

func (m *MsgGS2CLPrepareSoldierBaseNotice) GetBaseInfo() *PrepareSoldierBaseInfo {
	if m != nil {
		return m.BaseInfo
	}
	return nil
}

// 训练信息更新
type MsgGS2CLPrepareSoldierTrainNotice struct {
	Info *PrepareSoldierTrainWorkInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierTrainNotice) Reset()         { *m = MsgGS2CLPrepareSoldierTrainNotice{} }
func (m *MsgGS2CLPrepareSoldierTrainNotice) String() string { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierTrainNotice) ProtoMessage()    {}
func (*MsgGS2CLPrepareSoldierTrainNotice) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{2}
}

func (m *MsgGS2CLPrepareSoldierTrainNotice) GetInfo() *PrepareSoldierTrainWorkInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

// 请求训练
type MsgCL2GSPrepareSoldierTrainRequest struct {
	Count uint64 `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
}

func (m *MsgCL2GSPrepareSoldierTrainRequest) Reset()         { *m = MsgCL2GSPrepareSoldierTrainRequest{} }
func (m *MsgCL2GSPrepareSoldierTrainRequest) String() string { return proto.CompactTextString(m) }
func (*MsgCL2GSPrepareSoldierTrainRequest) ProtoMessage()    {}
func (*MsgCL2GSPrepareSoldierTrainRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{3}
}

func (m *MsgCL2GSPrepareSoldierTrainRequest) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type MsgGS2CLPrepareSoldierTrainReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierTrainReply) Reset()         { *m = MsgGS2CLPrepareSoldierTrainReply{} }
func (m *MsgGS2CLPrepareSoldierTrainReply) String() string { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierTrainReply) ProtoMessage()    {}
func (*MsgGS2CLPrepareSoldierTrainReply) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{4}
}

func (m *MsgGS2CLPrepareSoldierTrainReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

// 立即训练
type MsgCL2GSPrepareSoldierInstantTrainRequest struct {
	Count uint64 `protobuf:"varint,1,opt,name=count" json:"count,omitempty"`
}

func (m *MsgCL2GSPrepareSoldierInstantTrainRequest) Reset() {
	*m = MsgCL2GSPrepareSoldierInstantTrainRequest{}
}
func (m *MsgCL2GSPrepareSoldierInstantTrainRequest) String() string { return proto.CompactTextString(m) }
func (*MsgCL2GSPrepareSoldierInstantTrainRequest) ProtoMessage()    {}
func (*MsgCL2GSPrepareSoldierInstantTrainRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{5}
}

func (m *MsgCL2GSPrepareSoldierInstantTrainRequest) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type MsgGS2CLPrepareSoldierInstantTrainReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierInstantTrainReply) Reset() {
	*m = MsgGS2CLPrepareSoldierInstantTrainReply{}
}
func (m *MsgGS2CLPrepareSoldierInstantTrainReply) String() string { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierInstantTrainReply) ProtoMessage()    {}
func (*MsgGS2CLPrepareSoldierInstantTrainReply) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{6}
}

func (m *MsgGS2CLPrepareSoldierInstantTrainReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

// 立即加速完成
type MsgCL2GSPrepareSoldierInstantSpeedUpRequest struct {
}

func (m *MsgCL2GSPrepareSoldierInstantSpeedUpRequest) Reset() {
	*m = MsgCL2GSPrepareSoldierInstantSpeedUpRequest{}
}
func (m *MsgCL2GSPrepareSoldierInstantSpeedUpRequest) String() string {
	return proto.CompactTextString(m)
}
func (*MsgCL2GSPrepareSoldierInstantSpeedUpRequest) ProtoMessage() {}
func (*MsgCL2GSPrepareSoldierInstantSpeedUpRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{7}
}

type MsgGS2CLPrepareSoldierInstantSpeedUpReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierInstantSpeedUpReply) Reset() {
	*m = MsgGS2CLPrepareSoldierInstantSpeedUpReply{}
}
func (m *MsgGS2CLPrepareSoldierInstantSpeedUpReply) String() string { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierInstantSpeedUpReply) ProtoMessage()    {}
func (*MsgGS2CLPrepareSoldierInstantSpeedUpReply) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{8}
}

func (m *MsgGS2CLPrepareSoldierInstantSpeedUpReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

// 取消训练
type MsgCL2GSPrepareSoldierCancelRequest struct {
}

func (m *MsgCL2GSPrepareSoldierCancelRequest) Reset()         { *m = MsgCL2GSPrepareSoldierCancelRequest{} }
func (m *MsgCL2GSPrepareSoldierCancelRequest) String() string { return proto.CompactTextString(m) }
func (*MsgCL2GSPrepareSoldierCancelRequest) ProtoMessage()    {}
func (*MsgCL2GSPrepareSoldierCancelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{9}
}

type MsgGS2CLPrepareSoldierCancelReply struct {
	ErrorCode int32 `protobuf:"varint,99,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
}

func (m *MsgGS2CLPrepareSoldierCancelReply) Reset()         { *m = MsgGS2CLPrepareSoldierCancelReply{} }
func (m *MsgGS2CLPrepareSoldierCancelReply) String() string { return proto.CompactTextString(m) }
func (*MsgGS2CLPrepareSoldierCancelReply) ProtoMessage()    {}
func (*MsgGS2CLPrepareSoldierCancelReply) Descriptor() ([]byte, []int) {
	return fileDescriptor64, []int{10}
}

func (m *MsgGS2CLPrepareSoldierCancelReply) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

func init() {
	proto.RegisterType((*MsgGS2CLPrepareSoldierNotice)(nil), "protomsg.MsgGS2CLPrepareSoldierNotice")
	proto.RegisterType((*MsgGS2CLPrepareSoldierBaseNotice)(nil), "protomsg.MsgGS2CLPrepareSoldierBaseNotice")
	proto.RegisterType((*MsgGS2CLPrepareSoldierTrainNotice)(nil), "protomsg.MsgGS2CLPrepareSoldierTrainNotice")
	proto.RegisterType((*MsgCL2GSPrepareSoldierTrainRequest)(nil), "protomsg.MsgCL2GSPrepareSoldierTrainRequest")
	proto.RegisterType((*MsgGS2CLPrepareSoldierTrainReply)(nil), "protomsg.MsgGS2CLPrepareSoldierTrainReply")
	proto.RegisterType((*MsgCL2GSPrepareSoldierInstantTrainRequest)(nil), "protomsg.MsgCL2GSPrepareSoldierInstantTrainRequest")
	proto.RegisterType((*MsgGS2CLPrepareSoldierInstantTrainReply)(nil), "protomsg.MsgGS2CLPrepareSoldierInstantTrainReply")
	proto.RegisterType((*MsgCL2GSPrepareSoldierInstantSpeedUpRequest)(nil), "protomsg.MsgCL2GSPrepareSoldierInstantSpeedUpRequest")
	proto.RegisterType((*MsgGS2CLPrepareSoldierInstantSpeedUpReply)(nil), "protomsg.MsgGS2CLPrepareSoldierInstantSpeedUpReply")
	proto.RegisterType((*MsgCL2GSPrepareSoldierCancelRequest)(nil), "protomsg.MsgCL2GSPrepareSoldierCancelRequest")
	proto.RegisterType((*MsgGS2CLPrepareSoldierCancelReply)(nil), "protomsg.MsgGS2CLPrepareSoldierCancelReply")
}

func init() { proto.RegisterFile("prepare_solider.proto", fileDescriptor64) }

var fileDescriptor64 = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x41, 0x4b, 0xf3, 0x40,
	0x10, 0x86, 0xc9, 0x47, 0xfb, 0xd1, 0x8e, 0x17, 0x09, 0x0a, 0x45, 0x14, 0xe2, 0x4a, 0xb1, 0x22,
	0xf6, 0x50, 0x4f, 0x0a, 0x1e, 0xda, 0x1c, 0x6a, 0xa5, 0x15, 0x49, 0x14, 0x6f, 0x86, 0x6d, 0x32,
	0x0d, 0xa1, 0xc9, 0x4e, 0xdc, 0x4d, 0x29, 0xfe, 0x1c, 0xff, 0xa9, 0xb8, 0x49, 0xc5, 0x42, 0xa2,
	0x39, 0x78, 0xca, 0x26, 0x99, 0xf7, 0x7d, 0x9e, 0x85, 0x81, 0xfd, 0x54, 0x62, 0xca, 0x25, 0x7a,
	0x8a, 0xe2, 0x28, 0x40, 0xd9, 0x4f, 0x25, 0x65, 0x64, 0xb6, 0xf4, 0x23, 0x51, 0xe1, 0xc1, 0x6e,
	0xa2, 0x42, 0xcf, 0xa7, 0x24, 0x21, 0x91, 0xff, 0x63, 0xef, 0x06, 0x1c, 0xce, 0x54, 0x38, 0x76,
	0x07, 0xf6, 0xf4, 0x21, 0x4f, 0xbb, 0x14, 0x07, 0x11, 0xca, 0x7b, 0xca, 0x22, 0x1f, 0xcd, 0x1b,
	0x68, 0xcf, 0xb9, 0x42, 0x2f, 0x12, 0x0b, 0xea, 0x18, 0x96, 0xd1, 0xdb, 0x19, 0x58, 0xfd, 0x4d,
	0x61, 0x7f, 0x3b, 0x32, 0xe2, 0x0a, 0x27, 0x62, 0x41, 0x4e, 0x6b, 0x5e, 0x9c, 0xcc, 0x11, 0xb4,
	0xd7, 0x24, 0x97, 0x79, 0xfc, 0x9f, 0x8e, 0x77, 0xab, 0xe2, 0x8f, 0x92, 0x47, 0xe2, 0x99, 0xe4,
	0x32, 0xef, 0x58, 0x17, 0x27, 0xc6, 0xc1, 0x2a, 0x57, 0xfc, 0xe4, 0xfd, 0x89, 0x26, 0x7b, 0x81,
	0xe3, 0x72, 0x84, 0x76, 0x2a, 0x18, 0x57, 0xd0, 0xf8, 0x56, 0x5f, 0xf3, 0x1a, 0x3a, 0xc2, 0xae,
	0x81, 0xcd, 0x54, 0x68, 0x4f, 0x07, 0x63, 0xb7, 0x64, 0xd8, 0xc1, 0xd7, 0x15, 0xaa, 0xcc, 0xdc,
	0x83, 0xa6, 0x4f, 0x2b, 0x91, 0x69, 0x42, 0xc3, 0xc9, 0x5f, 0xd8, 0xb0, 0xea, 0xfa, 0x45, 0x36,
	0x8d, 0xdf, 0xcc, 0x23, 0x00, 0x94, 0x92, 0xa4, 0xe7, 0x53, 0x80, 0x1d, 0xdf, 0x32, 0x7a, 0x4d,
	0xa7, 0xad, 0xbf, 0xd8, 0x14, 0x20, 0x1b, 0xc2, 0x59, 0x39, 0x7e, 0x22, 0x54, 0xc6, 0x45, 0x56,
	0xc3, 0xe2, 0x16, 0x4e, 0xcb, 0x2d, 0xb6, 0x2b, 0x6a, 0xc8, 0x5c, 0xc0, 0xf9, 0x8f, 0x32, 0x6e,
	0x8a, 0x18, 0x3c, 0xa5, 0x85, 0x0e, 0xbb, 0xd3, 0xee, 0xd5, 0xe0, 0xaf, 0xf1, 0x1a, 0xe8, 0x2e,
	0x9c, 0x94, 0xa3, 0x6d, 0x2e, 0x7c, 0x8c, 0x37, 0xc8, 0x51, 0xd5, 0x36, 0x6c, 0xc6, 0x7e, 0x47,
	0xcd, 0xff, 0xeb, 0xed, 0xb8, 0xfc, 0x08, 0x00, 0x00, 0xff, 0xff, 0xaa, 0xc3, 0x80, 0x89, 0x94,
	0x03, 0x00, 0x00,
}
