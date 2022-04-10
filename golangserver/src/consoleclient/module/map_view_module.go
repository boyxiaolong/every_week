package module

import (
	"public/command"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

//GmapviewTest 测试用例对象
var GmapviewTest *mapviewTest

func init() {
	GmapviewTest = &mapviewTest{}
	GmapviewTest.InitCmds("mapview")
	command.GCommand.RegCommand("createmapview", createview, "mapview")
	GmapviewTest.RegCommand("createmapview2", createview2)
	//GmapviewTest.RegCommand("moveview", moveview)
}

//BaseTest 测试用例
type mapviewTest struct {
	TestBase
}

func createview(str *common.StringParse) (err error) {
	if str.Len() < 1 {
		common.GStdout.Error("login param error")
		return
	}

	regionid := uint64(str.GetUInt64(1))
	request := &protomsg.MsgCL2GSViewMapRequest{}
	request.KingdomId = 1
	request.Position = &protomsg.Vector2D{}
	request.Position.X = 3624000
	request.Position.Y = 3792000
	request.RegionId = regionid
	_ = GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLViewMapReply).(*protomsg.MsgGS2CLViewMapReply)

	return
}

func createview2() (bool, int32) {
	request := &protomsg.MsgCL2GSViewMapRequest{}
	request.KingdomId = 1
	request.Position = &protomsg.Vector2D{}
	request.Position.X = 2052000
	request.Position.Y = 156000

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLViewMapReply).(*protomsg.MsgGS2CLViewMapReply)

	return response.ErrorCode == 0, response.ErrorCode
}

// FindEmptyPos 查找空位
func FindEmptyPos(request *protomsg.MsgCL2GSSearchEmptyPosRequest) (bool, *protomsg.Vector2D) {
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSearchEmptyPosReply).(*protomsg.MsgGS2CLSearchEmptyPosReply)
	if response.ErrorCode != 0 {
		return false, nil
	}

	return true, response.Position
}

// SearchEntity 查找实体
func SearchEntity(searchType protomsg.MapSearchType, level uint32) (uint64, *protomsg.Vector2D) {
	command.GCommand.ExecuteCommand("pm mapsearch " + strconv.FormatUint(uint64(searchType), 10) + " " + strconv.FormatUint(uint64(level), 10) + " 1")
	response := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMapSearchReply).(*protomsg.MsgGS2CLMapSearchReply)
	if response.ErrorCode != 0 {
		return 0, nil
	}

	return response.EntityId, response.Position

}
