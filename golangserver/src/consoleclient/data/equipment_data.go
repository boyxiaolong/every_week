package data

import (
	"public/common"
	"public/message/protomsg"

	"github.com/golang/protobuf/proto"
)

// EquipmentData 装备数据
type EquipmentData struct {
	// 装备动态ID
	ID uint32
	// 装备ID
	EquipmentID uint32
	// 等级
	Level uint32
	// 英雄id
	HeroID uint32
}

// GEquipmentDatas 装备数据
var GEquipmentDatas map[uint32]*EquipmentData

// GHeroEquipmentDatas 英雄装备数据
var GHeroEquipmentDatas map[uint32][]uint32

func init() {
	GEquipmentDatas = make(map[uint32]*EquipmentData)
	GHeroEquipmentDatas = make(map[uint32][]uint32)
}

func setEquipmentDataInfo(result *EquipmentData, data *protomsg.EquipmentInfo) {
	result.ID = data.Id
	result.Level = data.Level
	result.EquipmentID = data.EquipmentId
	result.HeroID = data.HeroId
	if result.HeroID > 0 {
		var equipments []uint32
		if len(GHeroEquipmentDatas[result.HeroID]) == 0 {
			equipments = make([]uint32, 0, 10)
		} else {
			equipments = GHeroEquipmentDatas[result.HeroID]
		}
		equipments = append(equipments, result.ID)
		GHeroEquipmentDatas[result.HeroID] = equipments
	}
}

// OnMsgGS2CLAllEquipmentNotice  所有装备信息
func OnMsgGS2CLAllEquipmentNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLAllEquipmentNotice)
	for _, info := range msg.Equipments {
		data := &EquipmentData{}
		setEquipmentDataInfo(data, info)
		GEquipmentDatas[data.ID] = data
	}
}

// OnMsgGS2CLAddEquipmentsNotice  新增装备信息
func OnMsgGS2CLAddEquipmentsNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLAddEquipmentsNotice)
	for _, info := range msg.Equipments {
		data := &EquipmentData{}
		setEquipmentDataInfo(data, info)
		GEquipmentDatas[data.ID] = data
	}
}

// OnMsgGS2CLEquipmentNotice  更新装备通知
func OnMsgGS2CLEquipmentNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLEquipmentNotice)
	if data, ok := GEquipmentDatas[msg.Equipment.Id]; ok {
		if data.HeroID > 0 {
			if ids, ok := GHeroEquipmentDatas[data.HeroID]; ok {
				GHeroEquipmentDatas[data.HeroID] = deleteHeroEquipmentID(ids, data.ID)
			}
		}
		setEquipmentDataInfo(data, msg.Equipment)
	} else {
		data := &EquipmentData{}
		setEquipmentDataInfo(data, msg.Equipment)
		GEquipmentDatas[data.ID] = data
	}
}

// OnMsgGS2CLRemoveEquipmentsNotice  删除装备通知
func OnMsgGS2CLRemoveEquipmentsNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLRemoveEquipmentsNotice)
	for _, id := range msg.Ids {
		doRemoveEquipment(id)
	}
}

func doRemoveEquipment(id uint32) {
	if data, ok := GEquipmentDatas[id]; ok {
		if data.HeroID > 0 {
			if ids, ok := GHeroEquipmentDatas[data.HeroID]; ok {
				GHeroEquipmentDatas[data.HeroID] = deleteHeroEquipmentID(ids, data.ID)
			}
		}
		delete(GEquipmentDatas, id)
	}
}

func deleteHeroEquipmentID(src []uint32, id uint32) []uint32 {
	result := make([]uint32, 0, len(src))
	for _, val := range src {
		if val == id {
			result = append(result, val)
		}
	}
	return result
}

func doGetEquipmentsByEquipmentID(args ...interface{}) interface{} {
	datas := make([]*EquipmentData, 0, 10)

	equipmentID, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("doGetEquipmentsByEquipmentID error")
		return datas
	}

	for _, data := range GEquipmentDatas {
		if data.EquipmentID == equipmentID {
			datas = append(datas, data)
		}
	}
	return datas
}

// GetEquipmentsByEquipmentID 根据装备id获取装备列表
func GetEquipmentsByEquipmentID(equipmentID uint32) []*EquipmentData {
	data, ok := GDataCenter.GetDataParam(doGetEquipmentsByEquipmentID, equipmentID).([]*EquipmentData)
	if !ok {
		return nil
	}

	return data
}

func doCountEquipmentsByEquipmentID(args ...interface{}) interface{} {

	equipmentID, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("doCountEquipmentsByEquipmentID error")
		return 0
	}

	var count uint32
	count = 0
	for _, data := range GEquipmentDatas {
		if data.EquipmentID == equipmentID {
			count++
		}
	}
	return count
}

// CountEquipmentsByEquipmentID 根据装备id获取装备数量
func CountEquipmentsByEquipmentID(equipmentID uint32) uint32 {
	return GDataCenter.GetDataParam(doCountEquipmentsByEquipmentID, equipmentID).(uint32)
}

func doGeEquipment(args ...interface{}) interface{} {

	id, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("doGeEquipment error")
		return nil
	}

	if data, ok := GEquipmentDatas[id]; ok {
		return data
	}
	return nil
}

// GeEquipment 获取装备
func GeEquipment(id uint32) (*EquipmentData, bool) {
	if data, ok := GDataCenter.GetDataParam(doGeEquipment, id).(*EquipmentData); ok {
		return data, true
	}
	return nil, false
}

func doGeHeroEquipmentIDs(args ...interface{}) interface{} {
	heroID, ok := args[0].(uint32)
	if !ok {
		common.GStdout.Error("doGeHeroEquipmentIDs error")
	}

	if data, ok := GHeroEquipmentDatas[heroID]; ok {
		return data
	}
	return nil
}

// GeHeroEquipmentIDs 获取英雄穿戴的装备id列表
func GeHeroEquipmentIDs(heroID uint32) ([]uint32, bool) {
	if data, ok := GDataCenter.GetDataParam(doGeHeroEquipmentIDs, heroID).([]uint32); ok {
		return data, true
	}
	return nil, false
}
