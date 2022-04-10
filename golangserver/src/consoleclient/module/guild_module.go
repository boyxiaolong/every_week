package module

import (
	"consoleclient/data"
	"fmt"
	"public/command"
	"public/message/error_code"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
	"time"
)

const (
	otherPlayerID  uint64 = 100000
	helperPlayerID uint64 = 1000001
)

var (
	helperGuildID uint64
	prepared      bool
)

//GGuildTest 测试用例对象
var GGuildTest *GuildTest

func init() {

	helperGuildID = 0
	prepared = false

	GGuildTest = &GuildTest{}
	GGuildTest.InitCmds("guild")
	GGuildTest.RegCommand("create", createGuild)
	GGuildTest.RegCommand("join", joinGuild)

	GGuildTest.RegCommand("quitguild", quitGuild)
	GGuildTest.RegCommand("masterquitguild", masterQuitGuild)
	GGuildTest.RegCommand("listapplicant", listApplicant)

	GGuildTest.RegCommand("answerapply", answerApply)
	GGuildTest.RegCommand("queryguild", queryGuild)
	GGuildTest.RegCommand("searchguild", searchGuild)

	GGuildTest.RegCommand("changeposition", changePosition)
	GGuildTest.RegCommand("kickmember", kickMember)
	GGuildTest.RegCommand("setbulletin", setBulletin)
	GGuildTest.RegCommand("setjoincondition", setJoinCondition)

	GGuildTest.RegCommand("updateparam", updateParam)

	GGuildTest.RegCommand("recommendguildlist", recommendGuildList)
	GGuildTest.RegCommand("recommendguild", recommendGuild)

	GGuildTest.RegCommand("guildmail", sendGuildMail)
	GGuildTest.RegCommand("confirmmail", confirmMail)
	GGuildTest.RegCommand("welcome", getWelcomeMail)
	GGuildTest.RegCommand("updatewelcome", updateWelcomeMail)
	GGuildTest.RegCommand("memberrank", queryMemberRankBoard)

	GGuildTest.RegCommand("store", queryGuildStore)
	GGuildTest.RegCommand("queryguildbuildinglist", queryGuildBuildingList)

	GGuildTest.RegCommand("createbuild1", createGuildBuilding1)
	GGuildTest.RegCommand("createbuild", createGuildBuilding)
	GGuildTest.RegCommand("outfire", buildingOutfire)
	GGuildTest.RegCommand("passandtemple", passAndTemple)

	GGuildTest.RegCommand("addguildlabel", addGuildLabel)
	GGuildTest.RegCommand("cancelguildlabel", cancelGuildLabel)

	GGuildTest.RegCommand("guildshopadditem", guildShopAddItem)
	GGuildTest.RegCommand("guildshopbuyitem", guildShopBuyItem)
	GGuildTest.RegCommand("guildshopadditemloglist", guildShopAddItemLogList)
	GGuildTest.RegCommand("guildshopbuyitemloglist", guildShopBuyItemLogList)
	GGuildTest.RegCommand("guildshopitemlist", guildShopItemList)

	//GGuildTest.RegCommand("removebuild", removeBuild)

}

//GuildTest 测试用例
type GuildTest struct {
	TestBase
}

// PrepareHelperGuild 准备帮助公会
func PrepareHelperGuild() {
	if helperGuildID != 0 {
		return
	}

	QuitGuild()
	ForcePlayerQuitGuild(helperPlayerID)
	command.GCommand.ExecuteCommand("pm rich")

	CreateGuild("Help")
	helperGuildID = GetGuildID()
	ForcePlayerJoinGuild(helperGuildID, helperPlayerID)

	ChangePos(helperPlayerID, protomsg.GuildPosition_kGuildPositionMaster)

	QuitGuild()
}

func GetHelperGuildID() uint64 {
	return helperGuildID
}

// PrepareGuildCondition 准备公会条件
func PrepareGuildCondition() {

	if prepared {
		return
	}

	command.GCommand.ExecuteCommand("pm resetcity")
	command.GCommand.ExecuteCommand("pm powercity")

	prepared = true
}

