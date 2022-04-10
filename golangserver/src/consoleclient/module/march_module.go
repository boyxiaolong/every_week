package module

import (
	"public/command"
	"public/common"
	"public/message/error_code"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GmarchTest 测试用例对象
var GmarchTest *marchTest

func init() {
	GmarchTest = &marchTest{}
	GmarchTest.InitCmds("march")
	GmarchTest.RegCommand("marchcommand", marchcommand)
	GmarchTest.RegCommand("moveback", moveback)
	GmarchTest.RegCommand("createmarch", createmarch)
	GmarchTest.RegCommand("marchpos", marchPosition)
	GmarchTest.RegCommand("movecastle", marchMoveCastle)
	GmarchTest.RegCommand("createguildmarch", createGuildMarch)
	GmarchTest.RegCommand("pathfind", pathFind)
}

//BaseTest 测试用例
type marchTest struct {
	TestBase
}

// CreateMarch 派兵驻扎
func CreateMarch() (uint32, int32) {

	searchRequest := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	searchRequest.EntityType = protomsg.EntityType_kEntityType_March
	searchRequest.SearchLevel = 1
	searchRequest.SearchRange = 30
	searchRequest.NearbyPlayerCastle = true

	// 查找空位
	ok, pos := FindEmptyPos(searchRequest)
	if !ok {
		return 0, int32(error_code.ErrorCode_kECMapNotFoundEmptyPos)
	}

	request := &protomsg.MsgCL2GSCreateMarchRequest{}
	request.Command = &protomsg.MarchCommand{}
	request.Command.Position = pos
	request.Command.Position.X = 0
	request.Command.Position.Y = 0
	request.Command.TargetType = 9
	request.Command.TargetId = 0
	//request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_Collect
	//request.Command.TargetId = 177327

	request.Army = &protomsg.ArmyData{}
	request.Army.Hero1 = 1043
	request.Army.Hero2 = 0
	request.Army.CurrentTroops = make([]*protomsg.TroopData, 0)

	newtroop := &protomsg.TroopData{}
	newtroop.TroopId = 141
	newtroop.TroopUints = 100
	request.Army.CurrentTroops = append(request.Army.InitTroops, newtroop)

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateMarchReply).(*protomsg.MsgGS2CLCreateMarchReply)

	if response.ErrorCode != 0 {
		return 0, response.ErrorCode
	}

	common.GStdout.Success("Create station march success!")
	return response.March.MarchIndex, response.ErrorCode
}

var g_marchIndex uint32 = 0
var g_target uint64 = 0

func marchcommand() (bool, int32) {

	g_marchIndex, error := CreateMarch()
	if 0 == g_marchIndex {
		return false, error
	}

	request := &protomsg.MsgCL2GSMarchCommandRequest{}
	request.MarchIndex = g_marchIndex
	request.Command = &protomsg.MarchCommand{}

	//request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_PickUp
	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_JoinReenforce
	if g_target != 70016 {
		g_target = 70016
	} else {
		g_target = 70017
	}
	request.Command.TargetId = g_target
	request.RegionId = 0

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMarchCommandReply).(*protomsg.MsgGS2CLMarchCommandReply)
	//MoveBackMarch(marchIndex)

	return response.ErrorCode == 0, response.ErrorCode
}

func moveback() (bool, int32) {
	//MoveBackMarch(1)
	MoveBackMarch(g_marchIndex)
	g_marchIndex = 0
	return true, 0
}

func createmarch() (bool, int32) {
	marchIndex, error := CreateMarch()
	if 0 == marchIndex {
		return false, error
	}
	MoveBackMarch(marchIndex)

	return true, 0
}

