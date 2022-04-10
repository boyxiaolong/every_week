// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg_common_radar.proto

package protomsg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 雷达状态类型
type RadarStatus int32

const (
	RadarStatus_kRadarStatus_None     RadarStatus = 0
	RadarStatus_kRadarStatus_Accept   RadarStatus = 1
	RadarStatus_kRadarStatus_Complete RadarStatus = 2
	RadarStatus_kRadarStatus_Destroy  RadarStatus = 3
)

var RadarStatus_name = map[int32]string{
	0: "kRadarStatus_None",
	1: "kRadarStatus_Accept",
	2: "kRadarStatus_Complete",
	3: "kRadarStatus_Destroy",
}
var RadarStatus_value = map[string]int32{
	"kRadarStatus_None":     0,
	"kRadarStatus_Accept":   1,
	"kRadarStatus_Complete": 2,
	"kRadarStatus_Destroy":  3,
}

func (x RadarStatus) String() string {
	return proto.EnumName(RadarStatus_name, int32(x))
}
func (RadarStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{0} }

// 雷达时间类型
type RadarTimeType int32

const (
	RadarTimeType_kRadarTimeType_None   RadarTimeType = 0
	RadarTimeType_kRadarTimeType_Daily  RadarTimeType = 1
	RadarTimeType_kRadarTimeType_Weekly RadarTimeType = 2
	RadarTimeType_kRadarTimeType_Yearly RadarTimeType = 3
	RadarTimeType_kRadarTimeType_Once   RadarTimeType = 4
)

var RadarTimeType_name = map[int32]string{
	0: "kRadarTimeType_None",
	1: "kRadarTimeType_Daily",
	2: "kRadarTimeType_Weekly",
	3: "kRadarTimeType_Yearly",
	4: "kRadarTimeType_Once",
}
var RadarTimeType_value = map[string]int32{
	"kRadarTimeType_None":   0,
	"kRadarTimeType_Daily":  1,
	"kRadarTimeType_Weekly": 2,
	"kRadarTimeType_Yearly": 3,
	"kRadarTimeType_Once":   4,
}

func (x RadarTimeType) String() string {
	return proto.EnumName(RadarTimeType_name, int32(x))
}
func (RadarTimeType) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{1} }

// 雷达处理类型
type RadarType int32

const (
	RadarType_kRadarType_None     RadarType = 0
	RadarType_kRadarType_Monster  RadarType = 1
	RadarType_kRadarType_Common   RadarType = 2
	RadarType_kRadarType_Dialogue RadarType = 3
)

var RadarType_name = map[int32]string{
	0: "kRadarType_None",
	1: "kRadarType_Monster",
	2: "kRadarType_Common",
	3: "kRadarType_Dialogue",
}
var RadarType_value = map[string]int32{
	"kRadarType_None":     0,
	"kRadarType_Monster":  1,
	"kRadarType_Common":   2,
	"kRadarType_Dialogue": 3,
}

func (x RadarType) String() string {
	return proto.EnumName(RadarType_name, int32(x))
}
func (RadarType) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{2} }

// 怪物处理类型的雷达处理子类型
type RadarMonsterSubType int32

const (
	RadarMonsterSubType_kRadarMonsterSubType_None   RadarMonsterSubType = 0
	RadarMonsterSubType_kRadarMonsterSubType_Own    RadarMonsterSubType = 1
	RadarMonsterSubType_kRadarMonsterSubType_Common RadarMonsterSubType = 2
)

var RadarMonsterSubType_name = map[int32]string{
	0: "kRadarMonsterSubType_None",
	1: "kRadarMonsterSubType_Own",
	2: "kRadarMonsterSubType_Common",
}
var RadarMonsterSubType_value = map[string]int32{
	"kRadarMonsterSubType_None":   0,
	"kRadarMonsterSubType_Own":    1,
	"kRadarMonsterSubType_Common": 2,
}

