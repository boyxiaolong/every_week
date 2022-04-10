package player

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	BTree "public/behaviortree"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
	"robotclient/client/player/data"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("CreateGuild", CreateCreateGuildTask)
	BTree.RegisterTaskCreator("JoinGuild", CreateJoinGuildTask)
	BTree.RegisterTaskCreator("LeaveGuild", CreateLeaveGuildTask)
	BTree.RegisterTaskCreator("DismissGuild", CreateDismissGuildTask)
	BTree.RegisterTaskCreator("GuildBattle", CreateGuildBattleTask)
	BTree.RegisterTaskCreator("GuildBattleToTarget", CreateGuildBattleToTargetTask)
	BTree.RegisterTaskCreator("JoinGuildMarch", CreateJoinGuildMarchTask)
	BTree.RegisterTaskCreator("SendGuildMail", CreateSendGuildMailTask)
	BTree.RegisterTaskCreator("GuildAssistHelp", CreateGuildAssistHelpTask)
	BTree.RegisterConditionCallback("IsInGuild", IsInGuild)
	BTree.RegisterConditionCallback("CheckGuildPosition", CheckGuildPosition)
	BTree.RegisterTaskCreator("MoveCastleToGuild", CreateMoveCastleToGuildTask)
}

func IsInGuild(testobj BTree.ObjectInterface, params string) bool {
	player := testobj.(*Player)
	guild_id := player.GetData("GuildData.GuildID").(uint64)
	return guild_id != 0
}

func (player *Player) GetGuildId() uint64 {
	guild_id := player.GetData("GuildData.GuildID").(uint64)
	return guild_id
}

func CheckGuildPosition(testobj BTree.ObjectInterface, params string) bool {
	player := testobj.(*Player)
	check_position, error := strconv.Atoi(params)
	if error != nil {
		common.GStdout.Debug("CheckGuildPosition =====")
		return false
	}

	guild_position := (int)(player.GetData("GuildData.GuildPosition").(protomsg.GuildPosition))
	common.GStdout.Error("CheckGuildPosition %v %v %v ", guild_position, check_position, guild_position >= check_position)
	return guild_position >= check_position
}

// RangeEntities  通知地图数据
func (p *Player) RangeGuildEntities(fn func(entityID uint64, entityData *data.SyncEntityData) bool) {
	if mainRegion, ok := p.GetMainRegion(); ok {
		mainRegion.GetGuildCenter().RangeEntities(fn)
	}
}

// RangeGuildMarchs  通知地图数据
func (p *Player) RangeGuildMarchs(fn func(entityID uint64, entityData *data.SyncEntityData) bool) {
	p.RangeGuildEntities(func(entityID uint64, entityData *data.SyncEntityData) bool {
		if entityData.EntityType == (uint32)(protomsg.EntityType_kEntityType_GuildMarch) {
			fn(entityID, entityData)
		}
		return true
	})
}

// CheckGuildEntityExist  通知地图数据
func (p *Player) CheckGuildEntityExist(marchEntityID uint64) bool {
	if mainRegion, ok := p.GetMainRegion(); ok {
		_, ok := mainRegion.GetGuildCenter().GetEntity(marchEntityID)
		return ok
	}
	return false
}

// WaitMarchBack a
func (p *Player) WaitGuildEntityDisband(marchEntityID uint64, waitTime int) bool {
	if waitTime == 0 {
		for true {
			time.Sleep(10 * time.Second)
			if p.CheckGuildEntityExist(marchEntityID) {
				continue
			}
			return true
		}
	}

	for i := 0; i < waitTime; i++ {
		time.Sleep(1 * time.Second)
		if p.CheckGuildEntityExist(marchEntityID) {
			continue
		}
		return true
	}
	return false
}