func createGuild() (bool, int32) {
	PrepareGuildCondition()
	QuitGuild()
	command.GCommand.ExecuteCommand("pm rich")

	param := &protomsg.GuildParam{}
	param.Name = "ggggg" + strconv.FormatUint(data.GetPlayerID(), 10)
	param.ShortName = "gg" + strconv.FormatUint(data.GetPlayerID()%100, 10)
	param.Bulletin = "bulletin"
	param.JoinType = protomsg.GuildJoinType_kGuildJoinTypeAllow
	param.Icon = 1

	request := &protomsg.MsgCL2GSGuildCreateRequest{}
	request.Param = param

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildCreateReply).(*protomsg.MsgGS2CLGuildCreateReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func joinGuild() (bool, int32) {
	PrepareGuildCondition()
	QuitGuild()
	PrepareHelperGuild()

	request := &protomsg.MsgCL2GSGuildJoinRequest{}
	request.GuildId = helperGuildID

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildJoinReply).(*protomsg.MsgGS2CLGuildJoinReply)

	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func quitGuild() (bool, int32) {
	PrepareGuildCondition()
	QuitGuild()
	PrepareHelperGuild()
	command.GCommand.ExecuteCommand("pm clearapply")

	request := &protomsg.MsgCL2GSGuildJoinRequest{}
	request.GuildId = helperGuildID
	GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildJoinReply)

	request1 := &protomsg.MsgCL2GSGuildQuitRequest{}
	response1 := GGameInfo.SendAndWait(request1, msgtype.MsgType_kMsgGS2CLGuildQuitReply).(*protomsg.MsgGS2CLGuildQuitReply)

	//common.CHECK_ERROR(int32(response1.ErrorCode), uint16(msgtype.MsgType_kMsgGS2CLGuildQuitReply))

	return response1.ErrorCode == 0, response1.ErrorCode
}

func masterQuitGuild() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("a")
	}

	// 将其他玩家加入公会
	ForcePlayerJoinGuild(GetGuildID(), otherPlayerID)
	time.Sleep(1500 * time.Millisecond)

	// 尝试退出公会
	request := &protomsg.MsgCL2GSGuildQuitRequest{}
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildQuitReply).(*protomsg.MsgGS2CLGuildQuitReply)

	DestroyGuild()

	var result bool = response.ErrorCode == int32(error_code.ErrorCode_kECGuildMasterQuitForbidden) || response.ErrorCode == int32(error_code.ErrorCode_kECGuildNoPermission)
	return result, response.ErrorCode
}

func listApplicant() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("a")
	}

	request := &protomsg.MsgCL2GSGuildListApplicantRequest{}
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildListApplicantReply).(*protomsg.MsgGS2CLGuildListApplicantReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func answerApply() (bool, int32) {
	PrepareGuildCondition()
	QuitGuild()

	request := &protomsg.MsgCL2GSGuildCreateRequest{}
	param := &protomsg.GuildParam{}
	param.Name = "anot" + strconv.FormatUint(data.GetPlayerID()%100, 10)
	param.ShortName = "oi" + strconv.FormatUint(data.GetPlayerID()%100, 10)
	param.Bulletin = "another"
	param.JoinType = protomsg.GuildJoinType_kGuildJoinTypeApply
	param.Icon = 1

	request.Param = param

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildCreateReply).(*protomsg.MsgGS2CLGuildCreateReply)

	strPm := fmt.Sprintf("pm joinguild %v %v", response.GuildId, otherPlayerID)
	command.GCommand.ExecuteCommand(strPm)
	GGameInfo.Wait(msgtype.MsgType_kMsgGS2CLGuildJoinReply)

	request1 := &protomsg.MsgCL2GSGuildAnswerApplyRequest{}
	request1.PlayerId = otherPlayerID
	request1.Agree = false

	response1 := GGameInfo.SendAndWait(request1, msgtype.MsgType_kMsgGS2CLGuildAnswerApplyReply).(*protomsg.MsgGS2CLGuildAnswerApplyReply)

	DestroyGuild()

	return response1.ErrorCode == 0, response1.ErrorCode
}