func marchPosition() (bool, int32) {

	searchRequest := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	searchRequest.EntityType = protomsg.EntityType_kEntityType_March
	searchRequest.SearchLevel = 1
	searchRequest.SearchRange = 50
	searchRequest.NearbyPlayerCastle = true

	// 查找空位
	ok, pos := FindEmptyPos(searchRequest)
	if !ok {
		return false, int32(error_code.ErrorCode_kECMapNotFoundEmptyPos)
	}

	marchIndex, error := CreateMarch()
	if 0 == marchIndex {
		return false, error
	}

	request := &protomsg.MsgCL2GSMarchCommandRequest{}
	request.MarchIndex = marchIndex
	request.Command = &protomsg.MarchCommand{}
	request.Command.Position = pos
	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_Position
	request.Command.TargetId = 0
	request.RegionId = 0

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMarchCommandReply).(*protomsg.MsgGS2CLMarchCommandReply)
	MoveBackMarch(marchIndex)

	return response.ErrorCode == 0, response.ErrorCode
}

// DisbandGuildMarch 解散集结部队
func DisbandGuildMarch(entityID uint64) {
	if 0 == entityID {
		return
	}

	request := &protomsg.MsgCL2GSDisbandGuildMarchRequest{}
	request.GuildMarchEntityId = entityID

	GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLDisbandGuildMarchReply)
}

// MoveBackMarch 队伍返回
func MoveBackMarch(marchIndex uint32) {
	request := &protomsg.MsgCL2GSMarchCommandRequest{}
	request.MarchIndex = marchIndex
	request.Command = &protomsg.MarchCommand{}
	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_MoveBack
	request.Command.TargetId = 0
	request.RegionId = 0

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMarchCommandReply).(*protomsg.MsgGS2CLMarchCommandReply)
	if response.ErrorCode != 0 {
		return
	}

	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMarchRemoveNotice, 30)
}

func createGuildMarch() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("x")
	}

	entityID, pos := SearchEntity(protomsg.MapSearchType_kMapSearchType_BarbarianFort, 1)
	if entityID == 0 {
		return false, int32(error_code.ErrorCode_kECMapNotFoundEntity)
	}

	request := &protomsg.MsgCL2GSCreateGuildMarchRequest{}
	request.WaitTimeIndex = 1
	request.Command = &protomsg.MarchCommand{}
	request.Command.Position = pos
	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_Battle
	request.Command.TargetId = entityID

	request.ArmyData = &protomsg.ArmyData{}
	request.ArmyData.Hero1 = 1001
	request.ArmyData.Hero2 = 0
	request.ArmyData.CurrentTroops = make([]*protomsg.TroopData, 0)

	newtroop := &protomsg.TroopData{}
	newtroop.TroopId = 101
	newtroop.TroopUints = 100
	request.ArmyData.CurrentTroops = append(request.ArmyData.InitTroops, newtroop)

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateGuildMarchReply).(*protomsg.MsgGS2CLCreateGuildMarchReply)
	if response.ErrorCode != 0 {
		return false, response.ErrorCode
	}

	DisbandGuildMarch(response.GuildMarchEntityId)
	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMarchRemoveNotice, 30)

	return true, 0
}

func marchMoveCastle() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 9003 10")

	searchRequest := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	searchRequest.EntityType = protomsg.EntityType_kEntityType_Castle
	searchRequest.SearchLevel = 1
	searchRequest.SearchRange = 30
	searchRequest.NearbyPlayerCastle = false
	searchRequest.CenterPos = &protomsg.Vector2D{6089000, 5175000}

	// 查找空位
	ok, pos := FindEmptyPos(searchRequest)
	if !ok {
		return false, int32(error_code.ErrorCode_kECMapNotFoundEmptyPos)
	}

	request := &protomsg.MsgCL2GSMoveCastleRequest{}
	request.Type = 4
	request.Position = pos

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMoveCastleReply).(*protomsg.MsgGS2CLMoveCastleReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func pathFind() (bool, int32) {
	request := &protomsg.MsgCL2GSPathFindRequest{}
	request.RegionId = 281479271743488

	request.Start = &protomsg.Vector2D{}
	request.Start.X = 6089000
	request.Start.Y = 5175000

	request.End = &protomsg.Vector2D{}
	request.End.X = 6039000
	request.End.Y = 2607000

	request.IsCrossPass = true

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPathFindReply).(*protomsg.MsgGS2CLPathFindReply)
	return response.ErrorCode == 0, response.ErrorCode
}
