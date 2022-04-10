package module

import (
	//"public/command"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GMileStoneTest 测试用例对象
var GMileStoneTest *MileStoneTest

func init() {
	GMileStoneTest = &MileStoneTest{}
	GMileStoneTest.InitCmds("milestone")
	GMileStoneTest.RegCommand("milestoneinfo", mileStoneInfo)
	GMileStoneTest.RegCommand("milestonereward", mileStoneReward)
	GMileStoneTest.RegCommand("milestonerank", milestoneRank)
}

//MileStoneTest 测试用例
type MileStoneTest struct {
	TestBase
}

func mileStoneInfo() (bool, int32) {
	request := &protomsg.MsgCL2GSMileStoneInfoRequest{}
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMileStoneInfoReply).(*protomsg.MsgGS2CLMileStoneInfoReply)
	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func mileStoneReward() (bool, int32) {
	command.GCommand.ExecuteCommand("pm opennextmilestone")
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMileStoneChangeNotice)
	mileStoneInfo()

	request := &protomsg.MsgCL2GSMileStoneGetRewardRequest{}
	request.MileStoneId = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMileStoneGetRewardReply).(*protomsg.MsgGS2CLMileStoneGetRewardReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func milestoneRank() (bool, int32) {
	request := &protomsg.MsgCL2GSMileStoneRankListRequest{}
	request.MileStoneId = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMileStoneRankListReply).(*protomsg.MsgGS2CLMileStoneRankListReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
