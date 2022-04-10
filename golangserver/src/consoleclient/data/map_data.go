package data

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"public/common"
	"reflect"

	"public/message/protomsg"

	"github.com/golang/protobuf/proto"
)

//MapData struct
type MapData struct {
}

type CastleData struct {
	EntityID uint64
	Pos      protomsg.Vector2D
	PlayerID uint64
}

//GMapData xxxx
var GMapData *MapData

var GPlayerCastle map[uint64]*CastleData
var GCastle map[uint64]*CastleData

func init() {
	GMapData = &MapData{}
	GPlayerCastle = make(map[uint64]*CastleData)
	GCastle = make(map[uint64]*CastleData)
}

//GPropertyTypeName 属性名
var GPropertyTypeName = map[uint32]string{
	1:  "protomsg.MapData",
	2:  "protomsg.EntityOwnerData",
	4:  "protomsg.MoveData",
	5:  "protomsg.GuildBuildingData",
	6:  "protomsg.ArmyData",
	7:  "protomsg.PassData",
	8:  "protomsg.TempleData",
	9:  "protomsg.CollectData",
	10: "protomsg.DBReenforceData",
	11: "protomsg.BarbarianFortData",
	12: "protomsg.BarbariansData",
	14: "protomsg.CityWallData",
	15: "protomsg.EntityTypeData",
	16: "protomsg.MarchCommand",
	17: "protomsg.EntityMarchData",
	18: "protomsg.ResourceData",
	19: "protomsg.GuildMarchData",
	20: "protomsg.CastleExtendData",
	21: "protomsg.AttackTargetData",
	22: "protomsg.MoveSpeedData",
	23: "protomsg.MapAlarmSyncDatas",
	24: "protomsg.SkillStatusData",
	26: "protomsg.MultiStrikeData",
	27: "protomsg.BattleStatusData",
	28: "protomsg.ScoutCommand",
	29: "protomsg.EntityScoutData",
	30: "protomsg.SummonData",
	31: "protomsg.TrapData",
	32: "protomsg.GuildCollectBuildingData",
	33: "protomsg.GenericMonsterData",
	34: "protomsg.SummonMonsterData",
	35: "protomsg.ExpeStrongholdData",
	36: "protomsg.ReenforceData",
	37: "protomsg.ExpeEnemyFortData",
	38: "protomsg.MovePathData",
	39: "protomsg.PompeiiBuildData",
	43: "protomsg.EntityCarriageData",
	47: "protomsg.GuildMarchLimit",
	48: "protomsg.EntityKingdomTitle",
	49: "protomsg.ScoutExtendData",
	51: "protomsg.PompeiiLabData",
	52: "protomsg.PompeiiCultureTankData",
	53: "protomsg.MarchPickUpData",
}

func doGetPlayerCastleByPlayerID(args ...interface{}) interface{} {
	datas := make([]*EquipmentData, 0, 10)

	playerID, ok := args[0].(uint64)
	if !ok {
		common.GStdout.Error("doGetPlayerCastleByPlayerID error")
		return datas
	}
	if data, ok := doGetPlayerCastle(playerID); ok {
		return data
	}
	return nil
}

func GetPlayerCastle(playerID uint64) *CastleData {
	data, ok := GDataCenter.GetDataParam(doGetPlayerCastleByPlayerID, playerID).(*CastleData)
	if !ok {
		return nil
	}

	return data
}

func doSetPlayerCastle(data *CastleData) {
	if data.PlayerID > 0 {
		GPlayerCastle[data.PlayerID] = data
	}

	if data.EntityID > 0 {
		GCastle[data.EntityID] = data
	}
}

func doGetPlayerCastle(playerID uint64) (*CastleData, bool) {
	if data, ok := GPlayerCastle[playerID]; ok {
		return data, true
	}
	return nil, false
}

func doMutableCastle(entityID uint64) *CastleData {
	if data, ok := GCastle[entityID]; ok {
		return data
	}
	castle := &CastleData{}
	castle.EntityID = entityID
	doSetPlayerCastle(castle)
	return castle
}

// GetPropertyDataType content
func GetPropertyDataType(propertyType uint32) (bool, reflect.Type) {
	name, ok := GPropertyTypeName[propertyType]
	if !ok {
		common.GStdout.Error("propertyType not exist:%v", propertyType)
		return false, nil
	}

	t := proto.MessageType(name)
	if t == nil {
		common.GStdout.Error("propertyType message not exist:%v,name:%v", propertyType, name)
		return false, nil
	}
	return true, t
}