func (x RadarMonsterSubType) String() string {
	return proto.EnumName(RadarMonsterSubType_name, int32(x))
}
func (RadarMonsterSubType) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{3} }

// 普通运输处理类型的雷达处理子类型
type RadarCommonSubType int32

const (
	RadarCommonSubType_kRadarCommonSubType_None    RadarCommonSubType = 0
	RadarCommonSubType_kRadarCommonSubType_Assist  RadarCommonSubType = 1
	RadarCommonSubType_kRadarCommonSubType_Rescue  RadarCommonSubType = 2
	RadarCommonSubType_kRadarCommonSubType_Supply  RadarCommonSubType = 3
	RadarCommonSubType_kRadarCommonSubType_Explore RadarCommonSubType = 4
)

var RadarCommonSubType_name = map[int32]string{
	0: "kRadarCommonSubType_None",
	1: "kRadarCommonSubType_Assist",
	2: "kRadarCommonSubType_Rescue",
	3: "kRadarCommonSubType_Supply",
	4: "kRadarCommonSubType_Explore",
}
var RadarCommonSubType_value = map[string]int32{
	"kRadarCommonSubType_None":    0,
	"kRadarCommonSubType_Assist":  1,
	"kRadarCommonSubType_Rescue":  2,
	"kRadarCommonSubType_Supply":  3,
	"kRadarCommonSubType_Explore": 4,
}

func (x RadarCommonSubType) String() string {
	return proto.EnumName(RadarCommonSubType_name, int32(x))
}
func (RadarCommonSubType) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{4} }

// 雷达品质
type RadarQuality int32

const (
	RadarQuality_kRadarQuality_None   RadarQuality = 0
	RadarQuality_kRadarQuality_Green  RadarQuality = 1
	RadarQuality_kRadarQuality_Blue   RadarQuality = 2
	RadarQuality_kRadarQuality_Purple RadarQuality = 3
	RadarQuality_kRadarQuality_Orange RadarQuality = 4
	RadarQuality_kRadarQuality_Red    RadarQuality = 5
)

var RadarQuality_name = map[int32]string{
	0: "kRadarQuality_None",
	1: "kRadarQuality_Green",
	2: "kRadarQuality_Blue",
	3: "kRadarQuality_Purple",
	4: "kRadarQuality_Orange",
	5: "kRadarQuality_Red",
}
var RadarQuality_value = map[string]int32{
	"kRadarQuality_None":   0,
	"kRadarQuality_Green":  1,
	"kRadarQuality_Blue":   2,
	"kRadarQuality_Purple": 3,
	"kRadarQuality_Orange": 4,
	"kRadarQuality_Red":    5,
}

func (x RadarQuality) String() string {
	return proto.EnumName(RadarQuality_name, int32(x))
}
func (RadarQuality) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{5} }

// 雷达销毁原因
type RadarDestroyReason int32

const (
	RadarDestroyReason_kRadarDestroyReason_None         RadarDestroyReason = 0
	RadarDestroyReason_kRadarDestroyReason_Reward       RadarDestroyReason = 1
	RadarDestroyReason_kRadarDestroyReason_Timeout      RadarDestroyReason = 2
	RadarDestroyReason_kRadarDestroyReason_CastleChange RadarDestroyReason = 3
	RadarDestroyReason_kRadarDestroyReason_DayReset     RadarDestroyReason = 4
)

var RadarDestroyReason_name = map[int32]string{
	0: "kRadarDestroyReason_None",
	1: "kRadarDestroyReason_Reward",
	2: "kRadarDestroyReason_Timeout",
	3: "kRadarDestroyReason_CastleChange",
	4: "kRadarDestroyReason_DayReset",
}
var RadarDestroyReason_value = map[string]int32{
	"kRadarDestroyReason_None":         0,
	"kRadarDestroyReason_Reward":       1,
	"kRadarDestroyReason_Timeout":      2,
	"kRadarDestroyReason_CastleChange": 3,
	"kRadarDestroyReason_DayReset":     4,
}

