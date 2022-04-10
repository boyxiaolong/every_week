package module

import (
	"public/command"
	"public/message/msgtype"
	"public/message/protomsg"
)

// "public/command"

// GHeroTest 测试用例对象
var GHeroTest *HeroTest

func init() {
	GHeroTest = &HeroTest{}
	GHeroTest.InitCmds("hero")
	GHeroTest.RegCommand("summonhoro", summonHoro)
	GHeroTest.RegCommand("exchangeherodebris", exchangeHeroDebris)
	GHeroTest.RegCommand("herokillinfo", heroKillInfo)
	GHeroTest.RegCommand("upgradestar", upgradeStar)
	GHeroTest.RegCommand("upgradeskill", upgradeSkill)
	GHeroTest.RegCommand("awakskill", awakSkill)
	GHeroTest.RegCommand("usepoint", usePoint)
	GHeroTest.RegCommand("getglobalpoints", getGlobalPoints)
	GHeroTest.RegCommand("resetpoint", resetPoint)
}

// HeroTest 测试用例
type HeroTest struct {
	TestBase
}

// 召唤英雄
func summonHoro() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 51003 10")
	request := &protomsg.MsgCL2GSSummonHoroRequest{}
	request.HeroId = 1003

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSummonHoroReply).(*protomsg.MsgGS2CLSummonHoroReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 兑换碎片
func exchangeHeroDebris() (bool, int32) {
	command.GCommand.ExecuteCommand("pm addhero 1001 1")
	command.GCommand.ExecuteCommand("pm additem 51001 3")
	request := &protomsg.MsgCL2GSExchangeHeroDebrisRequest{}
	request.HeroId = 1001
	request.Exchange = 51001
	request.ExchangeCount = 3

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLExchangeHeroDebrisReply).(*protomsg.MsgGS2CLExchangeHeroDebrisReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 升级技能
func upgradeSkill() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 51001 100")
	command.GCommand.ExecuteCommand("pm resethero 1001")
	request := &protomsg.MsgCL2GSUpgradeSkillRequest{}
	request.HeroId = 1001
	request.SkillId = 10011
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLUpgradeSkillReply).(*protomsg.MsgGS2CLUpgradeSkillReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 获取英雄击杀数据
func heroKillInfo() (bool, int32) {
	request := &protomsg.MsgCL2GSHeroKillRequest{}
	request.HeroId = 1001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLHeroKillReply).(*protomsg.MsgGS2CLHeroKillReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 升星
func upgradeStar() (bool, int32) {
	command.GCommand.ExecuteCommand("pm upgradehero 1001 100000000")
	command.GCommand.ExecuteCommand("pm additem 52000 30")
	request := &protomsg.MsgCL2GSUpgradeStarRequest{}
	request.HeroId = 1001
	request.Items = make([]*protomsg.StarItemInfo, 0)
	request.Items = append(request.Items, &protomsg.StarItemInfo{ItemId: 52000, Count: 30})
	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLUpgradeStarReply).(*protomsg.MsgGS2CLUpgradeStarReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 觉醒技能
func awakSkill() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 51001 10000")

	for i := 0; i < 6; i++ {
		command.GCommand.ExecuteCommand("pm upgradeheroskill 1001")
	}

	request := &protomsg.MsgCL2GSAwakSkillRequest{}
	request.HeroId = 1001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLAwakSkillReply).(*protomsg.MsgGS2CLAwakSkillReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 加点
func usePoint() (bool, int32) {
	request := &protomsg.MsgCL2GSUsePointRequest{}
	request.HeroId = 1001
	request.Info = make([]*protomsg.PointInfo, 0)
	request.Info = append(request.Info, &protomsg.PointInfo{Type: 1, Count: 2})

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLUsePointReply).(*protomsg.MsgGS2CLUsePointReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

func getGlobalPoints() (bool, int32) {
	request := &protomsg.MsgCL2GSGetGlobalHeroPointsRequest{}
	request.HeroId = 1001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGetGlobalHeroPointsReply).(*protomsg.MsgGS2CLGetGlobalHeroPointsReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 重置点数
func resetPoint() (bool, int32) {
	command.GCommand.ExecuteCommand("pm additem 40120 1")
	request := &protomsg.MsgCL2GSResetPointRequest{}
	request.HeroId = 1001

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLResetPointReply).(*protomsg.MsgGS2CLResetPointReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
