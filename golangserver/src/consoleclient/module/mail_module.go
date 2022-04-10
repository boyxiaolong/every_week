package module

import (
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GMailTest 测试用例对象
var GMailTest *MailTest

func init() {
	GMailTest = &MailTest{}
	GMailTest.InitCmds("mail")
	GMailTest.RegCommand("operatormail", operatormail)
}

//MailTest 测试用例
type MailTest struct {
	TestBase
}

func operatormail() (bool, int32) {
	command.GCommand.ExecuteCommand("pm addmail 22 33")
	message := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMailNewNotice).(*protomsg.MsgGS2CLMailNewNotice)

	request1 := &protomsg.MsgCL2GSMailListRequest{}
	request1.Id = append(request1.Id, message.Mail.Id)

	readResponse1 := GGameInfo.SendAndWait(request1, msgtype.MsgType_kMsgGS2CLMailListReply).(*protomsg.MsgGS2CLMailListReply)

	if readResponse1.ErrorCode != 0 {
		return readResponse1.ErrorCode == 0, readResponse1.ErrorCode
	}

	request := &protomsg.MsgCL2GSMailOperatorRequest{}
	request.Id = append(request.Id, message.Mail.Id)
	request.ReadFlag = uint32(protomsg.MailReadFlag_kMailFlagRead)

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMailOperatorReply).(*protomsg.MsgGS2CLMailOperatorReply)

	if readResponse.ErrorCode != 0 {
		return readResponse.ErrorCode == 0, readResponse.ErrorCode
	}

	request.ReadFlag = uint32(protomsg.MailReadFlag_kMailFlagExtract)

	readResponse = GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMailOperatorReply).(*protomsg.MsgGS2CLMailOperatorReply)

	if readResponse.ErrorCode != 0 {
		return readResponse.ErrorCode == 0, readResponse.ErrorCode
	}

	request.ReadFlag = uint32(protomsg.MailReadFlag_kMailFlagExtract)

	readResponse = GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMailOperatorReply).(*protomsg.MsgGS2CLMailOperatorReply)

	if readResponse.ErrorCode != 0 {
		return readResponse.ErrorCode == 0, readResponse.ErrorCode
	}

	request.ReadFlag = uint32(protomsg.MailReadFlag_kMailFlagConfirm)

	readResponse = GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMailOperatorReply).(*protomsg.MsgGS2CLMailOperatorReply)

	if readResponse.ErrorCode != 0 {
		return readResponse.ErrorCode == 0, readResponse.ErrorCode
	}

	request.ReadFlag = uint32(protomsg.MailReadFlag_kMailFlagRefuse)

	readResponse = GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMailOperatorReply).(*protomsg.MsgGS2CLMailOperatorReply)

	if readResponse.ErrorCode != 0 {
		return readResponse.ErrorCode == 0, readResponse.ErrorCode
	}

	request.ReadFlag = uint32(protomsg.MailReadFlag_kMailFlagCollect)

	readResponse = GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMailOperatorReply).(*protomsg.MsgGS2CLMailOperatorReply)

	if readResponse.ErrorCode != 0 {
		return readResponse.ErrorCode == 0, readResponse.ErrorCode
	}

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
