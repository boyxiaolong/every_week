package module

import (
	//"public/command"

	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

//GRecruitTest 测试用例对象
var GPompeiiTest *PompeiiTest

func init() {
	GPompeiiTest = &PompeiiTest{}
	GPompeiiTest.InitCmds("pompeii")
	//GPompeiiTest.RegCommand("pompeiiEvent", pompeiiEvent)
	GPompeiiTest.RegCommand("cure", cure)
	GPompeiiTest.RegCommand("instantcure", instantCure)
	GPompeiiTest.RegCommand("speedup", speedUpCure)
}

//RecruitTest 测试用例
type PompeiiTest struct {
	TestBase
}

func pompeiiEvent() (bool, int32) {
	request := &protomsg.MsgCL2GSPompeiiEventQueryRequest{}
	request.RegionId = 0

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPompeiiEventQueryReply).(*protomsg.MsgGS2CLPompeiiEventQueryReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func cure() (bool, int32) {

	troopID := 101
	count := 1000

	// 准备公会
	PrepareGuildCondition()
	QuitGuild()
	PrepareHelperGuild()

	// 创建副本
	command.GCommand.ExecuteCommand("pm createpompeii 1 " + strconv.FormatUint(GetGuildID(), 10) + " " + strconv.FormatUint(GetHelperGuildID(), 10))

	notice := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLPMCreateInstanceNotice).(*protomsg.MsgGS2CLPMCreateInstanceNotice)
	if notice.ErrorCode != 0 {
		return false, notice.ErrorCode
	}

	regionID := notice.RegionId
	command.GCommand.ExecuteCommand("pm intoinstance " + strconv.FormatUint(regionID, 10) + " 1")
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMapDataReply)
	command.GCommand.ExecuteCommand("pm addpompeiiwound " + strconv.FormatUint(uint64(troopID), 10) + " " + strconv.FormatUint(uint64(count), 10))

	cureRequest := &protomsg.MsgCL2GSPompeiiHospitalCureRequest{}
	army := make([]*protomsg.TroopData, 0)
	troop := &protomsg.TroopData{}
	troop.TroopId = uint32(troopID)
	troop.TroopUints = uint32(count)
	army = append(army, troop)
	cureRequest.ArmyInfo = army
	cureResponse := GGameInfo.SendAndWait(cureRequest, msgtype.MsgType_kMsgGS2CLPompeiiHospitalCureReply).(*protomsg.MsgGS2CLPompeiiHospitalCureReply)
	command.GCommand.ExecuteCommand("pm closeinstance " + strconv.FormatUint(regionID, 10))

	return cureResponse.ErrorCode == 0, cureResponse.ErrorCode
}

func instantCure() (bool, int32) {

	troopID := 101
	count := 1000

	command.GCommand.ExecuteCommand("pm rich 100000")
	// 准备公会
	PrepareGuildCondition()
	QuitGuild()
	PrepareHelperGuild()

	// 创建副本
	command.GCommand.ExecuteCommand("pm createpompeii 1 " + strconv.FormatUint(GetGuildID(), 10) + " " + strconv.FormatUint(GetHelperGuildID(), 10))

	notice := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLPMCreateInstanceNotice).(*protomsg.MsgGS2CLPMCreateInstanceNotice)
	if notice.ErrorCode != 0 {
		return false, notice.ErrorCode
	}

	regionID := notice.RegionId
	command.GCommand.ExecuteCommand("pm intoinstance " + strconv.FormatUint(regionID, 10) + " 1")
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMapDataReply)
	command.GCommand.ExecuteCommand("pm addpompeiiwound " + strconv.FormatUint(uint64(troopID), 10) + " " + strconv.FormatUint(uint64(count), 10))

	cureRequest := &protomsg.MsgCL2GSPompeiiHospitalInstantCureRequest{}
	army := make([]*protomsg.TroopData, 0)
	troop := &protomsg.TroopData{}
	troop.TroopId = uint32(troopID)
	troop.TroopUints = uint32(count)
	army = append(army, troop)
	cureRequest.ArmyInfo = army
	cureResponse := GGameInfo.SendAndWait(cureRequest, msgtype.MsgType_kMsgGS2CLPompeiiHospitalInstantCureReply).(*protomsg.MsgGS2CLPompeiiHospitalInstantCureReply)

	command.GCommand.ExecuteCommand("pm closeinstance " + strconv.FormatUint(regionID, 10))

	return cureResponse.ErrorCode == 0, cureResponse.ErrorCode
}

func speedUpCure() (bool, int32) {

	troopID := 101
	count := 1000

	// 准备公会
	PrepareGuildCondition()
	QuitGuild()
	PrepareHelperGuild()

	// 创建副本
	command.GCommand.ExecuteCommand("pm createpompeii 1 " + strconv.FormatUint(GetGuildID(), 10) + " " + strconv.FormatUint(GetHelperGuildID(), 10))

	notice := GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLPMCreateInstanceNotice).(*protomsg.MsgGS2CLPMCreateInstanceNotice)
	if notice.ErrorCode != 0 {
		return false, notice.ErrorCode
	}

	regionID := notice.RegionId
	command.GCommand.ExecuteCommand("pm intoinstance " + strconv.FormatUint(regionID, 10) + " 1")
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLMapDataReply)
	command.GCommand.ExecuteCommand("pm addpompeiiwound " + strconv.FormatUint(uint64(troopID), 10) + " " + strconv.FormatUint(uint64(count), 10))

	cureRequest := &protomsg.MsgCL2GSPompeiiHospitalCureRequest{}
	army := make([]*protomsg.TroopData, 0)
	troop := &protomsg.TroopData{}
	troop.TroopId = uint32(troopID)
	troop.TroopUints = uint32(count)
	army = append(army, troop)
	cureRequest.ArmyInfo = army
	cureResponse := GGameInfo.SendAndWait(cureRequest, msgtype.MsgType_kMsgGS2CLPompeiiHospitalCureReply).(*protomsg.MsgGS2CLPompeiiHospitalCureReply)
	if cureResponse.ErrorCode != 0 {
		return false, cureResponse.ErrorCode
	}

	request := &protomsg.MsgCL2GSPompeiiHospitalInstantSpeedUpRequest{}
	request.SpeedUpSeconds = 10000
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPompeiiHospitalInstantSpeedUpReply).(*protomsg.MsgGS2CLPompeiiHospitalInstantSpeedUpReply)

	command.GCommand.ExecuteCommand("pm closeinstance " + strconv.FormatUint(regionID, 10))

	return response.ErrorCode == 0, response.ErrorCode
}