func (p *Player) CreateGuildMarch(army *protomsg.ArmyData, command *protomsg.MarchCommand, wait_time_index uint32) (uint64, bool) {
	request := &protomsg.MsgCL2GSCreateGuildMarchRequest{}
	request.Command = command
	request.ArmyData = army
	request.WaitTimeIndex = wait_time_index
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateGuildMarchReply).(*protomsg.MsgGS2CLCreateGuildMarchReply)
	if !ok {
		return 0, false
	}
	time.Sleep(1 * time.Second) // 等待数据同步
	if int32(response.ErrorCode) == 0 {
		return response.GuildMarchEntityId, true
	}
	return 0, false
}

func CreateCreateGuildTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &CreateGuildTask{}
	return task, nil
}

// CreateGuildTask comment
type CreateGuildTask struct {
}

func (t *CreateGuildTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return player.CreateGuild()
}

const beginLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		if i == 0 {
			b[i] = beginLetterBytes[rand.Intn(len(beginLetterBytes))]
		} else {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
	}
	return string(b)
}

func (player *Player) CreateGuild() bool {
	player.PM("rich")

	param := &protomsg.GuildParam{}
	param.Name = RandStringBytes(24)
	param.ShortName = RandStringBytes(5)
	param.Bulletin = "bulletin"
	param.JoinType = protomsg.GuildJoinType_kGuildJoinTypeAllow

	request := &protomsg.MsgCL2GSGuildCreateRequest{}
	request.Param = param

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildCreateReply).(*protomsg.MsgGS2CLGuildCreateReply)
	if !ok {
		return false
	}

	common.GStdout.Debug("=================== guild name %v ", param.Name)
	return int32(response.ErrorCode) == 0
}

// CreateJoinGuildTask comment
func CreateJoinGuildTask(params []string) (res BTree.BTTaskInterface, err error) {
	var guildID uint64
	guildID = 0
	if len(params) >= 1 {
		id, error := strconv.ParseUint(params[0], 10, 64)
		if error != nil {
			return nil, error
		}
		guildID = id
	}

	task := &JoinGuildTask{guildID: guildID}
	return task, nil
}

// JoinGuildTask comment
type JoinGuildTask struct {
	guildID uint64
}

// DoTask comment
func (t *JoinGuildTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	if player.GetGuildId() != 0 {
		return true
	}

	if t.guildID != 0 {
		return player.JoinGuild(t.guildID)
	}

	searchRequest := &protomsg.MsgCL2GSRecommendGuildListRequest{}
	searchResponse, ok := player.SendAndWait(searchRequest, msgtype.MsgType_kMsgGS2CLRecommendGuildListReply).(*protomsg.MsgGS2CLRecommendGuildListReply)
	if !ok {
		return false
	}

	if searchResponse.ErrorCode != 0 {
		return false
	}

	for _, v := range searchResponse.Guilds {
		if player.JoinGuild(v.GuildId) {
			return true
		}
	}
	return false
}

//JoinGuild a
func (p *Player) JoinGuild(join_guild_id uint64) bool {
	request := &protomsg.MsgCL2GSGuildJoinRequest{}
	request.GuildId = join_guild_id

	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildJoinReply).(*protomsg.MsgGS2CLGuildJoinReply)
	if !ok {
		common.GStdout.Debug("=================== join guild send wait failed timeout...  player_id: %v ", p.PlayerID)
		return false
	}

	if response.ErrorCode != 0 {
		common.GStdout.Debug("=================== join guild failed player_id:%v errorcode:%v ", p.PlayerID, response.ErrorCode)
	}
	return int32(response.ErrorCode) == 0
}

func CreateLeaveGuildTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &LeaveGuildTask{}
	return task, nil
}

// LeaveGuildTask comment
type LeaveGuildTask struct {
}

func (t *LeaveGuildTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return player.LeaveGuild()
}

//Login Login
func (player *Player) LeaveGuild() bool {
	request := &protomsg.MsgCL2GSGuildQuitRequest{}

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildQuitReply).(*protomsg.MsgGS2CLGuildQuitReply)
	if !ok {
		return false
	}

	return int32(response.ErrorCode) == 0
}

// DismissGuildTask comment
type DismissGuildTask struct {
}

func (t *DismissGuildTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return player.DismissGuild()
}

