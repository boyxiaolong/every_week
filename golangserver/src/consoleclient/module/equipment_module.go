package module

import (
	"consoleclient/data"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

// "public/command"

// GEquipmentTest 测试用例对象
var GEquipmentTest *EquipmentTest

func init() {
	GEquipmentTest = &EquipmentTest{}
	GEquipmentTest.InitCmds("equipment")
	GEquipmentTest.RegCommand("wear", wearEquipment)
	GEquipmentTest.RegCommand("takeoff", takeOffEquipment)
	GEquipmentTest.RegCommand("uprank", upRankEquipment)
	GEquipmentTest.RegCommand("decompose", decomposeEquipment)
	GEquipmentTest.RegCommand("materialcomposite", materialComposite)
	GEquipmentTest.RegCommand("materialdecompose", materialDecompose)
	GEquipmentTest.RegCommand("exchange", equipmentExchange)
	GEquipmentTest.RegCommand("add", addequipment)
}

// EquipmentTest 测试用例
type EquipmentTest struct {
	TestBase
}

func getEquipment(equipmentID uint32) uint32 {
	count := data.CountEquipmentsByEquipmentID(equipmentID)
	if count == 0 {
		command.GCommand.ExecuteCommand("pm addequipment " + strconv.Itoa(int(equipmentID)))
	}
	equipments := data.GetEquipmentsByEquipmentID(equipmentID)
	if equipments == nil || len(equipments) == 0 {
		return 0
	}

	return equipments[0].ID
}

func wearEquipment() (bool, int32) {
	command.GCommand.ExecuteCommand("pm addhero 1001")
	request := &protomsg.MsgCL2GSWearEquipmentRequest{}
	request.HeroId = 1001
	request.Id = getEquipment(61111)

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLWearEquipmentReply).(*protomsg.MsgGS2CLWearEquipmentReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 获取英雄穿戴的（第）一个装备
func getHeroOneEquipment(heroID uint32) uint32 {
	if ids, ok := data.GeHeroEquipmentIDs(heroID); ok {
		return ids[0]
	}
	return 0
}

func takeOffEquipment() (bool, int32) {
	request := &protomsg.MsgCL2GSTakeOffEquipmentRequest{}
	request.HeroId = 1001
	request.Id = getHeroOneEquipment(1001)

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLTakeOffEquipmentReply).(*protomsg.MsgGS2CLTakeOffEquipmentReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func upRankEquipment() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 60000 30")
	command.GCommand.ExecuteCommand("pm additem 60010 3")
	command.GCommand.ExecuteCommand("pm additem 60020 30")
	//command.GCommand.ExecuteCommand("pm additem 60001 3")
	//command.GCommand.ExecuteCommand("pm additem 60011 3")
	//command.GCommand.ExecuteCommand("pm additem 60021 3")
	command.GCommand.ExecuteCommand("pm addcurrency 4 1500")

	request := &protomsg.MsgCL2GSUpRankEquipmentRequest{}
	request.Id = getEquipment(61411)

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLUpRankEquipmentReply).(*protomsg.MsgGS2CLUpRankEquipmentReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func decomposeEquipment() (bool, int32) {
	request := &protomsg.MsgCL2GSEquipmentDecomposeRequest{}
	request.Id = getEquipment(61411)

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLEquipmentDecomposeReply).(*protomsg.MsgGS2CLEquipmentDecomposeReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func materialComposite() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 60000 4")

	request := &protomsg.MsgCL2GSMaterialCompositeRequest{}
	request.ItemId = 60001
	request.Count = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMaterialCompositeReply).(*protomsg.MsgGS2CLMaterialCompositeReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func materialDecompose() (bool, int32) {
	request := &protomsg.MsgCL2GSMaterialDecomposeRequest{}
	request.ItemId = 60001
	request.Count = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMaterialDecomposeReply).(*protomsg.MsgGS2CLMaterialDecomposeReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func equipmentExchange() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 60501 100")

	request := &protomsg.MsgCL2GSEquipmentExchangeRequest{}
	request.EquipmentId = 61011

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLEquipmentExchangeReply).(*protomsg.MsgGS2CLEquipmentExchangeReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func addequipment() (bool, int32) {
	// for i := 0; i < 998; i++ {
	// 	command.GCommand.ExecuteCommand("pm addequipment 61411")
	// 	if i%100 == 0 {
	// 		time.Sleep(2000 * time.Millisecond)
	// 	}
	// }

	return true, 0
}
