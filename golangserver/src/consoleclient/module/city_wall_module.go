package module

import (
	//"public/command"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GCityWallTest 测试用例对象
var GCityWallTest *CityWallTest

func init() {
	GCityWallTest = &CityWallTest{}
	GCityWallTest.InitCmds("citywall")
	GCityWallTest.RegCommand("repair", repair)
	GCityWallTest.RegCommand("putoutfire", putOutFire)
	GCityWallTest.RegCommand("garrison", garrison)
}

//CityWallTest 测试用例
type CityWallTest struct {
	TestBase
}

func repair() (bool, int32) {
	command.GCommand.ExecuteCommand("pm citywallburn")

	request := &protomsg.MsgCL2GSRepairWallRequest{}

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLRepairWallReply).(*protomsg.MsgGS2CLRepairWallReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func putOutFire() (bool, int32) {
	command.GCommand.ExecuteCommand("pm citywallburn")

	request := &protomsg.MsgCL2GSPutOutFireRequest{}
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPutOutFireReply).(*protomsg.MsgGS2CLPutOutFireReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func garrison() (bool, int32) {
	request := &protomsg.MsgCL2GSSelectGarrisonHeroRequest{}
	request.MainHeroId = 1001
	request.SecondHeroId = 1002
	request.RegionId = 0

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSelectGarrisonHeroReply).(*protomsg.MsgGS2CLSelectGarrisonHeroReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