//DismissGuild
func (player *Player) DismissGuild() bool {
	request := &protomsg.MsgCL2GSGuildDismissRequest{}

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildDismissReply).(*protomsg.MsgGS2CLGuildDismissReply)
	if !ok {
		return false
	}
	return int32(response.ErrorCode) == 0
}

func CreateDismissGuildTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &DismissGuildTask{}
	return task, nil
}

// GuildBattleTask comment
type GuildBattleTask struct {
	search_type      protomsg.MapSearchType
	wait_time_index  uint32
	waitMarchDisband bool
	minArmy          uint32
	maxArmy          uint32
}

// GuildBattleTaskData comment
type GuildBattleTaskData struct {
	guildMarchEntityID uint64
}

func CreateGuildBattleTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 5 {
		return nil, fmt.Errorf("create Guild Battle params error!")
	}
	search_type, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	wait_time_index, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	waitMarchDisband, error := strconv.ParseBool(params[2])
	if error != nil {
		return nil, error
	}

	minArmy, error := strconv.Atoi(params[3])
	if error != nil {
		return nil, error
	}

	maxArmy, error := strconv.Atoi(params[4])
	if error != nil {
		return nil, error
	}

	task := &GuildBattleTask{search_type: (protomsg.MapSearchType)(search_type), wait_time_index: (uint32)(wait_time_index), waitMarchDisband: waitMarchDisband, minArmy: (uint32)(minArmy), maxArmy: (uint32)(maxArmy)}
	return task, nil
}

func (t *GuildBattleTask) getTaskData(runtask *BTree.BTreeRunTask, taskindex uint32) *GuildBattleTaskData {
	data := runtask.GetTaskData(taskindex)
	if data != nil {
		return data.(*GuildBattleTaskData)
	}

	task_data := &GuildBattleTaskData{guildMarchEntityID: 0}
	return task_data
}

func (t *GuildBattleTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata := t.getTaskData(runtask, taskindex)
	if taskdata.guildMarchEntityID != 0 {
		if player.CheckGuildEntityExist(taskdata.guildMarchEntityID) {
			// 还在集结中
			return true
		} else {
			taskdata.guildMarchEntityID = 0
		}
	}

	target_id, _, ok := player.SearchTarget(t.search_type, 1)
	if !ok {
		common.GStdout.Debug("GuildBattleTask 1 ")
		return false
	}

	army_data := player.RandomArmyDataByRange(t.minArmy, t.maxArmy)
	if army_data == nil {
		common.GStdout.Debug("GuildBattleTask 2 ")
		return false
	}

	command_data := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, target_id, 0, 0)
	player.PM("addap 1000")
	marchEntityID, ok := player.CreateGuildMarch(army_data, command_data, t.wait_time_index)
	if !ok {
		common.GStdout.Debug("GuildBattleTask 3 ")
		return false
	}

	taskdata.guildMarchEntityID = marchEntityID

	if t.waitMarchDisband {
		player.WaitGuildEntityDisband(marchEntityID, 0)
	}
	return true
}

// GuildBattleToTargetTask comment
type GuildBattleToTargetTask struct {
	target_id        uint64
	wait_time_index  uint32
	waitMarchDisband bool
	minArmy          uint32
	maxArmy          uint32
}

func (t *GuildBattleToTargetTask) getTaskData(runtask *BTree.BTreeRunTask, taskindex uint32) *GuildBattleTaskData {
	data := runtask.GetTaskData(taskindex)
	if data != nil {
		return data.(*GuildBattleTaskData)
	}

	task_data := &GuildBattleTaskData{guildMarchEntityID: 0}
	return task_data
}

func CreateGuildBattleToTargetTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("create Guild Battle params error!")
	}

	target_id, error := strconv.ParseUint(params[0], 10, 64)
	if error != nil {
		return nil, error
	}

	wait_time_index, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	waitMarchDisband, error := strconv.ParseBool(params[2])
	if error != nil {
		return nil, error
	}

	minArmy, error := strconv.Atoi(params[3])
	if error != nil {
		return nil, error
	}

	maxArmy, error := strconv.Atoi(params[4])
	if error != nil {
		return nil, error
	}

	task := &GuildBattleToTargetTask{target_id: target_id, wait_time_index: (uint32)(wait_time_index), waitMarchDisband: waitMarchDisband, minArmy: (uint32)(minArmy), maxArmy: (uint32)(maxArmy)}
	return task, nil
}

