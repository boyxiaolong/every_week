package module

import (
	"consoleclient/data"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

const ()

//GTaskTest 测试用例对象
var GTaskTest *TaskTest

func init() {

	GTaskTest = &TaskTest{}
	GTaskTest.InitCmds("task")
	GTaskTest.RegCommand("awardMainTask", awardMainTask)
	GTaskTest.RegCommand("awardBranchTask", awardBranchTask)
	GTaskTest.RegCommand("awardDailyTask", awardDailyTask)
}

//TaskTest 测试用例
type TaskTest struct {
	TestBase
}

func awardMainTask() (bool, int32) {
	command.GCommand.ExecuteCommand("pm openchapter 1")
	command.GCommand.ExecuteCommand("pm setmainstatus 1 2")
	request := &protomsg.MsgCL2GSTaskMainTaskAwardRequest{}
	request.TaskId = 1
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLTaskMainTaskAwardReply).(*protomsg.MsgGS2CLTaskMainTaskAwardReply)
	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}

func awardBranchTask() (bool, int32) {
	command.GCommand.ExecuteCommand("pm acceptbranch 1")
	command.GCommand.ExecuteCommand("pm setbranchstatus 1 2")
	request := &protomsg.MsgCL2GSTaskBranchTaskAwardRequest{}
	request.TaskId = 1
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLTaskBranchTaskAwardReply).(*protomsg.MsgGS2CLTaskBranchTaskAwardReply)
	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}

func awardDailyTask() (bool, int32) {
	dailyTaskID := data.GetFirstDaiylTaskID()
	strDailyTaskID := strconv.FormatUint(uint64(dailyTaskID), 10)
	command.GCommand.ExecuteCommand("pm setdailystatus " + strDailyTaskID + " 2")
	request := &protomsg.MsgCL2GSTaskDailyTaskAwardRequest{}
	request.TaskId = dailyTaskID
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLTaskDailyTaskAwardReply).(*protomsg.MsgGS2CLTaskDailyTaskAwardReply)
	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}
