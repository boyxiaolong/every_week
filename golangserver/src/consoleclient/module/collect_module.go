package module

import (
	//"public/command"
	"public/command"
	"public/message/error_code"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

//GCollectTest 测试用例对象
var GCollectTest *CollectTest

func init() {
	GCollectTest = &CollectTest{}
	GCollectTest.InitCmds("collect")
	GCollectTest.RegCommand("begincollect", beginCollect)
}

//CollectTest 测试用例
type CollectTest struct {
	TestBase
}

func beginCollect() (bool, int32) {
	command.GCommand.ExecuteCommand("pm addhero 1001")
	command.GCommand.ExecuteCommand("pm addarmy 101 10000")

	var troopSpeedType uint64 = uint64(protomsg.PlayerExtendAttribute_kExtendBattleAttrTroopSpeed_All)
	command.GCommand.ExecuteCommand("pm addeffect " + strconv.FormatUint(troopSpeedType, 10) + " 10000")

	var collectSpeedType uint64 = uint64(protomsg.PlayerExtendAttribute_kExtendAttrkIncRate_CollectFood)
	command.GCommand.ExecuteCommand("pm addeffect " + strconv.FormatUint(collectSpeedType, 10) + " 100000")

	command.GCommand.ExecuteCommand("pm mapsearch 1 1 1")
	response := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMapSearchReply).(*protomsg.MsgGS2CLMapSearchReply)

	if response.ErrorCode != 0 {
		return response.ErrorCode == 0, response.ErrorCode
	}

	request := &protomsg.MsgCL2GSCreateMarchRequest{}
	request.Command = &protomsg.MarchCommand{}
	request.Command.Position = &protomsg.Vector2D{}
	request.Command.Position.X = 0
	request.Command.Position.Y = 0
	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_Collect
	request.Command.TargetId = response.EntityId

	request.Army = &protomsg.ArmyData{}
	request.Army.Hero1 = 1001
	request.Army.Hero2 = 0
	request.Army.CurrentTroops = make([]*protomsg.TroopData, 0)

	newtroop := &protomsg.TroopData{}
	newtroop.TroopId = 101
	newtroop.TroopUints = 1000
	request.Army.CurrentTroops = append(request.Army.InitTroops, newtroop)

	marchResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateMarchReply).(*protomsg.MsgGS2CLCreateMarchReply)
	if marchResponse.ErrorCode != 0 && marchResponse.ErrorCode != int32(error_code.ErrorCode_kECMapStateNotConnected) {
		return marchResponse.ErrorCode == 0, marchResponse.ErrorCode
	}
	MoveBackMarch(marchResponse.March.MarchIndex)

	command.GCommand.ExecuteCommand("pm removeeffect " + strconv.FormatUint(troopSpeedType, 10) + " 10000")
	command.GCommand.ExecuteCommand("pm removeeffect " + strconv.FormatUint(collectSpeedType, 10) + " 100000")

	return true, 0
}
