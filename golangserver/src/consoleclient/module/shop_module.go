package module

import (
	//"public/command"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GShopTest 测试用例对象
var GShopTest *ShopTest

func init() {
	GShopTest = &ShopTest{}
	GShopTest.InitCmds("shop")
	GShopTest.RegCommand("shopbuy", shopBuy)
	GShopTest.RegCommand("shoprefresh", shopRefresh)
}

//ShopTest 测试用例
type ShopTest struct {
	TestBase
}

func shopBuy() (bool, int32) {
	command.GCommand.ExecuteCommand("pm rich 100000")
	request := &protomsg.MsgCL2GSShopBuyItemRequest{}
	request.Id = 20001
	request.Count = 2
	request.Type = 2

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLShopBuyItemReply).(*protomsg.MsgGS2CLShopBuyItemReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func shopRefresh() (bool, int32) {
	command.GCommand.ExecuteCommand("pm rich 100000")

	request := &protomsg.MsgCL2GSShopRefreshRequest{}
	request.Type = 2

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLShopRefreshReply).(*protomsg.MsgGS2CLShopRefreshReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
