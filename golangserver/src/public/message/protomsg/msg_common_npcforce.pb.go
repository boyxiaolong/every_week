// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg_common_npcforce.proto

package protomsg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type NpcForceEvent int32

const (
	NpcForceEvent_kNpcForceEvent_None             NpcForceEvent = 0
	NpcForceEvent_kNpcForceEvent_AttackCity       NpcForceEvent = 1
	NpcForceEvent_kNpcForceEvent_AttackStronghold NpcForceEvent = 2
)

var NpcForceEvent_name = map[int32]string{
	0: "kNpcForceEvent_None",
	1: "kNpcForceEvent_AttackCity",
	2: "kNpcForceEvent_AttackStronghold",
}
var NpcForceEvent_value = map[string]int32{
	"kNpcForceEvent_None":             0,
	"kNpcForceEvent_AttackCity":       1,
	"kNpcForceEvent_AttackStronghold": 2,
}

func (x NpcForceEvent) String() string {
	return proto.EnumName(NpcForceEvent_name, int32(x))
}
func (NpcForceEvent) EnumDescriptor() ([]byte, []int) { return fileDescriptor51, []int{0} }

type NpcForceDataType int32

const (
	NpcForceDataType_kNpcForceData_None       NpcForceDataType = 0
	NpcForceDataType_kNpcForceData_FortNumber NpcForceDataType = 1
)

var NpcForceDataType_name = map[int32]string{
	0: "kNpcForceData_None",
	1: "kNpcForceData_FortNumber",
}
var NpcForceDataType_value = map[string]int32{
	"kNpcForceData_None":       0,
	"kNpcForceData_FortNumber": 1,
}

func (x NpcForceDataType) String() string {
	return proto.EnumName(NpcForceDataType_name, int32(x))
}
func (NpcForceDataType) EnumDescriptor() ([]byte, []int) { return fileDescriptor51, []int{1} }

type NpcForce struct {
	NpcForceId        uint64   `protobuf:"varint,1,opt,name=npc_force_id,json=npcForceId" json:"npc_force_id,omitempty"`
	NpcForceName      string   `protobuf:"bytes,2,opt,name=npc_force_name,json=npcForceName" json:"npc_force_name,omitempty"`
	NpcForceShortName string   `protobuf:"bytes,3,opt,name=npc_force_short_name,json=npcForceShortName" json:"npc_force_short_name,omitempty"`
	NpcForceFlag      uint32   `protobuf:"varint,4,opt,name=npc_force_flag,json=npcForceFlag" json:"npc_force_flag,omitempty"`
	DramaStoryId      uint32   `protobuf:"varint,6,opt,name=drama_story_id,json=dramaStoryId" json:"drama_story_id,omitempty"`
	Citys             []uint64 `protobuf:"varint,5,rep,packed,name=citys" json:"citys,omitempty"`
}

func (m *NpcForce) Reset()                    { *m = NpcForce{} }
func (m *NpcForce) String() string            { return proto.CompactTextString(m) }
func (*NpcForce) ProtoMessage()               {}
func (*NpcForce) Descriptor() ([]byte, []int) { return fileDescriptor51, []int{0} }

func (m *NpcForce) GetNpcForceId() uint64 {
	if m != nil {
		return m.NpcForceId
	}
	return 0
}

func (m *NpcForce) GetNpcForceName() string {
	if m != nil {
		return m.NpcForceName
	}
	return ""
}

func (m *NpcForce) GetNpcForceShortName() string {
	if m != nil {
		return m.NpcForceShortName
	}
	return ""
}

func (m *NpcForce) GetNpcForceFlag() uint32 {
	if m != nil {
		return m.NpcForceFlag
	}
	return 0
}

func (m *NpcForce) GetDramaStoryId() uint32 {
	if m != nil {
		return m.DramaStoryId
	}
	return 0
}

func (m *NpcForce) GetCitys() []uint64 {
	if m != nil {
		return m.Citys
	}
	return nil
}

type NpcForceFortData struct {
	Num     uint32 `protobuf:"varint,1,opt,name=num" json:"num,omitempty"`
	FortId  uint32 `protobuf:"varint,2,opt,name=fort_id,json=fortId" json:"fort_id,omitempty"`
	EventId uint32 `protobuf:"varint,3,opt,name=event_id,json=eventId" json:"event_id,omitempty"`
}

func (m *NpcForceFortData) Reset()                    { *m = NpcForceFortData{} }
func (m *NpcForceFortData) String() string            { return proto.CompactTextString(m) }
func (*NpcForceFortData) ProtoMessage()               {}
func (*NpcForceFortData) Descriptor() ([]byte, []int) { return fileDescriptor51, []int{1} }

func (m *NpcForceFortData) GetNum() uint32 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *NpcForceFortData) GetFortId() uint32 {
	if m != nil {
		return m.FortId
	}
	return 0
}

func (m *NpcForceFortData) GetEventId() uint32 {
	if m != nil {
		return m.EventId
	}
	return 0
}