func queryGuild() (bool, int32) {
	PrepareHelperGuild()

	request := &protomsg.MsgCL2GSGuildQueryRequest{}
	request.GuildId = helperGuildID
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildQuserReply).(*protomsg.MsgGS2CLGuildQuserReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func searchGuild() (bool, int32) {
	request := &protomsg.MsgCL2GSGuildSearchRequest{}

	condition := &protomsg.GuildSearchCondition{}
	condition.Name = "小李子"
	request.SearchCondition = condition
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildSearchReply).(*protomsg.MsgGS2CLGuildSearchReply)

	return response.ErrorCode == 0, response.ErrorCode

}

func kickMember() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("d")
	}

	ForcePlayerJoinGuild(GetGuildID(), otherPlayerID)

	request := &protomsg.MsgCL2GSGuildKickMemberRequest{}
	request.PlayerId = otherPlayerID

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildKickMemberReply).(*protomsg.MsgGS2CLGuildKickMemberReply)

	DestroyGuild()

	return response.ErrorCode == 0, response.ErrorCode

}

func setBulletin() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("c")
	}

	request := &protomsg.MsgCL2GSGuildSetBulletinRequest{}
	request.Bulletin = "bulletin"

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildSetBulletinReply).(*protomsg.MsgGS2CLGuildSetBulletinReply)

	DestroyGuild()

	return response.ErrorCode == 0, response.ErrorCode

}

func setJoinCondition() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("f")
	}

	request := &protomsg.MsgCL2GSGuildSetJoinConditionRequest{}
	condition := &protomsg.GuildJoinCondition{}
	request.Condition = condition
	request.JoinType = protomsg.GuildJoinType_kGuildJoinTypeApply

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildSetJoinConditionReply).(*protomsg.MsgGS2CLGuildSetJoinConditionReply)

	DestroyGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func changePosition() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("e")
	}

	ForcePlayerJoinGuild(GetGuildID(), otherPlayerID)

	request := &protomsg.MsgCL2GSGuildChangePositionRequest{}
	request.PlayerId = otherPlayerID
	request.Position = protomsg.GuildPosition_kGuildPositionSenior

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildChangePositionReply).(*protomsg.MsgGS2CLGuildChangePositionReply)

	DestroyGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func updateParam() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("f")
	}

	request := &protomsg.MsgCL2GSGuildUpdateParamRequest{}
	request.Param = &protomsg.GuildParam{}
	request.Param.Bulletin = "another bulletin"
	request.Param.Icon = 1
	request.Param.JoinType = protomsg.GuildJoinType_kGuildJoinTypeApply
	request.Param.Condition = &protomsg.GuildJoinCondition{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildUpdateParamReply).(*protomsg.MsgGS2CLGuildUpdateParamReply)

	DestroyGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func recommendGuildList() (bool, int32) {
	request := &protomsg.MsgCL2GSRecommendGuildListRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLRecommendGuildListReply).(*protomsg.MsgGS2CLRecommendGuildListReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func recommendGuild() (bool, int32) {
	request := &protomsg.MsgCL2GSRecommendGuildRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLRecommendGuildReply).(*protomsg.MsgGS2CLRecommendGuildReply)

	var result bool = response.ErrorCode == 0 || response.ErrorCode == int32(error_code.ErrorCode_kECGuildGuildNoExist)
	return result, response.ErrorCode
}

