package data

import (
	//"public/common"
	"sync"
	"public/message/msgtype"
	"public/message/protomsg"
	"github.com/golang/protobuf/proto"
)

// BuildData 建筑数据
type BuildData struct {
	// 建筑动态ID
	ID uint32
	// 建筑类型ID, city_building对应ID
	ConstructionType uint32
	// 建筑等级
	Level uint32
}

type CityData struct {
	Buildings  *sync.Map
}

func init() {
	RegisterDataCreator(CreateCityData)
}

func CreateCityData(data_center *DataCenter) {
	data := &CityData{Buildings:&sync.Map{}}
	data_center.DataRegister(data)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLCityInfoReply, data.OnMsgGS2CLCityInfoReply)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLCityBuildingUpdateNotice, data.OnMsgGS2CLCityBuildingUpdateNotice)
}

func setBuildDataInfo(buildData *BuildData, building *protomsg.BuildingInfo) {
	buildData.ID = building.Id
	buildData.Level = building.Level
	buildData.ConstructionType = building.Type
}

// OnMsgGS2CLCityInfoReply  内城建筑基础信息
func (data *CityData)OnMsgGS2CLCityInfoReply(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLCityInfoReply)
	for _, building := range msg.Map.Buildings {
		buildData := &BuildData{}
		setBuildDataInfo(buildData, building)
		data.Buildings.Store(buildData.ID,buildData)
	}
}

// OnMsgGS2CLCityBuildingUpdateNotice  建筑更新通知
func (data *CityData)OnMsgGS2CLCityBuildingUpdateNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLCityBuildingUpdateNotice)
	if buildData, ok := data.Buildings.Load(msg.Data.Id); ok {
		setBuildDataInfo(buildData.(*BuildData), msg.Data)
	} else {
		buildData := &BuildData{}
		setBuildDataInfo(buildData, msg.Data)
		data.Buildings.Store(buildData.ID,buildData)
	}
}

func (data *CityData)GetBuildMaxLevelByType(constructionType uint32) uint32 {
	var maxLevel uint32
	data.Buildings.Range(func(k, v interface{}) bool {
		buildingData:= v.(*BuildData)
		if buildingData.ConstructionType == constructionType {
			if buildingData.Level > maxLevel {
				maxLevel = buildingData.Level
			}
		}
		return true
	})
	return maxLevel
}
