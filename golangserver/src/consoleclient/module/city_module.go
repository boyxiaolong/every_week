package module

import (
	"fmt"
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

const (
	testConstructionType uint32 = 1004
)

//GCityTest 测试用例对象
var GCityTest *CityTest

func init() {

	GCityTest = &CityTest{}
	GCityTest.InitCmds("city")
	GCityTest.RegCommand("create", createBuilding)
	GCityTest.RegCommand("upgrade", upgradeBuilding)
	GCityTest.RegCommand("instant", instantBuilding)
	GCityTest.RegCommand("move", moveBuilding)
	GCityTest.RegCommand("querycity", queryCityInfo)
}

//CityTest 测试用例
type CityTest struct {
	TestBase
}

func findLocation(constructionType uint32) (bool, *protomsg.CityCoord) {
	msgCL2GSPMCommandRequest := &protomsg.MsgCL2GSPMCommandRequest{}
	msgCL2GSPMCommandRequest.Command = "findlocation " + strconv.FormatUint(uint64(testConstructionType), 10)
	reply := GGameInfo.SendAndWait(msgCL2GSPMCommandRequest, msgtype.MsgType_kMsgGS2CLCityFindLocationNotice).(*protomsg.MsgGS2CLCityFindLocationNotice)
	if reply.ErrorCode != 0 {
		return false, reply.Coord
	}

	return true, reply.Coord
}

//DoCreateBuilding 创建建筑
func DoCreateBuilding(constructionType uint32) (uint32, int32) {
	msgCL2GSPMCommandRequest := &protomsg.MsgCL2GSPMCommandRequest{}
	msgCL2GSPMCommandRequest.Command = "createbuilding " + strconv.FormatUint(uint64(constructionType), 10) + " 1 1 1"
	reply := GGameInfo.SendAndWait(msgCL2GSPMCommandRequest, msgtype.MsgType_kMsgGS2CLCityCreateBuildingReply).(*protomsg.MsgGS2CLCityCreateBuildingReply)
	if reply.ErrorCode != 0 {
		return 0, reply.ErrorCode
	}

	return reply.BuildingId, 0
}

// CreateBuilding 创建建筑
func CreateBuilding(constructionType uint32, level uint32) (uint32, int32) {
	msgCL2GSPMCommandRequest := &protomsg.MsgCL2GSPMCommandRequest{}
	msgCL2GSPMCommandRequest.Command = "createbuilding " + strconv.FormatUint(uint64(constructionType), 10) + " 1 1 " + strconv.FormatUint(uint64(level), 10)
	reply := GGameInfo.SendAndWait(msgCL2GSPMCommandRequest, msgtype.MsgType_kMsgGS2CLCityCreateBuildingReply).(*protomsg.MsgGS2CLCityCreateBuildingReply)
	if reply.ErrorCode != 0 {
		return 0, reply.ErrorCode
	}

	return reply.BuildingId, 0
}

// SetBuildingLevel 设置建筑等级
func SetBuildingLevel(buildingID uint32, level uint32) {
	command.GCommand.ExecuteCommand(fmt.Sprintf("pm buildinglevel %v %v", buildingID, level))
}

func createBuilding() (bool, int32) {
	command.GCommand.ExecuteCommand("pm resetcity")

	find, coord := findLocation(testConstructionType)
	if !find {
		return false, -1
	}

	request := &protomsg.MsgCL2GSCityCreateBuildingRequest{}
	request.Coord = coord
	request.ConstructionType = testConstructionType

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityCreateBuildingReply).(*protomsg.MsgGS2CLCityCreateBuildingReply)

	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}

func upgradeBuilding() (bool, int32) {
	command.GCommand.ExecuteCommand("pm resetcity")
	command.GCommand.ExecuteCommand("pm buildinglevel 1 1")

	request := &protomsg.MsgCL2GSCityUpgradeBuidlingRequest{}
	request.BuildingId = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityUpgradeBuidlingReply).(*protomsg.MsgGS2CLCityUpgradeBuidlingReply)

	command.GCommand.ExecuteCommand("pm reducebuildtime 1 360000")

	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLCityBuildingUpdateNotice)

	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}

func instantBuilding() (bool, int32) {
	command.GCommand.ExecuteCommand("pm resetcity")
	command.GCommand.ExecuteCommand("pm buildinglevel 1 1")
	command.GCommand.ExecuteCommand("pm rich")

	request := &protomsg.MsgCL2GSCityInstantUpgrageBuildingRequest{}
	request.BuildingId = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityInstantUpgrageBuildingReply).(*protomsg.MsgGS2CLCityInstantUpgrageBuildingReply)

	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}

func moveBuilding() (bool, int32) {
	command.GCommand.ExecuteCommand("pm resetcity")
	buildingID, error := DoCreateBuilding(testConstructionType)
	if buildingID == 0 {
		return false, error
	}

	find, coord := findLocation(testConstructionType)
	if !find {
		return false, -1
	}

	request := &protomsg.MsgCL2GSCityMoveBuildingRequest{}
	request.BuildingId = buildingID
	request.Coord = coord

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityMoveBuildingReply).(*protomsg.MsgGS2CLCityMoveBuildingReply)
	return int32(response.ErrorCode) == 0, int32(response.ErrorCode)
}

func queryCityInfo() (bool, int32) {
	request := &protomsg.MsgCL2GSQueryWorldPlayerCityInfoRequest{}
	request.PlayerId = GGameInfo.PlayerID
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLQueryWorldPlayerCityInfoReply).(*protomsg.MsgGS2CLQueryWorldPlayerCityInfoReply)

	return response.ErrorCode == 0, response.ErrorCode
}
