package module

import (
	"public/message/msgtype"
	"public/message/protomsg"
)

// ExpeTest comment
type ExpeTest struct {
	TestBase
}

//GExpeTest comment
var GExpeTest *ExpeTest

func init() {
	GExpeTest = &ExpeTest{}
	GExpeTest.InitCmds("expe")
	GExpeTest.RegCommand("level_list", testLevelList)
	GExpeTest.RegCommand("level_reward", testLevelReward)
	GExpeTest.RegCommand("batch_collect", testBatchCollect)
	GExpeTest.RegCommand("challenge_level", testChallengeLevel)
	GExpeTest.RegCommand("exit_level", testExitLevel)
}

//testLevelList comment
func testLevelList() (bool, int32) {
	request := &protomsg.MsgCL2GSPlayerExpeLevelListRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeLevelListReply).(*protomsg.MsgGS2CLPlayerExpeLevelListReply)

	return response.ErrorCode == 0, response.ErrorCode
}

//testLevelReward comment
func testLevelReward() (bool, int32) {
	var tmp = &protomsg.MsgCL2GSPMCommandRequest{}
	tmp.Command = "expe_refresh"
	GGameInfo.SendAndWait(tmp, msgtype.MsgType_kMsgGS2CLPMCommandReply)
	tmp.Command = "expe_reset"
	GGameInfo.SendAndWait(tmp, msgtype.MsgType_kMsgGS2CLPMCommandReply)
	tmp.Command = "expe_unlock_all"
	GGameInfo.SendAndWait(tmp, msgtype.MsgType_kMsgGS2CLPMCommandReply)

	request := &protomsg.MsgCL2GSPlayerExpeLevelRewardRequest{}
	request.LevelId = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeLevelRewardReply).(*protomsg.MsgGS2CLPlayerExpeLevelRewardReply)

	return response.ErrorCode == 0, response.ErrorCode
}

//testBatchCollect comment
func testBatchCollect() (bool, int32) {
	var tmp = &protomsg.MsgCL2GSPMCommandRequest{}
	tmp.Command = "expe_unlock_all"
	GGameInfo.SendAndWait(tmp, msgtype.MsgType_kMsgGS2CLPMCommandReply)

	request := &protomsg.MsgCL2GSPlayerExpeBatchCollectRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeBatchCollectReply).(*protomsg.MsgGS2CLPlayerExpeBatchCollectReply)

	return response.ErrorCode == 0, response.ErrorCode
}

//testChallengeLevel comment
func testChallengeLevel() (bool, int32) {
	request := &protomsg.MsgCL2GSPlayerExpeChallengeLevelRequest{}
	request.LevelId = 1
	// request.ArmyArray  // 我方阵容信息
	// request.HqArmy  // 我方主堡守军

	var armyData = &protomsg.ArmyData{}
	armyData.Hero1 = 1034
	var troopData = &protomsg.TroopData{}
	troopData.TroopId = 105
	troopData.TroopUints = 1000
	armyData.CurrentTroops = append(armyData.CurrentTroops, troopData)

	request.ArmyArray = append(request.ArmyArray, armyData)

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeChallengeLevelReply).(*protomsg.MsgGS2CLPlayerExpeChallengeLevelReply)

	return response.ErrorCode == 0, response.ErrorCode
}

//testExitLevel comment
func testExitLevel() (bool, int32) {
	request := &protomsg.MsgCL2GSPlayerExpeExitLevelRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeExitLevelReply).(*protomsg.MsgGS2CLPlayerExpeExitLevelReply)

	return response.ErrorCode == 0, response.ErrorCode
}