func (t *GuildBattleToTargetTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata := t.getTaskData(runtask, taskindex)

	if taskdata.guildMarchEntityID != 0 {
		if player.CheckGuildEntityExist(taskdata.guildMarchEntityID) {
			// 还在集结中
			return true
		} else {
			taskdata.guildMarchEntityID = 0
		}
	}

	army_data := player.RandomArmyDataByRange(t.minArmy, t.maxArmy)
	if army_data == nil {
		return false
	}

	commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, t.target_id, 0, 0)
	player.PM("addap 1000")
	marchEntityID, ok := player.CreateGuildMarch(army_data, commandData, t.wait_time_index)
	if !ok {
		return false
	}

	taskdata.guildMarchEntityID = marchEntityID
	if t.waitMarchDisband {
		player.WaitGuildEntityDisband(marchEntityID, 0)
	}
	return true
}

// JoinGuildMarchTask comment
type JoinGuildMarchTask struct {
	targetType       uint32
	waitMarchDisband bool
	minArmy          uint32
	maxArmy          uint32
}

// JoinGuildMarchTaskData comment
type JoinGuildMarchTaskData struct {
	marchEntityID      uint64
	guildMarchEntityID uint64
}

func (t *JoinGuildMarchTask) getTaskData(runtask *BTree.BTreeRunTask, taskindex uint32) *JoinGuildMarchTaskData {
	data := runtask.GetTaskData(taskindex)
	if data != nil {
		return data.(*JoinGuildMarchTaskData)
	}

	task_data := &JoinGuildMarchTaskData{marchEntityID: 0, guildMarchEntityID: 0}
	return task_data
}

func CreateJoinGuildMarchTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 4 {
		return nil, fmt.Errorf("create Join Guild March params error!")
	}

	targetType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	waitMarchDisband, error := strconv.ParseBool(params[1])
	if error != nil {
		return nil, error
	}

	minArmy, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	maxArmy, error := strconv.Atoi(params[3])
	if error != nil {
		return nil, error
	}

	task := &JoinGuildMarchTask{targetType: (uint32)(targetType), waitMarchDisband: waitMarchDisband, minArmy: (uint32)(minArmy), maxArmy: (uint32)(maxArmy)}
	return task, nil
}

// DoTask comment
func (t *JoinGuildMarchTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata := t.getTaskData(runtask, taskindex)
	if taskdata.marchEntityID != 0 {
		if player.CheckEntityExist(taskdata.marchEntityID) {
			// 还在集结中
			return true
		} else {
			taskdata.marchEntityID = 0
			taskdata.guildMarchEntityID = 0
		}
	}

	armyData := player.RandomArmyDataByRange(t.minArmy, t.maxArmy)
	if armyData == nil {
		return false
	}

	player.PM("addap 1000")
	var marchIndex uint32

	player.RangeGuildMarchs(func(entityID uint64, entityData *data.SyncEntityData) bool {
		if t.targetType != 0 {
			if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_GuildMarch)); ok {
				if property.(*protomsg.GuildMarchData).TargetType != t.targetType {
					return true
				}
			}
		}

		if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner)); ok {
			if property.(*protomsg.EntityOwnerData).PlayerId == player.PlayerID {
				return true
			}
		}

		commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_JoinMarch, entityID, 0, 0)
		newMarchIndex, createEntityID, ok := player.CreateMarch(armyData, commandData)
		if ok {
			marchIndex = newMarchIndex
			taskdata.marchEntityID = createEntityID
			taskdata.guildMarchEntityID = entityID
			return false
		}
		return true
	})

	if marchIndex == 0 {
		return false
	}

	if t.waitMarchDisband {
		player.WaitMarchBack(taskdata.marchEntityID, 0)
	}
	return false
}