func init() {
	proto.RegisterType((*NpcForce)(nil), "protomsg.npc_force")
	proto.RegisterType((*NpcForceFortData)(nil), "protomsg.NpcForceFortData")
	proto.RegisterEnum("protomsg.NpcForceEvent", NpcForceEvent_name, NpcForceEvent_value)
	proto.RegisterEnum("protomsg.NpcForceDataType", NpcForceDataType_name, NpcForceDataType_value)
}

func init() { proto.RegisterFile("msg_common_npcforce.proto", fileDescriptor51) }

var fileDescriptor51 = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcd, 0x6e, 0xa3, 0x30,
	0x1c, 0xc4, 0x97, 0x90, 0xcf, 0xbf, 0xc2, 0x8a, 0xf5, 0x46, 0x1b, 0x22, 0x6d, 0x55, 0x94, 0xf6,
	0x80, 0x72, 0x68, 0x0f, 0x7d, 0x82, 0xaa, 0x6d, 0x54, 0x2e, 0x1c, 0x48, 0x0f, 0xbd, 0x21, 0xc7,
	0x76, 0x08, 0x0a, 0xb6, 0x91, 0x71, 0x2a, 0xf1, 0xc2, 0x7d, 0x8e, 0xca, 0x8e, 0x28, 0x4d, 0xd5,
	0x13, 0xcc, 0xcc, 0x8f, 0xc1, 0x23, 0x19, 0x16, 0xbc, 0xce, 0x33, 0x22, 0x39, 0x97, 0x22, 0x13,
	0x15, 0xd9, 0x49, 0x45, 0xd8, 0x4d, 0xa5, 0xa4, 0x96, 0x68, 0x6c, 0x1f, 0xbc, 0xce, 0x97, 0xef,
	0x0e, 0x4c, 0x44, 0x45, 0x32, 0x9b, 0xa2, 0x10, 0xa6, 0x9f, 0x22, 0x2b, 0x68, 0xe0, 0x84, 0x4e,
	0xd4, 0x4f, 0x41, 0x54, 0x64, 0x6d, 0xac, 0x98, 0xa2, 0x6b, 0xf8, 0xdd, 0x11, 0x02, 0x73, 0x16,
	0xf4, 0x42, 0x27, 0x9a, 0xa4, 0xd3, 0x96, 0x49, 0x30, 0x67, 0xe8, 0x16, 0x66, 0x1d, 0x55, 0xef,
	0xa5, 0xd2, 0x27, 0xd6, 0xb5, 0xec, 0x9f, 0x96, 0xdd, 0x98, 0xc4, 0x7e, 0x70, 0x56, 0xbb, 0x2b,
	0x71, 0x1e, 0xf4, 0x43, 0x27, 0xf2, 0xba, 0xda, 0x75, 0x89, 0x73, 0x43, 0x51, 0x85, 0x39, 0xce,
	0x6a, 0x2d, 0x55, 0x63, 0x0e, 0x38, 0x3c, 0x51, 0xd6, 0xdd, 0x18, 0x33, 0xa6, 0x68, 0x06, 0x03,
	0x52, 0xe8, 0xa6, 0x0e, 0x06, 0xa1, 0x1b, 0xf5, 0xd3, 0x93, 0x58, 0xbe, 0x82, 0x9f, 0xb4, 0x5d,
	0x52, 0xe9, 0x47, 0xac, 0x31, 0xf2, 0xc1, 0x15, 0x47, 0x6e, 0x57, 0x7a, 0xa9, 0x79, 0x45, 0x73,
	0x18, 0xed, 0xcc, 0x69, 0x0b, 0x6a, 0x77, 0x79, 0xe9, 0xd0, 0xc8, 0x98, 0xa2, 0x05, 0x8c, 0xd9,
	0x1b, 0x13, 0x36, 0x71, 0x6d, 0x32, 0xb2, 0x3a, 0xa6, 0xab, 0x12, 0xbc, 0xb6, 0xf9, 0xc9, 0x58,
	0x68, 0x0e, 0x7f, 0x0f, 0x67, 0x4e, 0x96, 0x48, 0xc1, 0xfc, 0x5f, 0xe8, 0x02, 0x16, 0xdf, 0x82,
	0x7b, 0xad, 0x31, 0x39, 0x3c, 0x14, 0xba, 0xf1, 0x1d, 0x74, 0x05, 0x97, 0x3f, 0xc6, 0x1b, 0xad,
	0xa4, 0xc8, 0xf7, 0xb2, 0xa4, 0x7e, 0x6f, 0xf5, 0xdc, 0xed, 0x30, 0x1b, 0x5e, 0x9a, 0x8a, 0xa1,
	0x7f, 0x80, 0x0e, 0x5f, 0xcd, 0xf6, 0x7f, 0xff, 0x21, 0x38, 0xf7, 0xcd, 0xf2, 0xe4, 0xc8, 0xb7,
	0x4c, 0xf9, 0xce, 0x76, 0x68, 0x2f, 0xc1, 0xdd, 0x47, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe2, 0xca,
	0x21, 0xb8, 0x28, 0x02, 0x00, 0x00,
}
