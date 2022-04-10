package module

import (
	//"public/command"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

//GItemTest 测试用例对象
var GItemTest *ItemTest

func init() {
	GItemTest = &ItemTest{}
	GItemTest.InitCmds("item")
	GItemTest.RegCommand("useitem", useItem)
	GItemTest.RegCommand("buyitem", buyItem)
	GItemTest.RegCommand("buyanduseitem", buyAndUseItem)
}

//ItemTest 测试用例
type ItemTest struct {
	TestBase
}

func useItem() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 10000 10")

	request := &protomsg.MsgCL2GSItemUseRequest{}
	request.Id = 10000
	request.Count = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLItemUseReply).(*protomsg.MsgGS2CLItemUseReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func buyItem() (bool, int32) {
	command.GCommand.ExecuteCommand("pm rich 100000")

	request := &protomsg.MsgCL2GSItemBuyRequest{}
	request.Id = 10000
	request.Count = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLItemBuyReply).(*protomsg.MsgGS2CLItemBuyReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func buyAndUseItem() (bool, int32) {
	command.GCommand.ExecuteCommand("pm rich 100000")

	request := &protomsg.MsgCL2GSItemBuyAndUseRequest{}
	request.Id = 10000
	request.Count = 1

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLItemBuyAndUseReply).(*protomsg.MsgGS2CLItemBuyAndUseReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// DoUseItem 使用道具
func DoUseItem(id uint32, count uint32, params []uint64) (bool, int32) {
	request := &protomsg.MsgCL2GSItemUseRequest{}
	request.Id = id
	request.Count = count
	request.Params = params
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLItemUseReply).(*protomsg.MsgGS2CLItemUseReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