func sendGuildMail() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("g")
	}

	request := &protomsg.MsgCL2GSGuildSendMailRequest{}
	request.Subject = "notify mail"
	request.Content = "content"

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildSendMailReply).(*protomsg.MsgGS2CLGuildSendMailReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func confirmMail() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("h")
	}

	command.GCommand.ExecuteCommand("pm richguild 10000000")

	request := &protomsg.MsgCL2GSGuildSendConfirmMailRequest{}
	request.Subject = "notify mail"
	request.Content = "content"
	request.TimeType = protomsg.ConfirmMailTimeLength_kConfirmMailTimeLength10Minutes

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildSendConfirmMailReply).(*protomsg.MsgGS2CLGuildSendConfirmMailReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func getWelcomeMail() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("i")
	}

	request := &protomsg.MsgCL2GSGetWelcomeMailRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGetWelcomeMailReply).(*protomsg.MsgGS2CLGetWelcomeMailReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func updateWelcomeMail() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("j")
	}

	request := &protomsg.MsgCL2GSUpdateWelcomeMailRequest{}
	request.Content = "another"

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLUpdateWelcomeMailReply).(*protomsg.MsgGS2CLUpdateWelcomeMailReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func queryMemberRankBoard() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("k")
	}

	request := &protomsg.MsgCL2GSGuildQueryMemberRankBoardRequest{}
	//request.Type = protomsg.MemberRankBoardType_kMemberRankBoardTypeScience

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildQueryMemberRankBoardReply).(*protomsg.MsgGS2CLGuildQueryMemberRankBoardReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func queryGuildStore() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("l")
	}

	request := &protomsg.MsgCL2GSGuildQueryStoreRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildQueryStoreReply).(*protomsg.MsgGS2CLGuildQueryStoreReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func queryGuildBuildingList() (bool, int32) {
	if GetGuildID() == 0 {
		return true, 0
	}

	request := &protomsg.MsgCL2GSGuildBuildingListRequest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildBuildingListReply).(*protomsg.MsgGS2CLGuildBuildingListReply)

	return response.ErrorCode == 0, response.ErrorCode
}

// CreateGuildFort 创建联盟要塞，返回entity id和错误码
func CreateGuildFort1() (uint64, int32) {
	command.GCommand.ExecuteCommand("pm richguild 10000000")

	searchRequest := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	searchRequest.EntityType = protomsg.EntityType_kEntityType_GuildFort
	searchRequest.SearchLevel = 1
	searchRequest.SearchRange = 30
	searchRequest.NearbyPlayerCastle = true

	// 查找空位
	ok, pos := FindEmptyPos(searchRequest)
	if !ok {
		return 0, int32(error_code.ErrorCode_kECMapNotFoundEmptyPos)
	}

	// 创建行军
	marchIndex, err := CreateMarch()
	if marchIndex == 0 {
		return 0, err
	}

	request := &protomsg.MsgCL2GSConstructGuildBuildingRequest{}
	request.ConfigId = 1
	request.MarchIndex = marchIndex
	request.Position = pos
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLConstructGuildBuildingReply).(*protomsg.MsgGS2CLConstructGuildBuildingReply)
	if response.ErrorCode == 0 {
		return response.EntityId, 0
	}

	return 0, response.ErrorCode
}

func CreateGuildFort() (uint64, int32) {
	command.GCommand.ExecuteCommand("pm richguild 10000000")

	searchRequest := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	searchRequest.EntityType = protomsg.EntityType_kEntityType_GuildFort
	searchRequest.SearchLevel = 1
	searchRequest.SearchRange = 30
	searchRequest.NearbyPlayerCastle = true

	// 查找空位
	ok, pos := FindEmptyPos(searchRequest)
	if !ok {
		return 0, int32(error_code.ErrorCode_kECMapNotFoundEmptyPos)
	}

	request := &protomsg.MsgCL2GSConstructGuildBuildingRequest{}
	request.ConfigId = 8
	request.MarchIndex = 1
	request.Position = pos
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLConstructGuildBuildingReply).(*protomsg.MsgGS2CLConstructGuildBuildingReply)
	if response.ErrorCode == 0 {
		return response.EntityId, 0
	}

	return 0, response.ErrorCode
}

func createGuildBuilding1() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("n")
	}

	entityID, error := CreateGuildFort1()
	QuitGuild()

	if entityID == 0 {
		if error != int32(error_code.ErrorCode_kECGuildBuildingBadArea) {
			return false, error
		}
	}

	// 等待部队回城
	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMarchRemoveNotice, 30)

	return true, 0
}

func createGuildBuilding() (bool, int32) {
	//PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("n")
	}

	entityID, error := CreateGuildFort()
	//QuitGuild()

	if entityID == 0 {
		if error != int32(error_code.ErrorCode_kECGuildBuildingBadArea) {
			return false, error
		}
	}

	// 等待部队回城
	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMarchRemoveNotice, 30)

	return true, 0
}

