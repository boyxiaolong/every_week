package data

import (
	"public/message/msgtype"
	"public/message/protomsg"
	"github.com/golang/protobuf/proto"
	"public/common"
)

var ValidHeros = [...]uint32{1001,1002,1003,1004,1005,1007,1008,1011,1012,1013,1014,1015,1016,1018,1019,1020,1021,1022,1023,1024,1025,1031,1032,1033,1034,1036,1039,1040,1043,1061,1062,1070}
//var ValidTroops = [...]uint32{141,101,102,103,151,152,153,104,121,122,127,117,129,130,135,143,105,106,107,124,125,128,108,132,133,136,118,156,157,158,109,119,110,111,155,126,112,123,154,142,131,134,113,114,144,115,116,120}
//var ValidTroops = [...]uint32{141,143,109,113} //  一级兵
//var ValidTroops = [...]uint32{102,106,110,144} // 三级兵
//var ValidTroops = [...]uint32{113} //  一级兵
var TroopCapacity = [...]uint32{2500,3500,4500,6000,8500,11000,14000,17500,22000,27000,32000,37000,42000,48000,54000,60000,67000,74000,82000,90000,100000,110000,120000,130000,150000}
var ValidCaves = [...]uint32{50000, 50001, 50002, 50003, 50004, 50005, 50006}
var HeroSkillGroup = map[uint32]uint32{ 1001:10012,1002:10022,1003:10032,1004:10042,1005:10052,1007:10072,1008:10082,
																			1011:10112,1012:10122,1013:10132,1014:10142,1015:10152,1016:10162,1018:10182,
																			1019:10192,1020:10202,1021:10212,1022:10222,1023:10232,1024:10242,1025:10252,
																			1031:10312,1032:10322,1033:10332,1034:10342,1036:10362,1039:10392,1040:10402,
																			1043:10432,1061:10612,1062:10622,1070:10702}

var ValidTroops = [][]uint32{[]uint32{141,143,109,113},[]uint32{101,105,119,114},
														[]uint32{102,106,110,144},
														[]uint32{103,151,152,153,107,124,125,128,111,155,126,115},
														[]uint32{104,121,122,127,108,132,133,136,112,123,154,116},
														[]uint32{117,129,130,135,118,156,157,158,142,131,134,120},
														}											
var MaxScoutCount = 3

var InstancesArmy = map[uint32][2]uint32{1:{1,14000},
																			2:{1,14000},
																			3:{1,14000},
																			4:{1,14000},
																			5:{2,14000},
																			6:{2,90000},
																			7:{2,14000},
																			8:{2,14000},
																			9:{2,14000},
																			10:{2,14000},
																			11:{2,17000},
																			12:{2,17000},
																			13:{2,17000},
																			14:{2,17000},
																			15:{2,17000},
																			16:{2,17000},
																			17:{2,17000},
																			18:{2,17000},
																			19:{2,17000},
																			20:{2,17000},
																			21:{2,43000},
																			22:{2,43000},
																			23:{2,43000},
																			24:{2,43000},
																			25:{2,43000},
																			26:{2,43000},
																			27:{2,43000},
																			28:{2,43000},
																			29:{2,43000},
																			30:{2,43000},
																			31:{2,43000},
																			32:{2,43000},
																			33:{2,43000},
																			34:{2,43000},
																			35:{2,43000},
																			36:{2,43000},
																			37:{2,64500},
																			38:{2,64500},
																			39:{2,64500},
																			40:{2,64500},
																			41:{3,141000},
																			42:{3,141000},
																			43:{3,141000},
																			44:{3,141000},
																			45:{3,141000},
																			46:{3,141000},
																			47:{3,141000},
																			48:{3,141000},
																			49:{3,141000},
																			50:{3,141000},
																			51:{3,141000},
																			52:{3,141000},
																			53:{3,141000},
																			54:{3,141000},
																			55:{3,141000},
																			56:{3,141000},
																			57:{3,141000},
																			58:{3,141000},
																			59:{3,141000},
																			60:{3,141000},
																			61:{4,291000},
																			62:{4,291000},
																			63:{4,291000},
																			64:{4,291000},
																			65:{4,291000},
																			66:{4,291000},
																			67:{4,291000},
																			68:{4,291000},
																			69:{4,291000},
																			70:{4,291000},
																			71:{4,291000},
																			72:{4,291000},
																			73:{4,291000},
																			74:{4,291000},
																			75:{4,291000},
																			76:{4,291000},
																			77:{4,291000},
																			78:{4,291000},
																			79:{4,291000},
																			80:{4,291000},
																			81:{5,480000},
																			82:{5,480000},
																			83:{5,480000},
																			84:{5,480000},
																			85:{5,480000},
																			86:{5,480000},
																			87:{5,480000},
																			88:{5,480000},
																			89:{5,480000},
																			90:{5,480000},
																			91:{5,480000},
																			92:{5,480000},
																			93:{5,480000},
																			94:{5,480000},
																			95:{5,480000},
																			96:{5,480000},
																			97:{5,480000},
																			98:{5,480000},
																			99:{5,480000},
																			100:{5,480000} }

