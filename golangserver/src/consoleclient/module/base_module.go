package module

import (
	"consoleclient/data"
	"public/command"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GBaseTest 测试用例对象
var GBaseTest *BaseTest

var (
	isSuperman bool
)

func init() {
	GBaseTest = &BaseTest{}
	GBaseTest.InitCmds("base")
	GBaseTest.RegCommand("rename", rename)
	GBaseTest.RegCommand("querykillinfo", queryKillInfo)
	GBaseTest.RegCommand("querybattlestatistics", queryBattleStatistics)
	GBaseTest.RegCommand("queryplayerlordinfo", queryPlayerLordInfo)
	GBaseTest.RegCommand("setdisplayphoto", setDisplayPhoto)
	GBaseTest.RegCommand("setdisplayphotoframe", setDisplayPhotoFrame)
	GBaseTest.RegCommand("gem2res", gem2Res)
}

//BaseTest 测试用例
type BaseTest struct {
	TestBase
}

func rename() (bool, int32) {

	//data1 := data.GDataCenter.GetData(data.GetPlayerName)
	//data1 = data1.(string)
	data1 := data.GetPlayerName()
	common.GStdout.Success("get player base ,name %v", data1)
	data.GDataCenter.SetData(data.SetPlayerName, "2222222")

	return true, 0
}

func queryKillInfo() (bool, int32) {
	request := &protomsg.MsgCL2GSKillRequest{}
	request.PlayerId = data.GetPlayerID()

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLKillReply).(*protomsg.MsgGS2CLKillReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func queryBattleStatistics() (bool, int32) {
	request := &protomsg.MsgCL2GSBattleStatisticsRequest{}
	request.PlayerId = data.GetPlayerID()

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLBattleStatisticsReply).(*protomsg.MsgGS2CLBattleStatisticsReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func queryPlayerLordInfo() (bool, int32) {
	request := &protomsg.MsgCL2GSQueryPlayerLordInfoRequest{}
	request.PlayerId = GGameInfo.PlayerID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLQueryPlayerLordInfoReply).(*protomsg.MsgGS2CLQueryPlayerLordInfoReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func setDisplayPhoto() (bool, int32) {
	request := &protomsg.MsgCL2GSSetDisplayPhotoRequest{}
	request.Displayphoto = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSetDisplayPhotoReply).(*protomsg.MsgGS2CLSetDisplayPhotoReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func setDisplayPhotoFrame() (bool, int32) {
	request := &protomsg.MsgCL2GSSetDisplayPhotoFrameRequest{}
	request.DisplayPhotoFrame = 2

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSetDisplayPhotoFrameReply).(*protomsg.MsgGS2CLSetDisplayPhotoFrameReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func gem2Res() (bool, int32) {
	command.GCommand.ExecuteCommand("pm addcurrency 2 1000")
	//GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLPlayerAttribute)

	request := &protomsg.MsgCL2GSGemToResource{}
	list := make([]*protomsg.MsgCL2GSGemToResource_BuyResource, 0)

	resource := &protomsg.MsgCL2GSGemToResource_BuyResource{}
	resource.Type = protomsg.CurrencyType(protomsg.CurrencyType_kCurrencyTypeOil)
	resource.Num = 20

	list = append(list, resource)
	request.Resources = list

	GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerAttribute)
	return true, 0
}