func (x RadarDestroyReason) String() string {
	return proto.EnumName(RadarDestroyReason_name, int32(x))
}
func (RadarDestroyReason) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{6} }

// 雷达完成类型
type RadarCompleteType int32

const (
	RadarCompleteType_kRadarCompleteType_None        RadarCompleteType = 0
	RadarCompleteType_kRadarCompleteType_QueueAccpet RadarCompleteType = 1
)

var RadarCompleteType_name = map[int32]string{
	0: "kRadarCompleteType_None",
	1: "kRadarCompleteType_QueueAccpet",
}
var RadarCompleteType_value = map[string]int32{
	"kRadarCompleteType_None":        0,
	"kRadarCompleteType_QueueAccpet": 1,
}

func (x RadarCompleteType) String() string {
	return proto.EnumName(RadarCompleteType_name, int32(x))
}
func (RadarCompleteType) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{7} }

// 虚拟队列类型
type VirtualQueueType int32

const (
	VirtualQueueType_kVirtualQueueType_None  VirtualQueueType = 0
	VirtualQueueType_kVirtualQueueType_Radar VirtualQueueType = 1
)

var VirtualQueueType_name = map[int32]string{
	0: "kVirtualQueueType_None",
	1: "kVirtualQueueType_Radar",
}
var VirtualQueueType_value = map[string]int32{
	"kVirtualQueueType_None":  0,
	"kVirtualQueueType_Radar": 1,
}

func (x VirtualQueueType) String() string {
	return proto.EnumName(VirtualQueueType_name, int32(x))
}
func (VirtualQueueType) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{8} }

// 虚拟队列状态
type VirtualQueueStatus int32

const (
	VirtualQueueStatus_kVirtualQueueStatus_None     VirtualQueueStatus = 0
	VirtualQueueStatus_kVirtualQueueStatus_Accpet   VirtualQueueStatus = 1
	VirtualQueueStatus_kVirtualQueueStatus_Complete VirtualQueueStatus = 2
)

var VirtualQueueStatus_name = map[int32]string{
	0: "kVirtualQueueStatus_None",
	1: "kVirtualQueueStatus_Accpet",
	2: "kVirtualQueueStatus_Complete",
}
var VirtualQueueStatus_value = map[string]int32{
	"kVirtualQueueStatus_None":     0,
	"kVirtualQueueStatus_Accpet":   1,
	"kVirtualQueueStatus_Complete": 2,
}

func (x VirtualQueueStatus) String() string {
	return proto.EnumName(VirtualQueueStatus_name, int32(x))
}
func (VirtualQueueStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor53, []int{9} }

// 雷达数据
type RadarData struct {
	RadarId       uint32    `protobuf:"varint,1,opt,name=radar_id,json=radarId" json:"radar_id,omitempty"`
	RadarConfigId uint32    `protobuf:"varint,2,opt,name=radar_config_id,json=radarConfigId" json:"radar_config_id,omitempty"`
	Status        int32     `protobuf:"varint,3,opt,name=status" json:"status,omitempty"`
	Param         int64     `protobuf:"varint,4,opt,name=param" json:"param,omitempty"`
	Param2        int64     `protobuf:"varint,5,opt,name=param2" json:"param2,omitempty"`
	Position      *Vector2D `protobuf:"bytes,6,opt,name=position" json:"position,omitempty"`
	BeginTime     int64     `protobuf:"varint,7,opt,name=begin_time,json=beginTime" json:"begin_time,omitempty"`
	EndTime       int64     `protobuf:"varint,8,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	StartTime     int64     `protobuf:"varint,9,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	StopTime      int64     `protobuf:"varint,10,opt,name=stop_time,json=stopTime" json:"stop_time,omitempty"`
	RewardId      uint32    `protobuf:"varint,11,opt,name=reward_id,json=rewardId" json:"reward_id,omitempty"`
	Give          bool      `protobuf:"varint,12,opt,name=give" json:"give,omitempty"`
	DialogIndex   int32     `protobuf:"varint,13,opt,name=dialog_index,json=dialogIndex" json:"dialog_index,omitempty"`
	ShowId        int32     `protobuf:"varint,14,opt,name=show_id,json=showId" json:"show_id,omitempty"`
}

