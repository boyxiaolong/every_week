package module

import (
	"consoleclient/data"
	"public/command"
	"public/common"
	"public/message/error_code"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

// "public/command"

// GScoutTest 测试用例对象
var GScoutTest *ScoutTest

func init() {
	GScoutTest = &ScoutTest{}
	GScoutTest.InitCmds("scout")
	GScoutTest.RegCommand("scoutmap", scoutMap)
	//GScoutTest.RegCommand("exploremist", exploreMist)
	GScoutTest.RegCommand("explorecave", exploreCave)
}

// ScoutTest 测试用例
type ScoutTest struct {
	TestBase
}

// 侦查
func scoutMap() (bool, int32) {
	DoCreateBuilding(uint32(protomsg.ConstructionType_kConstructionType_ScoutCamp))
	request := &protomsg.MsgCL2GSScoutRequest{}
	request.Command = &protomsg.ScoutCommand{}
	request.ScoutIndex = 1
	request.Command.TargetId = 70001 // 神庙
	request.Command.TargetType = protomsg.ScoutCommandTarget_kScoutCommandTarget_DoScout
	request.RegionId = 0

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLScoutReply).(*protomsg.MsgGS2CLScoutReply)
	if readResponse.ErrorCode != 0 {
		return false, readResponse.ErrorCode
	}

	//scoutMoveback()

	return true, 0
}

// 斥候撤回
func scoutMoveback() (bool, int32) {
	request := &protomsg.MsgCL2GSScoutMoveBackRequest{}
	request.ScoutIndex = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLScoutMoveBackReply).(*protomsg.MsgGS2CLScoutMoveBackReply)
	if response.ErrorCode != 0 {
		return false, response.ErrorCode
	}

	waitResponse := GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLSyncEntitiesDataRemoveNotice, 30)
	if waitResponse == nil {
		return false, int32(error_code.ErrorCode_kEcGeneralTimeOut)
	}

	common.GStdout.Success("Scout move back!")
	return true, 0
}

// 探索迷雾
func exploreMist() (bool, int32) {
	DoCreateBuilding(uint32(protomsg.ConstructionType_kConstructionType_ScoutCamp))

	var scoutSpeedType uint64 = uint64(protomsg.PlayerExtendAttribute_kExtendBattleAttrScoutSpeed)
	command.GCommand.ExecuteCommand("pm addeffect " + strconv.FormatUint(scoutSpeedType, 10) + " 10000")

	var x int64 = data.GetCastlePosX()
	var y int64 = data.GetCastlePosY()

	x = (x/270)*270 + 1
	y = (y/270)*270 + 1

	request := &protomsg.MsgCL2GSMistExploreRequest{}
	request.Route = make([]*protomsg.Vector2D, 0)
	var i int64 = 0
	var j int64 = 0
	for i = 0; i < 9; i++ {
		for j = 0; j < 9; j++ {
			pos := &protomsg.Vector2D{}
			pos.X = (x + i*30) * 1000
			pos.Y = (y + j*30) * 1000

			request.Route = append(request.Route, pos)
		}
	}
	request.ScoutIndex = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMistExploreReply).(*protomsg.MsgGS2CLMistExploreReply)
	if response.ErrorCode != 0 {
		return false, response.ErrorCode
	}

	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMistOpenNotice, 30)
	scoutMoveback()
	command.GCommand.ExecuteCommand("pm removeeffect " + strconv.FormatUint(scoutSpeedType, 10) + " 10000")

	return true, 0
}

// 探索山洞
func exploreCave() (bool, int32) {
	var caveID uint64 = 50000
	var caveX int64 = 957
	var caveY int64 = 2893

	// 创建建筑
	DoCreateBuilding(uint32(protomsg.ConstructionType_kConstructionType_ScoutCamp))

	// 寻找位置，迁城
	searchRequest := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	searchRequest.EntityType = protomsg.EntityType_kEntityType_Castle
	searchRequest.SearchLevel = 1
	searchRequest.SearchRange = 60
	searchRequest.NearbyPlayerCastle = false
	searchRequest.CenterPos = &protomsg.Vector2D{X: caveX * 1000, Y: caveY * 1000}

	ok, pos := FindEmptyPos(searchRequest)
	if !ok {
		return false, int32(error_code.ErrorCode_kECMapNotFoundEmptyPos)
	}

	var castleX int64 = pos.X / 6000
	var castleY int64 = pos.Y / 6000

	command.GCommand.ExecuteCommand("pm movecastle 99 " + strconv.FormatInt(castleX, 10) + " " + strconv.FormatInt(castleY, 10))
	common.GStdout.Success("exploreCave move castle [%v, %v] success!", castleX, castleY)

	request := &protomsg.MsgCL2GSExploreVisitRequest{}
	request.ScoutIndex = 1
	request.ObjectId = caveID
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgCL2GSExploreVisitReply).(*protomsg.MsgCL2GSExploreVisitReply)
	if response.ErrorCode != 0 {
		return false, response.ErrorCode
	}
	common.GStdout.Success("exploreCave visit cave successs!")

	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLExploreUpdateCaveNotice, 30)
	scoutMoveback()

	return true, 0
}
