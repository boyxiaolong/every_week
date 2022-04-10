package module

import (
	"consoleclient/data"
	"public/message/msgtype"
	"public/message/protomsg"
)

// "public/command"

// GCarriageTest 测试用例对象
var GCarriageTest *CarriageTest

func init() {
	GCarriageTest = &CarriageTest{}
	GCarriageTest.InitCmds("carriage")
	//GCarriageTest.RegCommand("getcarriagecapacity", getCarriageCapacity)
	//GCarriageTest.RegCommand("createcarriage", createCarriage)
	//GCarriageTest.RegCommand("moveback", movebackCarriage)
	//GCarriageTest.RegCommand("createcarriage", createCarriage)
}

// CarriageTest 测试用例
type CarriageTest struct {
	TestBase
}

func getCarriageCapacity() (bool, int32) {
	playerID := getTargetPlayerID()

	request := &protomsg.MsgCL2GSGetCarriageCapacityRequest{}
	request.TargetPlayerId = playerID
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGetCarriageCapacityReply).(*protomsg.MsgGS2CLGetCarriageCapacityReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func buildTradingPost() {
	buildType := uint32(protomsg.ConstructionType_kConstructionType_TradingPost)
	count := data.CountBuildsByType(buildType)
	if count == 0 {
		DoCreateBuilding(buildType)
	}
}

func getTargetPlayerID() uint64 {
	return uint64(20010)
}

func getTargetEntityID() uint64 {
	playerID := getTargetPlayerID()

	buildTradingPost()

	castle := data.GetPlayerCastle(playerID)
	if castle == nil {
		return 0
	}

	return castle.EntityID
}

var gCarriageIndex uint32

func createCarriage() (bool, int32) {
	targetID := getTargetEntityID()

	request := &protomsg.MsgCL2GSCreateCarriageRequest{}
	request.Command = &protomsg.MarchCommand{}
	request.Command.Position = &protomsg.Vector2D{}
	request.Command.Position.X = 0
	request.Command.Position.Y = 0
	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_CarriageSend
	request.Command.TargetId = targetID
	request.Resources = &protomsg.ResourceSet{}
	request.Resources.Res = make([]*protomsg.Resource, 0, 10)
	request.Resources.Res = append(request.Resources.Res, &protomsg.Resource{ResType: 1, SubType: 3, Value: 2000})
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateCarriageReply).(*protomsg.MsgGS2CLCreateCarriageReply)
	if readResponse.ErrorCode == 0 {
		gCarriageIndex = readResponse.Data.MarchIndex
	}

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func movebackCarriage() (bool, int32) {
	marchIndex := gCarriageIndex

	request := &protomsg.MsgCL2GSMarchCommandRequest{}
	request.MarchIndex = marchIndex
	request.Command = &protomsg.MarchCommand{}

	request.Command.TargetType = protomsg.MarchCommandTarget_kMarchCommandTarget_ForceMoveBack
	request.Command.TargetId = 0
	request.RegionId = 0

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMarchCommandReply).(*protomsg.MsgGS2CLMarchCommandReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