func (m *RadarData) Reset()                    { *m = RadarData{} }
func (m *RadarData) String() string            { return proto.CompactTextString(m) }
func (*RadarData) ProtoMessage()               {}
func (*RadarData) Descriptor() ([]byte, []int) { return fileDescriptor53, []int{0} }

func (m *RadarData) GetRadarId() uint32 {
	if m != nil {
		return m.RadarId
	}
	return 0
}

func (m *RadarData) GetRadarConfigId() uint32 {
	if m != nil {
		return m.RadarConfigId
	}
	return 0
}

func (m *RadarData) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *RadarData) GetParam() int64 {
	if m != nil {
		return m.Param
	}
	return 0
}

func (m *RadarData) GetParam2() int64 {
	if m != nil {
		return m.Param2
	}
	return 0
}

func (m *RadarData) GetPosition() *Vector2D {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *RadarData) GetBeginTime() int64 {
	if m != nil {
		return m.BeginTime
	}
	return 0
}

func (m *RadarData) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *RadarData) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *RadarData) GetStopTime() int64 {
	if m != nil {
		return m.StopTime
	}
	return 0
}

func (m *RadarData) GetRewardId() uint32 {
	if m != nil {
		return m.RewardId
	}
	return 0
}

func (m *RadarData) GetGive() bool {
	if m != nil {
		return m.Give
	}
	return false
}

func (m *RadarData) GetDialogIndex() int32 {
	if m != nil {
		return m.DialogIndex
	}
	return 0
}

func (m *RadarData) GetShowId() int32 {
	if m != nil {
		return m.ShowId
	}
	return 0
}

// 雷达村庄数据
type RadarVillageData struct {
	ValliageId uint64    `protobuf:"varint,1,opt,name=valliage_id,json=valliageId" json:"valliage_id,omitempty"`
	Position   *Vector2D `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
}

func (m *RadarVillageData) Reset()                    { *m = RadarVillageData{} }
func (m *RadarVillageData) String() string            { return proto.CompactTextString(m) }
func (*RadarVillageData) ProtoMessage()               {}
func (*RadarVillageData) Descriptor() ([]byte, []int) { return fileDescriptor53, []int{1} }

func (m *RadarVillageData) GetValliageId() uint64 {
	if m != nil {
		return m.ValliageId
	}
	return 0
}

func (m *RadarVillageData) GetPosition() *Vector2D {
	if m != nil {
		return m.Position
	}
	return nil
}

// 雷达据点数据
type RadarStrongholdData struct {
	StrongholdId uint64    `protobuf:"varint,1,opt,name=stronghold_id,json=strongholdId" json:"stronghold_id,omitempty"`
	Position     *Vector2D `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
}

func (m *RadarStrongholdData) Reset()                    { *m = RadarStrongholdData{} }
func (m *RadarStrongholdData) String() string            { return proto.CompactTextString(m) }
func (*RadarStrongholdData) ProtoMessage()               {}
func (*RadarStrongholdData) Descriptor() ([]byte, []int) { return fileDescriptor53, []int{2} }

func (m *RadarStrongholdData) GetStrongholdId() uint64 {
	if m != nil {
		return m.StrongholdId
	}
	return 0
}

func (m *RadarStrongholdData) GetPosition() *Vector2D {
	if m != nil {
		return m.Position
	}
	return nil
}

