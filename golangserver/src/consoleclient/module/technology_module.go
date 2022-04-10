package module

import (
	"consoleclient/data"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
	"time"
)

// "public/command"

// GtechnologyTest 测试用例对象
var GtechnologyTest *technologyTest

func init() {
	GtechnologyTest = &technologyTest{}
	GtechnologyTest.InitCmds("tech")
	GtechnologyTest.RegCommand("reserch", research)
	GtechnologyTest.RegCommand("cancel", cancelReserch)
	GtechnologyTest.RegCommand("reserch", research)
	GtechnologyTest.RegCommand("deal", dealReserch)
	GtechnologyTest.RegCommand("instantreserch", instantResearch)
	GtechnologyTest.RegCommand("reserch", research)
	GtechnologyTest.RegCommand("vip", researchVipFree)
	GtechnologyTest.RegCommand("reserch", research)
	GtechnologyTest.RegCommand("instantspeedup", instantSpeedUpResearch)
}

// technologyTest 测试用例
type technologyTest struct {
	TestBase
}

func resetAllTech() {
	command.GCommand.ExecuteCommand("pm resetalltech")
}

func prepareScienceCenter() uint32 {

	var castleLevel uint32
	castleLevel = 4

	castles := data.GetBuildsByType(uint32(protomsg.ConstructionType_kConstructionType_MainCastle))
	if castles[0].Level < castleLevel {
		SetBuildingLevel(castles[0].ID, castleLevel)
	}

	buildType := uint32(protomsg.ConstructionType_kConstructionType_ScienceCenter)

	count := data.CountBuildsByType(buildType)
	if count == 0 {
		_, error := DoCreateBuilding(buildType)
		if error != 0 {
			return 0
		}
	}
	buildings := data.GetBuildsByType(buildType)
	if nil == buildings {
		return 0
	}

	return buildings[0].ID
}

// 研究
func research() (bool, int32) {
	prepareScienceCenter()
	resetAllTech()

	request := &protomsg.MsgCL2GSResearchRequest{}
	request.TechId = 5001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLResearchnReply).(*protomsg.MsgGS2CLResearchnReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 收割
func dealReserch() (bool, int32) {
	ok, error := research()
	if !ok {
		return ok, error
	}

	speedUpResearch()
	time.Sleep(time.Duration(1100) * time.Millisecond)

	request := &protomsg.MsgCL2GSDealResearchRequest{}

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLDealResearchReply).(*protomsg.MsgGS2CLDealResearchReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 切换
func switchReserch() (bool, int32) {
	ok, error := research()
	if !ok {
		return ok, error
	}

	request := &protomsg.MsgCL2GSSwitchResearchRequest{}
	request.OldTechId = 5
	request.NewTechId = 6

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSwitchResearchReply).(*protomsg.MsgGS2CLSwitchResearchReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 取消研究
func cancelReserch() (bool, int32) {
	ok, error := research()
	if !ok {
		return ok, error
	}

	request := &protomsg.MsgCL2GSCancelResearchRequest{}

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCancelResearchReply).(*protomsg.MsgGS2CLCancelResearchReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 立即研究
func instantResearch() (bool, int32) {

	prepareScienceCenter()
	resetAllTech()

	request := &protomsg.MsgCL2GSInstantResearchRequest{}
	request.TechId = 5001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLInstantResearchnReply).(*protomsg.MsgGS2CLInstantResearchnReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func researchVipFree() (bool, int32) {
	request := &protomsg.MsgCL2GSResearchVipFreeTimeRequest{}
	request.TechId = 5001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLResearchVipFreeTimeReply).(*protomsg.MsgGS2CLResearchVipFreeTimeReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func speedUpResearch() {
	command.GCommand.ExecuteCommand("pm instantspeedupreserch")
}

func instantSpeedUpResearch() (bool, int32) {
	request := &protomsg.MsgCL2GSInstantSpeedUpResearchRequest{}
	request.TechId = 5001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLInstantSpeedUpResearchReply).(*protomsg.MsgGS2CLInstantSpeedUpResearchReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