// SendGuildMailTask comment
type SendGuildMailTask struct {
}

func CreateSendGuildMailTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &SendGuildMailTask{}
	return task, nil
}

// DoTask comment
func (t *SendGuildMailTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	request := &protomsg.MsgCL2GSGuildSendMailRequest{}
	request.Subject = RandStringBytes(10)
	request.Content = RandStringBytes(20)
	if response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildSendMailReply).(*protomsg.MsgGS2CLGuildSendMailReply); ok {
		common.GStdout.Debug("send guild mail result %v", response.ErrorCode)
		return response.ErrorCode == 0
	}
	return false
}

// CreateMoveCastleToGuildTask comment
func CreateMoveCastleToGuildTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("move castle to guild task params error")
	}

	CenterType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	SearchRange, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	task := &MoveCastleToGuildTask{CenterType: (uint32)(CenterType), SearchRange: (uint32)(SearchRange)}
	return task, nil
}

// MoveCastleToGuildTask 采集
type MoveCastleToGuildTask struct {
	CenterType  uint32 // 1 盟主主堡
	SearchRange uint32
}

func (t *MoveCastleToGuildTask) FindCenterPos(runtask *BTree.BTreeRunTask, taskindex uint32) (*protomsg.Vector2D, bool) {
	player := runtask.TaskObj.(*Player)
	guildId := player.GetGuildId()
	if guildId == 0 {
		common.GStdout.Debug("FindCenterPos 1")
		return nil, false
	}

	switch t.CenterType {
	case 1:
		request := &protomsg.MsgCL2GSGuildQueryRequest{}
		request.GuildId = guildId
		common.GStdout.Debug("FindCenterPos send")
		if response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGuildQuserReply).(*protomsg.MsgGS2CLGuildQuserReply); ok {
			var leaderId uint64
			for _, v := range response.Members {
				if v.Position == 5 {
					leaderId = v.PlayerId
					break
				}
			}

			if leaderId == 0 {
				common.GStdout.Debug("FindCenterPos 3")
				return nil, false
			}

			var castleEntity *data.SyncEntityData
			player.RangeGuildEntities(func(entityID uint64, entityData *data.SyncEntityData) bool {
				if entityData.EntityType == (uint32)(protomsg.EntityType_kEntityType_Castle) {
					if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner)); ok {
						if property.(*protomsg.EntityOwnerData).PlayerId == leaderId {
							castleEntity = entityData
							return false
						}
					}
				}
				return true
			})

			if castleEntity == nil {
				common.GStdout.Debug("FindCenterPos 5")
				return nil, false
			}

			if property, ok := castleEntity.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Map)); ok {
				return property.(*protomsg.MapData).Position, true
			}
		}
	}
	common.GStdout.Debug("FindCenterPos 6")
	return nil, false
}

// DoTask comment
func (t *MoveCastleToGuildTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	if centerPos, ok := t.FindCenterPos(runtask, taskindex); ok {
		request := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
		request.EntityType = protomsg.EntityType_kEntityType_Castle
		request.SearchLevel = 1
		request.SearchRange = t.SearchRange
		request.CenterPos = &protomsg.Vector2D{}
		request.CenterPos.X = centerPos.X
		request.CenterPos.Y = centerPos.Y

		response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSearchEmptyPosReply).(*protomsg.MsgGS2CLSearchEmptyPosReply)
		if !ok {
			return false
		}

		player.PM("movecastle 99 " + strconv.FormatInt(response.Position.X/6000, 10) + " " + strconv.FormatInt(response.Position.Y/6000, 10))
		return true
	}
	common.GStdout.Debug("FindCenterPos 7")
	return false
}

// GuildAssistHelpTask comment
type GuildAssistHelpTask struct {
}

func CreateGuildAssistHelpTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &GuildAssistHelpTask{}
	return task, nil
}

// DoTask comment
func (t *GuildAssistHelpTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	request := &protomsg.MsgCL2GSGuildAssistHelpOtherRequest{}
	player.Send(request)
	return true
}