// 虚拟队列
type VirtualQueueData struct {
	QueueId   uint32      `protobuf:"varint,1,opt,name=queue_id,json=queueId" json:"queue_id,omitempty"`
	StartTime int64       `protobuf:"varint,2,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	Type      int32       `protobuf:"varint,3,opt,name=type" json:"type,omitempty"`
	Param1    int64       `protobuf:"varint,4,opt,name=param1" json:"param1,omitempty"`
	Param2    int64       `protobuf:"varint,5,opt,name=param2" json:"param2,omitempty"`
	Param3    int64       `protobuf:"varint,6,opt,name=param3" json:"param3,omitempty"`
	Paths     []*Vector2D `protobuf:"bytes,7,rep,name=paths" json:"paths,omitempty"`
	Status    int32       `protobuf:"varint,8,opt,name=status" json:"status,omitempty"`
}

func (m *VirtualQueueData) Reset()                    { *m = VirtualQueueData{} }
func (m *VirtualQueueData) String() string            { return proto.CompactTextString(m) }
func (*VirtualQueueData) ProtoMessage()               {}
func (*VirtualQueueData) Descriptor() ([]byte, []int) { return fileDescriptor53, []int{3} }

func (m *VirtualQueueData) GetQueueId() uint32 {
	if m != nil {
		return m.QueueId
	}
	return 0
}

func (m *VirtualQueueData) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *VirtualQueueData) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *VirtualQueueData) GetParam1() int64 {
	if m != nil {
		return m.Param1
	}
	return 0
}

func (m *VirtualQueueData) GetParam2() int64 {
	if m != nil {
		return m.Param2
	}
	return 0
}

func (m *VirtualQueueData) GetParam3() int64 {
	if m != nil {
		return m.Param3
	}
	return 0
}

func (m *VirtualQueueData) GetPaths() []*Vector2D {
	if m != nil {
		return m.Paths
	}
	return nil
}

func (m *VirtualQueueData) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func init() {
	proto.RegisterType((*RadarData)(nil), "protomsg.RadarData")
	proto.RegisterType((*RadarVillageData)(nil), "protomsg.RadarVillageData")
	proto.RegisterType((*RadarStrongholdData)(nil), "protomsg.RadarStrongholdData")
	proto.RegisterType((*VirtualQueueData)(nil), "protomsg.VirtualQueueData")
	proto.RegisterEnum("protomsg.RadarStatus", RadarStatus_name, RadarStatus_value)
	proto.RegisterEnum("protomsg.RadarTimeType", RadarTimeType_name, RadarTimeType_value)
	proto.RegisterEnum("protomsg.RadarType", RadarType_name, RadarType_value)
	proto.RegisterEnum("protomsg.RadarMonsterSubType", RadarMonsterSubType_name, RadarMonsterSubType_value)
	proto.RegisterEnum("protomsg.RadarCommonSubType", RadarCommonSubType_name, RadarCommonSubType_value)
	proto.RegisterEnum("protomsg.RadarQuality", RadarQuality_name, RadarQuality_value)
	proto.RegisterEnum("protomsg.RadarDestroyReason", RadarDestroyReason_name, RadarDestroyReason_value)
	proto.RegisterEnum("protomsg.RadarCompleteType", RadarCompleteType_name, RadarCompleteType_value)
	proto.RegisterEnum("protomsg.VirtualQueueType", VirtualQueueType_name, VirtualQueueType_value)
	proto.RegisterEnum("protomsg.VirtualQueueStatus", VirtualQueueStatus_name, VirtualQueueStatus_value)
}

func init() { proto.RegisterFile("msg_common_radar.proto", fileDescriptor53) }

var fileDescriptor53 = []byte{
	// 894 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xe1, 0x6e, 0xe3, 0x44,
	0x10, 0xc6, 0x71, 0xda, 0xa4, 0x93, 0x86, 0xdb, 0xdb, 0xde, 0xb5, 0x6e, 0x7b, 0xbd, 0x0b, 0x05,
	0xa1, 0xc8, 0x3f, 0x2a, 0xd1, 0x3e, 0x41, 0x49, 0x11, 0x8a, 0x10, 0x94, 0x73, 0xab, 0x43, 0xfc,
	0x8a, 0xb6, 0xf6, 0xe2, 0x98, 0x38, 0x5e, 0xe3, 0x5d, 0xb7, 0x97, 0xc7, 0xe0, 0x15, 0xf8, 0xc3,
	0x03, 0xf0, 0x93, 0x37, 0xe2, 0x29, 0xd0, 0x8e, 0xbd, 0x8e, 0xdd, 0x38, 0x42, 0xe2, 0x57, 0x32,
	0xdf, 0x37, 0x9e, 0x19, 0x7f, 0xf3, 0x79, 0xe0, 0x70, 0x29, 0xc3, 0x99, 0x2f, 0x96, 0x4b, 0x91,
	0xcc, 0x32, 0x16, 0xb0, 0xec, 0x22, 0xcd, 0x84, 0x12, 0xb4, 0x8f, 0x3f, 0x4b, 0x19, 0x9e, 0x90,
	0x75, 0x46, 0xc1, 0x9d, 0xff, 0x69, 0xc3, 0x9e, 0xa7, 0x73, 0x6f, 0x98, 0x62, 0xf4, 0x18, 0xfa,
	0xf8, 0xe0, 0x2c, 0x0a, 0x1c, 0x6b, 0x64, 0x8d, 0x87, 0x5e, 0x0f, 0xe3, 0x69, 0x40, 0xbf, 0x84,
	0x17, 0x05, 0xe5, 0x8b, 0xe4, 0x97, 0x28, 0xd4, 0x19, 0x1d, 0xcc, 0x18, 0x22, 0x3c, 0x41, 0x74,
	0x1a, 0xd0, 0x43, 0xd8, 0x95, 0x8a, 0xa9, 0x5c, 0x3a, 0xf6, 0xc8, 0x1a, 0xef, 0x78, 0x65, 0x44,
	0x5f, 0xc1, 0x4e, 0xca, 0x32, 0xb6, 0x74, 0xba, 0x23, 0x6b, 0x6c, 0x7b, 0x45, 0xa0, 0xb3, 0xf1,
	0xcf, 0xa5, 0xb3, 0x83, 0x70, 0x19, 0xd1, 0x0b, 0xe8, 0xa7, 0x42, 0x46, 0x2a, 0x12, 0x89, 0xb3,
	0x3b, 0xb2, 0xc6, 0x83, 0x4b, 0x7a, 0x61, 0xde, 0xe2, 0xe2, 0x91, 0xfb, 0x4a, 0x64, 0x97, 0x81,
	0x57, 0xe5, 0xd0, 0x33, 0x80, 0x07, 0x1e, 0x46, 0xc9, 0x4c, 0x45, 0x4b, 0xee, 0xf4, 0xb0, 0xd6,
	0x1e, 0x22, 0xf7, 0xd1, 0x92, 0xeb, 0xf7, 0xe2, 0x49, 0x50, 0x90, 0x7d, 0x24, 0x7b, 0x3c, 0x09,
	0x90, 0x3a, 0x03, 0x90, 0x8a, 0x65, 0xaa, 0x20, 0xf7, 0x8a, 0x27, 0x11, 0x41, 0xfa, 0x14, 0xf6,
	0xa4, 0x12, 0x69, 0xc1, 0x02, 0xb2, 0x7d, 0x0d, 0x18, 0x32, 0xe3, 0x4f, 0x2c, 0x0b, 0xb4, 0x1a,
	0x03, 0x54, 0xa3, 0x5f, 0x00, 0xd3, 0x80, 0x52, 0xe8, 0x86, 0xd1, 0x23, 0x77, 0xf6, 0x47, 0xd6,
	0xb8, 0xef, 0xe1, 0x7f, 0xfa, 0x19, 0xec, 0x07, 0x11, 0x8b, 0x45, 0x38, 0x8b, 0x92, 0x80, 0x7f,
	0x74, 0x86, 0x28, 0xd1, 0xa0, 0xc0, 0xa6, 0x1a, 0xa2, 0x47, 0xd0, 0x93, 0x73, 0xf1, 0xa4, 0x2b,
	0x7e, 0x5a, 0x0a, 0x38, 0x17, 0x4f, 0xd3, 0xe0, 0xdc, 0x07, 0x82, 0x8b, 0xfa, 0x10, 0xc5, 0x31,
	0x0b, 0x39, 0xee, 0xeb, 0x1d, 0x0c, 0x1e, 0x59, 0x1c, 0x47, 0x2c, 0xe4, 0x66, 0x65, 0x5d, 0x0f,
	0x0c, 0x34, 0x0d, 0x1a, 0x3a, 0x76, 0xfe, 0x5b, 0xc7, 0xf3, 0x5f, 0xe1, 0x00, 0x9b, 0xdc, 0xa9,
	0x4c, 0x24, 0xe1, 0x5c, 0xc4, 0x01, 0xf6, 0xf9, 0x1c, 0x86, 0xb2, 0x42, 0xd6, 0x9d, 0xf6, 0xd7,
	0xe0, 0xff, 0xe8, 0xf5, 0x8f, 0x05, 0xe4, 0x43, 0x94, 0xa9, 0x9c, 0xc5, 0xef, 0x73, 0x9e, 0x73,
	0xe3, 0xc0, 0xdf, 0x74, 0x50, 0x73, 0x20, 0xc6, 0xd3, 0xe0, 0xd9, 0xa6, 0x3a, 0xcf, 0x37, 0x45,
	0xa1, 0xab, 0x56, 0x29, 0x2f, 0x6d, 0x87, 0xff, 0x2b, 0x7b, 0x7d, 0x55, 0xba, 0xae, 0x8c, 0xb6,
	0xda, 0xce, 0xe0, 0x57, 0x68, 0x3a, 0x83, 0x5f, 0xd1, 0xb1, 0x36, 0xaf, 0x9a, 0x4b, 0xa7, 0x37,
	0xb2, 0xb7, 0xbc, 0x57, 0x91, 0x50, 0xb3, 0x7f, 0xbf, 0x6e, 0x7f, 0x37, 0x83, 0x41, 0x29, 0x2c,
	0x7e, 0x0d, 0xaf, 0xe1, 0xe5, 0xa2, 0x16, 0xcf, 0x7e, 0x10, 0x09, 0x27, 0x9f, 0xd0, 0x23, 0x38,
	0x68, 0xc0, 0xd7, 0xbe, 0xcf, 0x53, 0x45, 0x2c, 0x7a, 0x0c, 0xaf, 0x1b, 0xc4, 0x44, 0x2c, 0xd3,
	0x98, 0x2b, 0x4e, 0x3a, 0xd4, 0x81, 0x57, 0x0d, 0xea, 0x86, 0xeb, 0xad, 0xac, 0x88, 0xed, 0xfe,
	0x6e, 0xc1, 0x10, 0x19, 0xad, 0xcf, 0xbd, 0xd6, 0xa3, 0xaa, 0x6f, 0x10, 0xd3, 0xb8, 0x2a, 0x52,
	0x11, 0x37, 0x2c, 0x8a, 0x57, 0xf5, 0xce, 0x15, 0xf3, 0x13, 0xe7, 0x8b, 0x78, 0x45, 0x3a, 0x2d,
	0xd4, 0xcf, 0x9c, 0x65, 0xf1, 0x8a, 0xd8, 0x2d, 0x8d, 0x6e, 0x13, 0x9f, 0x93, 0xae, 0x3b, 0x2f,
	0xcf, 0x0d, 0x8e, 0x73, 0x00, 0x2f, 0x16, 0x55, 0x64, 0x46, 0x39, 0x04, 0x5a, 0x03, 0xbf, 0x17,
	0x89, 0x54, 0x3c, 0x23, 0xd6, 0x5a, 0x32, 0xc4, 0x27, 0x78, 0xc4, 0x48, 0xa7, 0xd6, 0x09, 0xa7,
	0xc6, 0x2f, 0x29, 0xe7, 0xc4, 0x76, 0x65, 0x69, 0xe5, 0xb2, 0xc2, 0x5d, 0xfe, 0x80, 0x3d, 0xcf,
	0xe0, 0x78, 0xd1, 0x82, 0x9b, 0xee, 0x6f, 0xc0, 0x69, 0xa5, 0x6f, 0x9f, 0x12, 0x62, 0xd1, 0x77,
	0x70, 0xda, 0xca, 0x9a, 0x69, 0xdc, 0xbf, 0x2c, 0xa0, 0x5e, 0x71, 0x0f, 0x35, 0x62, 0x9a, 0x56,
	0x55, 0x1b, 0xb0, 0xe9, 0xf9, 0x16, 0x4e, 0xda, 0xd8, 0x6b, 0x29, 0x23, 0xa9, 0x97, 0xbf, 0x85,
	0xf7, 0xb8, 0xf4, 0x73, 0xed, 0x80, 0x2d, 0xfc, 0x5d, 0x9e, 0xa6, 0xb8, 0x8c, 0x6a, 0xea, 0x26,
	0xff, 0xcd, 0xc7, 0x34, 0x16, 0x99, 0x5e, 0xca, 0x1f, 0x16, 0xec, 0x63, 0xc2, 0xfb, 0x9c, 0xc5,
	0x91, 0x5a, 0xad, 0x77, 0x50, 0x02, 0x1b, 0xfe, 0x34, 0xf8, 0xb7, 0x19, 0xe7, 0x5a, 0x98, 0x8d,
	0x07, 0xbe, 0x8e, 0xf3, 0x86, 0x39, 0x0d, 0xfe, 0x63, 0x9e, 0xa5, 0x31, 0x27, 0xf6, 0x26, 0x73,
	0x9b, 0xb1, 0x24, 0xe4, 0xa4, 0xbb, 0x5e, 0xb4, 0x61, 0x3c, 0x1e, 0x90, 0x1d, 0xf7, 0x6f, 0x23,
	0x6d, 0x69, 0x70, 0x8f, 0x33, 0x29, 0x92, 0xb5, 0xb4, 0x0d, 0x78, 0x43, 0xda, 0x26, 0xeb, 0xe1,
	0x91, 0xae, 0x2f, 0xb4, 0xc9, 0x6b, 0xd7, 0x8a, 0x5c, 0x91, 0x0e, 0xfd, 0x02, 0x46, 0x6d, 0x09,
	0x13, 0x26, 0x55, 0xcc, 0x27, 0x73, 0x1c, 0xd9, 0xa6, 0x23, 0x78, 0xd3, 0x96, 0x75, 0xc3, 0x56,
	0x1e, 0x97, 0x5c, 0x91, 0xae, 0x7b, 0x0f, 0x2f, 0xcd, 0x0a, 0xf0, 0xc3, 0x45, 0x5b, 0x9c, 0xc2,
	0xd1, 0x62, 0x03, 0x35, 0xa3, 0x9f, 0xc3, 0xdb, 0x16, 0x12, 0x2f, 0xe5, 0xb5, 0xef, 0xa7, 0x5c,
	0x11, 0xcb, 0xfd, 0xae, 0x79, 0x41, 0xb1, 0xe8, 0x09, 0x1c, 0x2e, 0x9e, 0x83, 0xa6, 0xa6, 0x6e,
	0xb8, 0xc1, 0x61, 0x13, 0x62, 0xb9, 0x0a, 0x68, 0x9d, 0x2b, 0x2f, 0x95, 0xd6, 0x77, 0x13, 0xae,
	0xeb, 0xdb, 0xc2, 0x9a, 0x01, 0x51, 0x98, 0x16, 0x7e, 0x7d, 0xbe, 0x1e, 0x76, 0xf1, 0x94, 0x5e,
	0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x81, 0xc3, 0x8b, 0xa0, 0xbd, 0x08, 0x00, 0x00,
}
