package module

import (
	//"public/command"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GRecruitTest 测试用例对象
var GRecruitTest *RecruitTest

func init() {
	GRecruitTest = &RecruitTest{}
	GRecruitTest.InitCmds("recruit")
	GRecruitTest.RegCommand("beginrecruit", beginRecruit)
}

//RecruitTest 测试用例
type RecruitTest struct {
	TestBase
}

func beginRecruit() (bool, int32) {
	command.GCommand.ExecuteCommand("pm createbuilding 1012 1 1 1")

	command.GCommand.ExecuteCommand("pm additem 50000 1000")
	request := &protomsg.MsgCL2GSRecruitRequest{}
	request.Type = protomsg.RecruitType_kRecruitType_Normal
	request.IsAll = true

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLRecruitReply).(*protomsg.MsgGS2CLRecruitReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
