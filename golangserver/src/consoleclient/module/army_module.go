package module

import (
	"consoleclient/data"
	"public/command"
	"public/message/error_code"
	"public/message/msgtype"
	"public/message/protomsg"
	"time"
)

// "public/command"

// GArmyTest 测试用例对象
var GArmyTest *ArmyTest

func init() {
	GArmyTest = &ArmyTest{}
	GArmyTest.InitCmds("army")
	GArmyTest.RegCommand("trainarmy", trainArmy)
	GArmyTest.RegCommand("cancel", cancelTrain)
	GArmyTest.RegCommand("trainarmy", trainArmy)
	GArmyTest.RegCommand("instantspeedup", instantSpeedUp)
	//GArmyTest.RegCommand("trainarmy", trainArmy)
	//GArmyTest.RegCommand("dealarmy", dealArmy)
	GArmyTest.RegCommand("uprank", upRank)
	GArmyTest.RegCommand("cancel", cancelTrain)
	GArmyTest.RegCommand("disband", disband)

	GArmyTest.RegCommand("instanttrain", instantTrain)
	GArmyTest.RegCommand("instantuprank", instantUpRank)
}

// ArmyTest 测试用例
type ArmyTest struct {
	TestBase
}

func getBarrackID() uint32 {
	buildType := uint32(protomsg.ConstructionType_kConstructionType_Barrack)
	count := data.CountBuildsByType(buildType)
	if count == 0 {
		_, error := DoCreateBuilding(buildType)
		if error != 0 {
			return 0
		}
	}
	buildings := data.GetBuildsByType(buildType)
	if buildings == nil {
		return 0
	}
	return buildings[0].ID
}

// 训练
func trainArmy() (bool, int32) {
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}

	request := &protomsg.MsgCL2GSTrainRequest{}
	request.Army = &protomsg.ArmyInfo{}
	request.Army.ArmyId = 141
	request.Army.Count = 10
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLTrainReply).(*protomsg.MsgGS2CLTrainReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 收割
func dealArmy() (bool, int32) {
	time.Sleep(2 * time.Second)
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}
	request := &protomsg.MsgCL2GSDealArmyRequest{}
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLDealArmyReply).(*protomsg.MsgGS2CLDealArmyReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 进阶
func upRank() (bool, int32) {
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}

	command.GCommand.ExecuteCommand("pm upgradetechnologylv 5004")

	request := &protomsg.MsgCL2GSUpRankRequest{}
	request.Army = &protomsg.ArmyInfo{}
	request.Army.ArmyId = 141
	request.Army.Count = 50
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLUpRankReply).(*protomsg.MsgGS2CLUpRankReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 解散
func disband() (bool, int32) {
	request := &protomsg.MsgCL2GSDisbandRequest{}
	request.Army = &protomsg.ArmyInfo{}
	request.Army.ArmyId = 141
	request.Army.Count = 10

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLDisbandReply).(*protomsg.MsgGS2CLDisbandReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 取消训练
func cancelTrain() (bool, int32) {
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}
	request := &protomsg.MsgCL2GSCancelTrainRequest{}
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCancelTrainReply).(*protomsg.MsgGS2CLCancelTrainReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 立即训练
func instantTrain() (bool, int32) {
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}
	request := &protomsg.MsgCL2GSInstantTrainRequest{}
	request.Army = &protomsg.ArmyInfo{}
	request.Army.ArmyId = 141
	request.Army.Count = 10
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLInstantTrainReply).(*protomsg.MsgGS2CLInstantTrainReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 立即进阶
func instantUpRank() (bool, int32) {
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}
	request := &protomsg.MsgCL2GSInstantUpRankRequest{}
	request.Army = &protomsg.ArmyInfo{}
	request.Army.ArmyId = 141
	request.Army.Count = 10
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLInstantUpRankReply).(*protomsg.MsgGS2CLInstantUpRankReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}

// 立即加速
func instantSpeedUp() (bool, int32) {
	barrackID := getBarrackID()
	if barrackID == 0 {
		return false, int32(error_code.ErrorCode_kECArmyNoBarrack)
	}
	request := &protomsg.MsgCL2GSInstantSpeedUpRequest{}
	request.BarrackId = barrackID

	readResponse := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLInstantSpeedUpReply).(*protomsg.MsgGS2CLInstantSpeedUpReply)

	return readResponse.ErrorCode == 0, readResponse.ErrorCode
}