func buildingOutfire() (bool, int32) {

	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("o")
	}

	entityID, error := CreateGuildFort()
	if 0 == entityID {
		return true, error
	}

	command.GCommand.ExecuteCommand("pm changebuildstatus " + strconv.FormatUint(entityID, 10) + " 2")
	command.GCommand.ExecuteCommand("pm changebuildstatus " + strconv.FormatUint(entityID, 10) + " 3")

	request := &protomsg.MsgCL2GSOutfireRequest{}
	request.ConfigId = 1
	request.EntityId = entityID
	request.Outfire = protomsg.OutfireType_kOutfireTypeIntegral
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLOutfireReply).(*protomsg.MsgGS2CLOutfireReply)

	QuitGuild()
	// 等待部队回城
	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMarchRemoveNotice, 30)

	return response.ErrorCode == 0, response.ErrorCode
}

func passAndTemple() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("p")
	}

	request := &protomsg.MsgCL2GSGuildTempleAndPassRequest{}
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildTempleAndPassReply).(*protomsg.MsgGS2CLGuildTempleAndPassReply)
	QuitGuild()

	return response.ErrorCode == 0, response.ErrorCode
}

func removeBuild() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("q")
	}

	entityID, error := CreateGuildFort()
	if 0 == entityID {
		return true, error
	}

	request := &protomsg.MsgCL2GSGuildBuildingRemoveRequest{}
	request.EntityId = entityID
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildBuildingRemoveReply).(*protomsg.MsgGS2CLGuildBuildingRemoveReply)
	QuitGuild()
	// 等待部队回城
	GGameInfo.WaitSeconds(msgtype.MsgType_kMsgGS2CLMarchRemoveNotice, 30)

	return response.ErrorCode == 0, response.ErrorCode
}

func addGuildLabel() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("r")
	}

	request := &protomsg.MsgCL2GSAddGuildLabelRequest{}
	request.LabelId = 1
	request.Content = "1111"
	request.Pos = &protomsg.Vector2D{}
	request.Pos.X = 1
	request.Pos.Y = 2
	request.RegionId = 1
	request.NoticeMail = false
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLAddGuildLabelReply).(*protomsg.MsgGS2CLAddGuildLabelReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func cancelGuildLabel() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("s")
	}

	request := &protomsg.MsgCL2GSCancelGuildLabelRequest{}
	request.LabelId = 1
	request.RegionId = 1
	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCancelGuildLabelReply).(*protomsg.MsgGS2CLCancelGuildLabelReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func guildShopAddItem() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("t")
	}

	command.GCommand.ExecuteCommand("pm richguild 1000000")

	request := &protomsg.MsgCL2GSGuildShopAddItemResquest{}
	request.ShopItemId = 1002
	request.Count = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildShopAddItemReply).(*protomsg.MsgGS2CLGuildShopAddItemReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func guildShopBuyItem() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("u")
	}

	request := &protomsg.MsgCL2GSGuildShopBuyItemResquest{}
	request.ShopItemId = 1002
	request.Count = 1

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildShopBuyItemReply).(*protomsg.MsgGS2CLGuildShopBuyItemReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func guildShopItemList() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("v")
	}

	request := &protomsg.MsgCL2GSGuildShopBuyItemListResquest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildShopBuyItemListReply).(*protomsg.MsgGS2CLGuildShopBuyItemListReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func guildShopAddItemLogList() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("w")
	}

	request := &protomsg.MsgCL2GSGuildShopAddItemLogListResquest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildShopAddItemLogListReply).(*protomsg.MsgGS2CLGuildShopAddItemLogListReply)

	return response.ErrorCode == 0, response.ErrorCode
}

func guildShopBuyItemLogList() (bool, int32) {
	PrepareGuildCondition()
	if GetGuildID() == 0 {
		CreateGuild("x")
	}

	request := &protomsg.MsgCL2GSGuildShopBuyItemLogListResquest{}

	response := GGameInfo.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildShopBuyItemLogListReply).(*protomsg.MsgGS2CLGuildShopBuyItemLogListReply)

	return response.ErrorCode == 0, response.ErrorCode
}