// MakeProperty content
func MakeProperty(propertyType uint32, data []byte) proto.Message {
	ok, t := GetPropertyDataType(propertyType)
	if !ok {
		return nil
	}

	if propertyMsg, ok := reflect.New(t.Elem()).Interface().(proto.Message); ok {
		err := proto.Unmarshal(data, propertyMsg)
		if err != nil {
			common.GStdout.Error("protobuf 解码出错:%v", err.Error())
			return nil
		}
		return propertyMsg
	}
	return nil
}

// PrintPropertyDetail content
func PrintPropertyDetail(entityID uint64, entityType uint32, propertyType uint32, propertyMsg proto.Message) {
	name, ok := GPropertyTypeName[propertyType]
	if !ok {
		common.GStdout.Error("propertyType not exist:%v", propertyType)
		return
	}

	data, err := json.Marshal(propertyMsg)
	if err != nil {
		common.GStdout.Error("JSON marshaling failed: %s", err)
		return
	}
	common.GStdout.Info("receive property success:property_name:%v,property_type:%v,data:%s", name, propertyType, data)
	if propertyType == 20 {
		extenddata := propertyMsg.(*protomsg.CastleExtendData)
		if extenddata.CampId == 0 {
			common.GStdout.Info("===============================CastleExtendData error================================================")
		}

	}
	if entityType == uint32(protomsg.EntityType_kEntityType_Castle) {
		if propertyType == uint32(protomsg.EntityPropertyType_kEntityPropertyType_Owner) {
			ownerData := propertyMsg.(*protomsg.EntityOwnerData)

			castle := doMutableCastle(entityID)
			castle.PlayerID = ownerData.PlayerId
			doSetPlayerCastle(castle)
		}
		if propertyType == uint32(protomsg.EntityPropertyType_kEntityPropertyType_Map) {
			mapData := propertyMsg.(*protomsg.MapData)

			castle := doMutableCastle(entityID)
			castle.Pos.X = mapData.Position.X
			castle.Pos.Y = mapData.Position.Y
		}
	}

}

// OnMsgGS2CLUpdateSyncEntitieNotice  公会基础信息
func OnMsgGS2CLUpdateSyncEntitieNotice(msginterface proto.Message) {
	common.GStdout.Info("========================================= recv UpdateSyncEntitieNotice start\n")
	msg := msginterface.(*protomsg.MsgGS2CLSyncEntitiesDataNotice)

	dataType := msg.ViewLod
	//common.GStdout.Info("%v\n", msg.SyncData)
	buf := bytes.NewReader(msg.SyncData)
	var entityCount uint32
	err := binary.Read(buf, binary.LittleEndian, &entityCount)
	if err != nil {
		common.GStdout.Info("%v\n", err)
		return
	}

	common.GStdout.Info(" region id %v, type %v , entity count %v \n", msg.RegionId, dataType, entityCount)

	for i := 0; i < int(entityCount); i++ {
		var entityID uint64
		var entitytype uint32
		var version uint32
		err = binary.Read(buf, binary.LittleEndian, &entitytype)
		if err != nil {
			common.GStdout.Info("%v\n", err)
			return
		}

		err = binary.Read(buf, binary.LittleEndian, &entityID)
		if err != nil {
			common.GStdout.Info("%v\n", err)
			return
		}

		err = binary.Read(buf, binary.LittleEndian, &version)
		if err != nil {
			common.GStdout.Info("%v\n", err)
			return
		}

		var propertyCount uint32
		err = binary.Read(buf, binary.LittleEndian, &propertyCount)
		if err != nil {
			common.GStdout.Info("%v\n", err)
			return
		}

		common.GStdout.Info("entity id %v , entity type %v , property count %v\n", entityID, entitytype, propertyCount)

		for j := 0; j < int(propertyCount); j++ {
			var propertyType uint32
			var propertyLength uint32
			err = binary.Read(buf, binary.LittleEndian, &propertyType)
			if err != nil {
				common.GStdout.Info("%v\n", err)
				return
			}

			err = binary.Read(buf, binary.LittleEndian, &propertyLength)
			if err != nil {
				common.GStdout.Info("%v\n", err)
				return
			}

			propertyData := make([]byte, propertyLength)

			err = binary.Read(buf, binary.LittleEndian, propertyData)
			if err != nil {
				common.GStdout.Info("%v\n", err)
				return
			}
			//common.GStdout.Info("property type %v , property length %v , data[%v]\n", propertyType, propertyLength, propertyData)
			propertyMsg := MakeProperty(propertyType, propertyData)
			if propertyMsg != nil {
				PrintPropertyDetail(entityID, entitytype, propertyType, propertyMsg)
			}
		}
	}
	common.GStdout.Info("========================================= recv UpdateSyncEntitieNotice end\n")
}
