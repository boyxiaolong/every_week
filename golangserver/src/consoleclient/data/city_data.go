package data

import (
	"public/common"
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

// GBuildDatas 内城建筑数据
var GBuildDatas map[uint32]*BuildData

func init() {
	GBuildDatas = make(map[uint32]*BuildData)
}

func setBuildDataInfo(buildData *BuildData, building *protomsg.BuildingInfo) {
	buildData.ID = building.Id
	buildData.Level = building.Level
	buildData.ConstructionType = building.Type
}

// OnMsgGS2CLCityInfoReply  内城建筑基础信息
func OnMsgGS2CLCityInfoReply(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLCityInfoReply)
	for _, building := range msg.Map.Buildings {
		buildData := &BuildData{}
		setBuildDataInfo(buildData, building)
		GBuildDatas[buildData.ID] = buildData
	}
}

// OnMsgGS2CLCityBuildingUpdateNotice  建筑更新通知
func OnMsgGS2CLCityBuildingUpdateNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLCityBuildingUpdateNotice)
	if buildData, ok := GBuildDatas[msg.Data.Id]; ok {
		setBuildDataInfo(buildData, msg.Data)
	} else {
		buildData := &BuildData{}
		setBuildDataInfo(buildData, msg.Data)
		GBuildDatas[buildData.ID] = buildData
	}
}

func doGetBuildsByType(args ...interface{}) interface{} {
	datas := make([]*BuildData, 0, 10)

	constructionType, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("DoCountBuildsByType error")
		return datas
	}

	for _, buildData := range GBuildDatas {
		if buildData.ConstructionType == constructionType {
			datas = append(datas, buildData)
		}
	}
	return datas
}

// GetBuildsByType 根据类型获取建筑列表
func GetBuildsByType(constructionType uint32) []*BuildData {
	data, ok := GDataCenter.GetDataParam(doGetBuildsByType, constructionType).([]*BuildData)
	if !ok {
		return nil
	}

	return data
}

func doCountBuildsByType(args ...interface{}) interface{} {

	constructionType, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("DoCountBuildsByType error")
		return 0
	}

	var count uint32
	count = 0
	for _, buildData := range GBuildDatas {
		if buildData.ConstructionType == constructionType {
			count++
		}
	}
	return count
}

// CountBuildsByType 根据类型获取建筑数量
func CountBuildsByType(constructionType uint32) uint32 {
	return GDataCenter.GetDataParam(doCountBuildsByType, constructionType).(uint32)
}

func doGetBuild(args ...interface{}) interface{} {

	id, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("DoGetBuild error")
		return nil
	}

	if buildData, ok := GBuildDatas[id]; ok {
		return buildData
	}
	return nil
}

// GetBuild 获取建筑
func GetBuild(id uint32) (*BuildData, bool) {
	if building, ok := GDataCenter.GetDataParam(doGetBuild, id).(*BuildData); ok {
		return building, true
	}
	return nil, false
}