func init() {
	RegisterDataCreator(CreateMapData)
}

func CreateMapData(data_center *DataCenter) {
	data := CreateRegionDataCenter()
	data_center.DataRegister(data)
	// data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLAddPlayerEntityNotice, data.OnMsgGS2CLMapBaseInfoReply)
	// data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLPlayerEntityUpdateNotice, data.OnMsgGS2CLPlayerEntityUpdateNotice)
	// data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLRemovePlayerEntityNotice, data.OnMsgGS2CLRemovePlayerEntityNotice)

	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLMapDataReply, data.OnMsgGS2CLMapDataReply)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLSyncEntitiesDataRemoveNotice, data.OnMsgGS2CLSyncEntitiesDataRemoveNotice)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLSyncEntitiesDataNotice, data.OnMsgGS2CLSyncEntitiesDataNotice)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLViewMapReply, data.OnMsgGS2CLViewMapReply)
	data_center.RegDispatch(msgtype.MsgType_kMsgG2CLIntoRegionNotice, data.OnMsgG2CLIntoRegionNotice)
	data_center.RegDispatch(msgtype.MsgType_kMsgG2CLLeaveRegionNotice, data.OnMsgG2CLLeaveRegionNotice)
}

// CheckIsMainRegion a
func CheckIsMainRegion(regionid uint64) bool {
	regiontype := regionid >> 48
	if regiontype == (uint64)(protomsg.RegionType_kRegionType_Kingdom) || regiontype == (uint64)(protomsg.RegionType_kRegionType_KVK) {
		return true
	}
	return false
}

// OnMsgGS2CLMapDataReply  通知地图数据
func (data *RegionDataCenter) OnMsgGS2CLMapDataReply(msg proto.Message) {
	dataMsg := msg.(*protomsg.MsgGS2CLMapDataReply)
	if !CheckIsMainRegion(dataMsg.RegionId) {
		return
	}
	data.MainRegionId = dataMsg.RegionId
}

// OnMsgGS2CLViewMapReply  查看地图数据
func (data *RegionDataCenter) OnMsgGS2CLViewMapReply(msg proto.Message) {
	dataMsg := msg.(*protomsg.MsgGS2CLViewMapReply)
	if dataMsg.ErrorCode == 0{
		return
	}

	if data.ViewRegionId != dataMsg.RegionId {
		if oldViewRegion,ok := data.GetRegion(data.ViewRegionId);ok{
			if viewCenter := oldViewRegion.GetEntityCenter(1); viewCenter != nil{
				viewCenter.Clear()
			}

			if viewCenter := oldViewRegion.GetEntityCenter(2); viewCenter != nil{
				viewCenter.Clear()
			}
		}
	}
	data.ViewRegionId = dataMsg.RegionId
}

func (c *RegionDataCenter) OnMsgGS2CLSyncEntitiesDataNotice(msg proto.Message) {
	dataMsg := msg.(*protomsg.MsgGS2CLSyncEntitiesDataNotice)
	regionData := c.GetRegionCenter(dataMsg.RegionId,dataMsg.ViewLod)
	if regionData != nil {
		regionData.UpdateSyncEntitiesData(dataMsg.SyncData)
	}else{
		common.GStdout.Info("OnMsgGS2CLSyncEntitiesDataNotice ========================================  %v", dataMsg.RegionId)
	}
}

func (c *RegionDataCenter) OnMsgGS2CLSyncEntitiesDataRemoveNotice(msg proto.Message) {
	dataMsg := msg.(*protomsg.MsgGS2CLSyncEntitiesDataRemoveNotice)
	regionData:= c.GetRegionCenter(dataMsg.RegionId,dataMsg.ViewLod)
	if regionData != nil {
		regionData.RemoveEntitiesData(dataMsg.EntityIds)
	}
}

func (c *RegionDataCenter) OnMsgG2CLIntoRegionNotice(msg proto.Message) {
	dataMsg := msg.(*protomsg.MsgG2CLIntoRegionNotice)
	c.CreateRegion(dataMsg.RegionId)
	common.GStdout.Info("OnMsgG2CLIntoRegionNotice ========================================  %v", dataMsg.RegionId)
}

func (c *RegionDataCenter) OnMsgG2CLLeaveRegionNotice(msg proto.Message) {
	dataMsg := msg.(*protomsg.MsgG2CLLeaveRegionNotice)
	c.DeleteRegion(dataMsg.RegionId)
}

