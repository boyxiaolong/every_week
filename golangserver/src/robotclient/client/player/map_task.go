package player

import (
	"fmt"
	"math/rand"
	BTree "public/behaviortree"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
	"robotclient/client/player/data"
	"strconv"
	"time"
)

var mapSize int64
var mapMaxPos int64

// Init comment
func init() {
	BTree.RegisterTaskCreator("ViewMap", CreateViewMapTask)
	BTree.RegisterTaskCreator("MapSearch", CreateMapSearchTask)
	BTree.RegisterTaskCreator("Battle", CreateBattleTask)
	BTree.RegisterTaskCreator("Collect", CreateCollectTask)
	BTree.RegisterTaskCreator("MarchMove", CreateMarchMoveTask)
	BTree.RegisterTaskCreator("MarchMoveBack", CreateMarchMoveBackTask)
	BTree.RegisterTaskCreator("QuarterMarch", CreateQuarterMarchTask)
	BTree.RegisterTaskCreator("BattleToTarget", CreateBattleToTargetTask)
	BTree.RegisterTaskCreator("MoveCastle", CreateMoveCastleTask)
	BTree.RegisterTaskCreator("Scout", CreateScoutTask)
	BTree.RegisterTaskCreator("ScoutTarget", CreateScoutTargetTask)
	BTree.RegisterTaskCreator("ExploreMist", CreateExploreMistTask)
	BTree.RegisterTaskCreator("Reenforce", CreateReenforceTask) // 驻防
	BTree.RegisterTaskCreator("ReenforceSelf", CreateReenforceSelfTask)

	BTree.RegisterConditionCallback("CheckMarchCount", CheckMarchCount)
	BTree.RegisterConditionCallback("CheckMapInit", CheckMapInit)

	mapSize = 1200
	mapMaxPos = (mapSize - 1) * 6 * 1000
}

func (p *Player) SearchTarget(searchType protomsg.MapSearchType, searchLevel uint32) (uint64, *protomsg.Vector2D, bool) {

	// cmd := fmt.Sprintf("mapsearch %v %v 4", (uint32)(searchType), searchLevel)
	// if searchType == protomsg.MapSearchType_kMapSearchType_Castle || searchType == protomsg.MapSearchType_kMapSearchType_March {
	// 	cmd = fmt.Sprintf("mapsearch %v %v 0", (uint32)(searchType), searchLevel)
	// }
	// request := &protomsg.MsgCL2GSPMCommandRequest{}
	// request.Command = cmd

	request := &protomsg.MsgCL2GSMapSearchRequest{}
	request.SearchType = searchType
	request.SearchLevel = searchLevel
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMapSearchReply).(*protomsg.MsgGS2CLMapSearchReply)
	if ok {
		return response.EntityId, response.Position, int32(response.ErrorCode) == 0
	}

	// if searchType == protomsg.MapSearchType_kMapSearchType_Castle || searchType == protomsg.MapSearchType_kMapSearchType_March {
	// 	return 0, nil, false
	// }

	// cmd = fmt.Sprintf("mapsearch %v %v 2", (uint32)(searchType), searchLevel)
	// request.Command = cmd
	// response, ok = p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMapSearchReply).(*protomsg.MsgGS2CLMapSearchReply)
	// if ok {
	// 	return response.EntityId, response.Position, int32(response.ErrorCode) == 0
	// }
	return 0, nil, false
}

func (p *Player) GetMainRegionId() uint64 {
	regionDataCenter := p.GetData("RegionDataCenter").(*data.RegionDataCenter)
	return regionDataCenter.MainRegionId
}

func (p *Player) GetMainRegion() (*data.RegionData, bool) {
	regionDataCenter := p.GetData("RegionDataCenter").(*data.RegionDataCenter)
	return regionDataCenter.GetMainRegion()
}

func (p *Player) GetRegion(regionId uint64) (*data.RegionData, bool) {
	regionDataCenter := p.GetData("RegionDataCenter").(*data.RegionDataCenter)
	return regionDataCenter.GetRegion(regionId)
}

func (p *Player) MaxTroopCapacity() uint32 {
	tccount := (uint32)(len(data.TroopCapacity))
	if tccount == 0 {
		return 0
	}

	cityData := p.GetData("CityData").(*data.CityData)
	castleLevel := cityData.GetBuildMaxLevelByType(1000)
	if castleLevel > 0 {
		castleLevel--
	}

	if castleLevel >= tccount {
		castleLevel = tccount - 1
	}

	return data.TroopCapacity[castleLevel]
}

func (p *Player) RandomArmyData() *protomsg.ArmyData {
	return p.RandomArmyDataByRange(1000, 1000)
}

func (p *Player) RandomArmyDataByRange(minArmyRate uint32, maxArmyRate uint32) *protomsg.ArmyData {
	var hero1 uint32
	var hero2 uint32
	regionData, ok := p.GetMainRegion()
	if !ok {
		return nil
	}

	usedheros := regionData.UsedHeros

	validHeros := make([]uint32, 0)
	for _, v := range data.ValidHeros {
		if _, ok := usedheros.Load(v); !ok {
			validHeros = append(validHeros, v)
		}
	}

	if len(validHeros) == 0 {
		return nil
	}

	for k, v := range validHeros {
		randomIndex := rand.Intn(len(validHeros))
		t := validHeros[randomIndex]
		validHeros[randomIndex] = v
		validHeros[k] = t
	}

	hero1 = validHeros[0]
	if len(validHeros) > 1 {
		hero2 = validHeros[1]
	}

	if len(data.ValidTroops) == 0 {
		return nil
	}

	rank_troops := data.ValidTroops[rand.Intn(len(data.ValidTroops))]
	if len(rank_troops) == 0 {
		return nil
	}

	troopid := rank_troops[rand.Intn(len(rank_troops))]

	if minArmyRate == maxArmyRate && maxArmyRate == 0 {
		minArmyRate = 1000
		maxArmyRate = 1000
	}

	if minArmyRate > maxArmyRate {
		maxArmyRate, minArmyRate = minArmyRate, maxArmyRate
	}
	randomRange := uint32(maxArmyRate - minArmyRate)

	troop_uints := p.MaxTroopCapacity() * (minArmyRate + randomRange) / 1000

	//common.GStdout.Info("MarchArmy range %v,%v,%v,%v,%v",maxArmyRate,minArmyRate, randomRange,p.MaxTroopCapacity(),troop_uints)
	//troop_uints = 20000

	pm := fmt.Sprintf("addhero %v", hero1)
	p.PM(pm)

	pm = fmt.Sprintf("upgradehero %v %v", hero1, rand.Intn(100000))
	p.PM(pm)

	if hero2 != 0 {
		pm = fmt.Sprintf("addhero %v", hero2)
		p.PM(pm)
		pm = fmt.Sprintf("upgradehero %v %v", hero2, rand.Intn(100000))
	}

	pm = fmt.Sprintf("addarmy %v %v", troopid, troop_uints)
	p.PM(pm)
	time.Sleep(1 * time.Second)
	return MakeArmyData(hero1, hero2, troopid, troop_uints)
}

// WaitMarchBack a
func (p *Player) WaitMarchBack(marchEntityID uint64, waitTime int) bool {
	if waitTime == 0 {
		for true {
			time.Sleep(10 * time.Second)
			if p.CheckEntityExist(marchEntityID) {
				continue
			}
			return true
		}
	}

	for i := 0; i < waitTime; i++ {
		time.Sleep(1 * time.Second)
		if p.CheckEntityExist(marchEntityID) {
			continue
		}
		return true
	}
	return false
}

// WaitScoutBack a
func (p *Player) WaitScoutBack(scoutIndex uint32, waitTime int) bool {
	if waitTime == 0 {
		for true {
			time.Sleep(10 * time.Second)
			if p.CheckScoutExist(scoutIndex) {
				continue
			}
			return true
		}
	}

	for i := 0; i < waitTime; i++ {
		time.Sleep(1 * time.Second)
		if p.CheckScoutExist(scoutIndex) {
			continue
		}
		return true
	}
	return false
}

// CheckEntityExist CheckEntityExist
func (p *Player) CheckEntityExist(marchEntityID uint64) bool {
	if mainRegion, ok := p.GetMainRegion(); ok {
		_, ok := mainRegion.GetPlayerCenter().GetEntity(marchEntityID)
		return ok
	}
	return false
}

// CheckEntityExist CheckEntityExist
func (p *Player) GetEntity(entityID uint64) (*data.SyncEntityData, bool) {
	if mainRegion, ok := p.GetMainRegion(); ok {
		entity, ok := mainRegion.GetPlayerCenter().GetEntity(entityID)
		return entity, ok
	}
	return nil, false
}

// CheckEntityExist CheckEntityExist
func (p *Player) GetViewEntity(entityID uint64) (*data.SyncEntityData, bool) {
	if mainRegion, ok := p.GetMainRegion(); ok {
		entity, ok := mainRegion.GetViewCenter().GetEntity(entityID)
		return entity, ok
	}
	return nil, false
}

// GetMarchCount  通知地图数据
func (p *Player) GetMarchCount() uint32 {
	if mainRegion, ok := p.GetMainRegion(); ok {
		return mainRegion.GetPlayerCenter().GetEntityCount((uint32)(protomsg.EntityType_kEntityType_March))
	}
	return 0
}

// RangeEntities  通知地图数据
func (p *Player) RangeEntities(fn func(entityID uint64, entityData *data.SyncEntityData) bool) {
	if mainRegion, ok := p.GetMainRegion(); ok {
		mainRegion.GetPlayerCenter().RangeEntities(fn)
	}
}

// CheckMarchCount CheckMarchCount
func CheckMarchCount(testobj BTree.ObjectInterface, params string) bool {
	player := testobj.(*Player)
	checkCount, error := strconv.Atoi(params)
	if error != nil {
		return false
	}

	return player.GetMarchCount() >= (uint32)(checkCount)
}

func CheckMapInit(testobj BTree.ObjectInterface, params string) bool {
	player := testobj.(*Player)
	if mainRegion, ok := player.GetMainRegion(); ok {
		return mainRegion.GetPlayerCenter().GetEntityCount((uint32)(protomsg.EntityType_kEntityType_Castle)) > 0
	}
	return false
}

func MakeArmyData(hero1 uint32, hero2 uint32, troop_id uint32, troop_uints uint32) *protomsg.ArmyData {
	army := &protomsg.ArmyData{}
	army.Hero1 = hero1
	army.Hero2 = hero2
	army.CurrentTroops = make([]*protomsg.TroopData, 0)

	troop_data := &protomsg.TroopData{}
	troop_data.TroopId = troop_id
	troop_data.TroopUints = troop_uints
	army.CurrentTroops = append(army.CurrentTroops, troop_data)
	return army
}

func MakeCommand(target_type protomsg.MarchCommandTarget, target_id uint64, posx int64, posy int64) *protomsg.MarchCommand {
	command := &protomsg.MarchCommand{}
	command.TargetType = target_type
	command.TargetId = target_id
	command.Position = &protomsg.Vector2D{}
	command.Position.X = posx
	command.Position.Y = posy
	return command
}

func (p *Player) CreateMarch(army *protomsg.ArmyData, command *protomsg.MarchCommand) (uint32, uint64, bool) {
	request := &protomsg.MsgCL2GSCreateMarchRequest{}
	request.Command = command
	request.Army = army
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateMarchReply).(*protomsg.MsgGS2CLCreateMarchReply)
	if !ok {
		return 0, 0, false
	}
	common.GStdout.Info("=================== CreateMarch response %v", response.ErrorCode)
	time.Sleep(1 * time.Second) // 等待数据同步
	if int32(response.ErrorCode) == 0 {

		return response.March.MarchIndex, response.March.EntityId, true
	}
	return 0, 0, false
}

func (p *Player) CreateTkMarch(deploy_id uint32, command *protomsg.MarchCommand) (uint32, uint64, bool) {
	request := &protomsg.MsgCL2GSCreateMarchTKRequest{}
	request.Command = command
	request.DeployId = deploy_id
	request.RegionId = p.GetMainRegionId()
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCreateMarchReply).(*protomsg.MsgGS2CLCreateMarchReply)
	if !ok {
		return 0, 0, false
	}

	common.GStdout.Info("=================== CreateMarch response %v", response.ErrorCode)
	time.Sleep(1 * time.Second) // 等待数据同步
	if int32(response.ErrorCode) == 0 {
		return response.March.MarchIndex, response.March.EntityId, true
	}
	return 0, 0, false
}

func (p *Player) MarchCommand(marchIndex uint32, command *protomsg.MarchCommand) bool {
	request := &protomsg.MsgCL2GSMarchCommandRequest{}
	request.Command = command
	request.MarchIndex = marchIndex
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMarchCommandReply).(*protomsg.MsgGS2CLMarchCommandReply)
	if !ok {
		common.GStdout.Debug(" march move 3 ")
		return false
	}

	common.GStdout.Info("=================== MarchCommand response %v", response.ErrorCode)
	return int32(response.ErrorCode) == 0
}

// CheckScoutExist CheckScoutExist
func (p *Player) CheckScoutExist(scoutIndex uint32) bool {
	isExist := false
	p.RangeEntities(func(entityID uint64, entityData *data.SyncEntityData) bool {
		if entityData.EntityType == (uint32)(protomsg.EntityType_kEntityType_Scout) {
			if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Scout)); ok {
				if property.(*protomsg.EntityScoutData).ScoutIndex == scoutIndex {
					isExist = true
				}
				return !isExist
			}
		}
		return true
	})
	return isExist
}

// GetCastleEntity GetCastleEntity
func (p *Player) GetCastleEntity() (*data.SyncEntityData, bool) {
	var castleEntity *data.SyncEntityData
	p.RangeEntities(func(entityID uint64, entityData *data.SyncEntityData) bool {
		if entityData.EntityType == (uint32)(protomsg.EntityType_kEntityType_Castle) {
			castleEntity = entityData
			return false
		}
		return true
	})
	return castleEntity, castleEntity != nil
}

// GetCastleEntity GetCastleEntity
func (p *Player) GetCastlePos() *protomsg.Vector2D {
	var castleEntity *data.SyncEntityData
	p.RangeEntities(func(entityID uint64, entityData *data.SyncEntityData) bool {
		if entityData.EntityType == (uint32)(protomsg.EntityType_kEntityType_Castle) {
			castleEntity = entityData
			return false
		}
		return true
	})
	if castleEntity == nil {
		return &protomsg.Vector2D{}
	}
	return castleEntity.GetPosition()
}

func CreateViewMapTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("view map task params error!")
	}

	viewType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	rangeType, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	moveRange, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	srcPosition := &protomsg.Vector2D{}
	switch viewType {
	case 2:
		if len(params) < 5 {
			return nil, fmt.Errorf("view map task params error!")
		}

		x, error := strconv.ParseInt(params[3], 10, 64)
		if error != nil {
			return nil, error
		}

		y, error := strconv.ParseInt(params[4], 10, 64)
		if error != nil {
			return nil, error
		}
		srcPosition.X = x * 1000 * 6
		srcPosition.Y = y * 1000 * 6
	case 3:
		srcPosition.X = rand.Int63n(mapMaxPos)
		srcPosition.Y = rand.Int63n(mapMaxPos)
	}

	task := &ViewMapTask{ViewType: (uint32)(viewType), RangeType: (uint32)(rangeType), MoveRange: ((int64)(moveRange)) * 6 * 1000, SrcPosition: srcPosition}
	return task, nil
}

// ViewMapTask comment
type ViewMapTask struct {
	ViewType     uint32 // 1 主堡为中心 2 指定坐标 3 随机中心坐标
	RangeType    uint32 // 1 以中心随机范围内查看 2 已之前的查看位置按随机范围内查看
	MoveRange    int64
	SrcPosition  *protomsg.Vector2D
	ViewPosition *protomsg.Vector2D
}

// ViewMapTaskData comment
type ViewMapTaskData struct {
	SrcPosition  *protomsg.Vector2D
	ViewPosition *protomsg.Vector2D
}

func (t *ViewMapTask) getTaskData(runtask *BTree.BTreeRunTask, taskindex uint32) *ViewMapTaskData {
	data := runtask.GetTaskData(taskindex)
	if data != nil {
		return data.(*ViewMapTaskData)

	}

	task_data := &ViewMapTaskData{SrcPosition: &protomsg.Vector2D{}, ViewPosition: &protomsg.Vector2D{}}

	switch t.ViewType {
	case 1:
		player := runtask.TaskObj.(*Player)
		castlePos := player.GetCastlePos()
		common.GStdout.Debug("get castle pos(%v,%v)", castlePos.X, castlePos.Y)
		task_data.SrcPosition.X = castlePos.X
		task_data.SrcPosition.Y = castlePos.Y
		task_data.ViewPosition = &protomsg.Vector2D{castlePos.X, castlePos.Y}
	case 2:
		task_data.SrcPosition.X = t.SrcPosition.X
		task_data.SrcPosition.Y = t.SrcPosition.Y
	case 3:
		task_data.SrcPosition.X = rand.Int63n(mapMaxPos)
		task_data.SrcPosition.Y = rand.Int63n(mapMaxPos)
	}

	task_data.ViewPosition.X = task_data.SrcPosition.X
	task_data.ViewPosition.Y = task_data.SrcPosition.Y
	runtask.SetTaskData(taskindex, task_data)

	return task_data
}

func (t *ViewMapTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata := t.getTaskData(runtask, taskindex)
	if t.ViewPosition == nil {
	}

	switch t.RangeType {
	case 1:
		taskdata.ViewPosition.X = taskdata.SrcPosition.X + rand.Int63n(t.MoveRange*2) - t.MoveRange
		taskdata.ViewPosition.Y = taskdata.SrcPosition.Y + rand.Int63n(t.MoveRange*2) - t.MoveRange
	case 2:
		taskdata.ViewPosition.X = taskdata.ViewPosition.X + rand.Int63n(t.MoveRange*2) - t.MoveRange
		taskdata.ViewPosition.Y = taskdata.ViewPosition.Y + rand.Int63n(t.MoveRange*2) - t.MoveRange
	}

	if taskdata.ViewPosition.X < 0 {
		taskdata.ViewPosition.X = 0
	}

	if taskdata.ViewPosition.X > mapMaxPos {
		taskdata.ViewPosition.X = mapMaxPos
	}

	if taskdata.ViewPosition.Y < 0 {
		taskdata.ViewPosition.Y = 0
	}

	if taskdata.ViewPosition.Y > mapMaxPos {
		taskdata.ViewPosition.Y = mapMaxPos
	}

	common.GStdout.Debug("view map type (%v,%v,%v) src(%v,%v) view(%v,%v)", t.ViewType, t.RangeType, t.MoveRange, taskdata.SrcPosition.X, taskdata.SrcPosition.Y, taskdata.ViewPosition.X, taskdata.ViewPosition.Y)

	request := &protomsg.MsgCL2GSViewMapRequest{}
	request.KingdomId = 0
	request.ViewLod = 0
	request.Position = taskdata.ViewPosition

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLViewMapReply).(*protomsg.MsgGS2CLViewMapReply)
	if !ok {
		return false
	}

	return int32(response.ErrorCode) == 0
}

func CreateMapSearchTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &MapSearchTask{}
	return task, nil
}

// MapSearchTask comment
type MapSearchTask struct {
}

func (t *MapSearchTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	rand.Seed(time.Now().Unix())

	request := &protomsg.MsgCL2GSMapSearchRequest{}
	request.SearchType = (protomsg.MapSearchType)(rand.Int31n(5))
	request.SearchLevel = 1

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMapSearchReply).(*protomsg.MsgGS2CLMapSearchReply)
	if !ok {
		return false
	}

	return int32(response.ErrorCode) == 0
}

type ReenforceSelfTask struct {
}

type ReenforceTask struct {
	searchType    protomsg.MapSearchType
	searchLevel   uint32
	waitMarchBack bool
	viewBattle    bool
	isQuarter     int32
	isCapture     bool
}

// BattleTask comment
type BattleTask struct {
	searchType    protomsg.MapSearchType
	searchLevel   uint32
	waitMarchBack bool
	viewBattle    bool
	isQuarter     int32
	minArmy       uint32
	maxArmy       uint32
}

// BattleTaskData comment
type BattleTaskData struct {
	marchEntityID uint64
	marchIndex    uint32
	targetID      uint64
	skillGroupId  uint32
}

// ReenforceTaskData comment
type ReenforceTaskData struct {
	marchEntityID uint64
	marchIndex    uint32
	targetID      uint64
}

func (player *Player) CastSkill(marchIndex uint32, skillGroupId uint32, targetID uint64) bool {
	entityRequest := &protomsg.MsgCL2GSGetEntityBaseDataRequest{}
	entityRequest.EntityId = targetID
	entityRequest.RegionId = player.GetMainRegionId()
	entityResponse, ok := player.SendAndWait(entityRequest, msgtype.MsgType_kMsgGS2CLEntityBaseDataReply).(*protomsg.MsgGS2CLEntityBaseDataReply)
	if !ok || entityResponse.ErrorCode != 0 {
		common.GStdout.Debug("cast skill error 1")
		return false
	}

	request := &protomsg.MsgCL2GSSkillCastRequest{}
	request.MarchIndex = marchIndex
	request.SkillGroupId = skillGroupId
	request.TargetEntityId = targetID
	request.TarPos = entityResponse.EntityData.Position
	request.RegionId = player.GetMainRegionId()

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSkillCastReply).(*protomsg.MsgGS2CLSkillCastReply)
	if !ok || response.ErrorCode != 0 {
		return false
	}
	return true
}

func (player *Player) CheckTarget(targetID uint64) bool {
	entityRequest := &protomsg.MsgCL2GSGetEntityBaseDataRequest{}
	entityRequest.EntityId = targetID
	entityRequest.RegionId = player.GetMainRegionId()
	entityResponse, ok := player.SendAndWait(entityRequest, msgtype.MsgType_kMsgGS2CLEntityBaseDataReply).(*protomsg.MsgGS2CLEntityBaseDataReply)
	if !ok || entityResponse.ErrorCode != 0 {
		return false
	}
	return true
}

// DoTask comment
func (t *BattleTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata, ok := runtask.GetTaskData(taskindex).(*BattleTaskData)
	if !ok {
		taskdata = &BattleTaskData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	newCommand := false
	newMarch := false
	if taskdata.marchEntityID == 0 {
		newCommand = true
		newMarch = true
	} else if t.isQuarter != 0 {
		if !player.CheckEntityExist(taskdata.marchEntityID) {
			common.GStdout.Debug("BattleTask %v not exist", taskdata.marchEntityID)
			newCommand = true
			newMarch = true
		} else {
			common.GStdout.Debug("BattleTask %v exist", taskdata.marchEntityID)
			if entityData, ok := player.GetEntity(taskdata.marchEntityID); ok {
				if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_MarchCommand)); ok {
					marchCommand := property.(*protomsg.MarchCommand)
					if marchCommand.TargetType == protomsg.MarchCommandTarget_kMarchCommandTarget_Station {
						newCommand = true
					}
				}
			} else {
				newCommand = true
				newMarch = true
			}
		}
	}

	var commandData *protomsg.MarchCommand
	if newCommand {
		var targetID uint64
		if taskdata.targetID != 0 && player.CheckTarget(taskdata.targetID) {
			targetID = taskdata.targetID
		} else {
			targetID, _, ok = player.SearchTarget(t.searchType, t.searchLevel)
			if !ok {
				return false
			}
		}

		player.PM("addap 1000")
		commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, targetID, 0, 0)
		commandData.IsQuarter = t.isQuarter
		taskdata.targetID = targetID
	}

	if newMarch {
		taskdata.marchIndex = 0
		taskdata.marchEntityID = 0
		taskdata.targetID = 0
		taskdata.skillGroupId = 0

		armyData := player.RandomArmyDataByRange(t.minArmy, t.maxArmy)
		if armyData == nil {
			return false
		}

		marchIndex, marchEntityID, ok := player.CreateMarch(armyData, commandData)
		if !ok {
			return false
		}

		taskdata.marchEntityID = marchEntityID
		taskdata.marchIndex = marchIndex
		taskdata.targetID = commandData.TargetId

		if skillGroupId, ok := data.HeroSkillGroup[armyData.Hero1]; ok {
			taskdata.skillGroupId = skillGroupId
		}

		if t.viewBattle {
			request := &protomsg.MsgCL2GSViewMapRequest{}
			request.KingdomId = 0
			request.ViewLod = 0
			request.Position = commandData.Position
			player.Send(request)
		}
	} else if newCommand {
		ok = player.MarchCommand(taskdata.marchIndex, commandData)
		if !ok {
			return false
		}
	}

	if taskdata.skillGroupId != 0 {
		if entityData, ok := player.GetEntity(taskdata.marchEntityID); ok {
			if entityproperty, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus)); ok {
				battleStatus := entityproperty.(*protomsg.BattleStatusData)
				if battleStatus.IsBattle > 0 {
					player.CastSkill(taskdata.marchIndex, taskdata.skillGroupId, taskdata.targetID)
				}
			}
		}
	}

	if t.waitMarchBack {
		player.WaitMarchBack(taskdata.marchEntityID, 0)
	}
	return true
}

func CreateReenforceSelfTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &ReenforceSelfTask{}
	return task, nil
}

func (t *ReenforceSelfTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata, ok := runtask.GetTaskData(taskindex).(*BattleTaskData)
	if !ok {
		taskdata = &BattleTaskData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	deply_id := player.GetDeployInfoId()
	if deply_id == nil {
		return false
	}

	data, ok := player.GetCastleEntity()
	if !ok {
		return false
	}

	commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_CastleSelfJoinReenforce, data.EntityId, data.GetPosition().X, data.GetPosition().Y)
	marchIndex, marchEntityID, ok := player.CreateTkMarch(*deply_id, commandData)
	if !ok {
		return false
	}

	taskdata.marchEntityID = marchEntityID
	taskdata.marchIndex = marchIndex
	return true
}

func CreateReenforceTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 6 {
		return nil, fmt.Errorf("reenforce task params error %v (%v)", len(params), params)
	}
	searchType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	searchLevel, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	if searchLevel == 0 {
		searchLevel = 1
	}

	waitMarchBack, error := strconv.ParseBool(params[2])
	if error != nil {
		return nil, error
	}

	viewBattle, error := strconv.ParseBool(params[3])
	if error != nil {
		return nil, error
	}

	isQuarter, error := strconv.Atoi(params[4])
	if error != nil {
		return nil, error
	}

	isCapture, error := strconv.ParseBool(params[5])
	if error != nil {
		return nil, error
	}

	task := &ReenforceTask{searchType: (protomsg.MapSearchType)(searchType), searchLevel: (uint32)(searchLevel), waitMarchBack: waitMarchBack, viewBattle: viewBattle, isQuarter: (int32)(isQuarter), isCapture: isCapture}
	return task, nil
}

func (t *ReenforceTask) DoCastle(player *Player, owernData *protomsg.EntityOwnerData, targetID uint64, position *protomsg.Vector2D) (*protomsg.MarchCommand, bool) {
	var commandData *protomsg.MarchCommand
	if owernData.PlayerId != 0 {
		if owernData.PlayerId == player.PlayerID {
			// 自己的城堡
			return nil, true
		}

		// 搜索的城堡不存在公会
		if owernData.GuildId == 0 {
			return nil, true
		}

		if player.GetGuildId() != 0 {
			// 玩家存在公会
			if player.GetGuildId() != owernData.GuildId {
				return nil, true
			}
		} else {
			// 加入公会，驻防
			if !player.JoinGuild(owernData.GuildId) {
				// 加入公会失败
				return nil, false
			}
		}

		// 加入友军城堡驻防
		commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_CastleFriendJoinReenforce, targetID, position.X, position.Y)
	} else {
		return nil, true
	}
	return commandData, true
}

func (t *ReenforceTask) DoStronghold(player *Player, owernData *protomsg.EntityOwnerData, targetID uint64, position *protomsg.Vector2D) (*protomsg.MarchCommand, bool) {
	var commandData *protomsg.MarchCommand
	if owernData.PlayerId != 0 {
		if owernData.PlayerId == player.PlayerID {
			// 已经被自己占领了
			return nil, true
		}

		// 占领的玩家不存在公会
		if owernData.GuildId == 0 {
			return nil, true
		}

		// 占领的玩家存在公会
		// 玩家的公会和占领玩家公会不一致
		if player.GetGuildId() != 0 {
			// 玩家存在公会
			if player.GetGuildId() != owernData.GuildId {
				return nil, true
			}
		} else {
			// 加入公会，驻防
			if !player.JoinGuild(owernData.GuildId) {
				// 加入公会失败
				return nil, false
			}
		}

		// 加入据点驻防
		commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_JoinReenforce, targetID, position.X, position.Y)
		// 公会占领
	} else if owernData.GuildId != 0 {
		if player.GetGuildId() != owernData.GuildId {
			return nil, true
		}

		if player.GetGuildId() != 0 {
			return nil, true
		}

		// 加入公会，驻防
		if !player.JoinGuild(owernData.GuildId) {
			// 加入公会失败
			return nil, false
		}
		// 加入据点驻防
		commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_JoinReenforce, targetID, position.X, position.Y)
	} else {
		// 需要攻打下据点，才能驻防
		// 攻打据点
		if t.isCapture {
			player.PM("capturestronghold " + strconv.Itoa(int(targetID)))
			commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_JoinReenforce, targetID, position.X, position.Y)
		} else {
			commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, targetID, position.X, position.Y)
		}
	}
	return commandData, true
}

func (t *ReenforceTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata, ok := runtask.GetTaskData(taskindex).(*ReenforceTaskData)
	if !ok {
		taskdata = &ReenforceTaskData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	newCommand := false
	newMarch := false
	if taskdata.marchEntityID == 0 {
		newCommand = true
		newMarch = true
	} else if t.isQuarter != 0 {
		if !player.CheckEntityExist(taskdata.marchEntityID) {
			common.GStdout.Debug("Reenforcetask %v not exist", taskdata.marchEntityID)
			newCommand = true
			newMarch = true
		} else {
			common.GStdout.Debug("Reenforcetask %v exist", taskdata.marchEntityID)
			if entityData, ok := player.GetEntity(taskdata.marchEntityID); ok {
				if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_MarchCommand)); ok {
					marchCommand := property.(*protomsg.MarchCommand)
					if marchCommand.TargetType == protomsg.MarchCommandTarget_kMarchCommandTarget_Station {
						newCommand = true
					}
				}
			} else {
				newCommand = true
				newMarch = true
			}
		}
	}

	var commandData *protomsg.MarchCommand
	var targetID uint64
	var position *protomsg.Vector2D

	if newCommand {
		if taskdata.targetID != 0 && player.CheckTarget(taskdata.targetID) {
			targetID = taskdata.targetID
		} else {
			targetID, position, ok = player.SearchTarget(t.searchType, t.searchLevel)
			if !ok {
				common.GStdout.Debug("search targetID failed id %v\n", targetID)
				return false
			}
		}

		player.PM("addap 1000")
	}

	if t.viewBattle {
		request := &protomsg.MsgCL2GSViewMapRequest{}
		request.KingdomId = 0
		request.ViewLod = 0
		request.Position = position
		response, ok := player.SendAndWait(
			request,
			msgtype.MsgType_kMsgGS2CLViewMapReply).(*protomsg.MsgGS2CLViewMapReply)

		if !ok {
			return false
		}

		if response.ErrorCode != 0 {
			return false
		}
	}

	if newMarch {
		taskdata.marchIndex = 0
		taskdata.marchEntityID = 0
		taskdata.targetID = 0

		// 获取未出征的队伍id
		deply_id := player.GetDeployInfoId()
		if deply_id == nil {
			return true
		}

		entityData, ok := player.GetViewEntity(targetID)
		if !ok {
			return false
		}

		if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner)); ok {
			owernData := property.(*protomsg.EntityOwnerData)
			if t.searchType == protomsg.MapSearchType_kMapSearchType_Stronghold_1 || t.searchType == protomsg.MapSearchType_kMapSearchType_Stronghold_2 || t.searchType == protomsg.MapSearchType_kMapSearchType_Stronghold_3 || t.searchType == protomsg.MapSearchType_kMapSearchType_Stronghold_4 || t.searchType == protomsg.MapSearchType_kMapSearchType_Stronghold_5 {
				commandData, ok = t.DoStronghold(player, owernData, targetID, position)
				if !ok || commandData == nil {
					return false
				}

			} else if t.searchType == protomsg.MapSearchType_kMapSearchType_Castle {
				commandData, ok = t.DoCastle(player, owernData, targetID, position)
				if !ok || commandData == nil {
					return false
				}
			} else {
				return false
			}

			commandData.IsQuarter = t.isQuarter
			taskdata.targetID = targetID

			marchIndex, marchEntityID, ok := player.CreateTkMarch(*deply_id, commandData)
			if !ok {
				return false
			}

			taskdata.marchEntityID = marchEntityID
			taskdata.marchIndex = marchIndex
			taskdata.targetID = commandData.TargetId
		}
	} else if newCommand {
		commandData = MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, targetID, 0, 0)
		commandData.IsQuarter = t.isQuarter
		taskdata.targetID = targetID

		ok = player.MarchCommand(taskdata.marchIndex, commandData)
		if !ok {
			return false
		}
	}

	if t.waitMarchBack {
		player.WaitMarchBack(taskdata.marchEntityID, 0)
	}
	return true
}

func CreateBattleTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 7 {
		return nil, fmt.Errorf("battle task params error %v (%v)", len(params), params)
	}

	searchType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	searchLevel, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	if searchLevel == 0 {
		searchLevel = 1
	}

	waitMarchBack, error := strconv.ParseBool(params[2])
	if error != nil {
		return nil, error
	}

	viewBattle, error := strconv.ParseBool(params[3])
	if error != nil {
		return nil, error
	}

	isQuarter, error := strconv.Atoi(params[4])
	if error != nil {
		return nil, error
	}

	minArmy, error := strconv.Atoi(params[5])
	if error != nil {
		return nil, error
	}

	maxArmy, error := strconv.Atoi(params[6])
	if error != nil {
		return nil, error
	}

	task := &BattleTask{searchType: (protomsg.MapSearchType)(searchType), searchLevel: (uint32)(searchLevel), waitMarchBack: waitMarchBack, viewBattle: viewBattle, isQuarter: (int32)(isQuarter), minArmy: (uint32)(minArmy), maxArmy: (uint32)(maxArmy)}
	return task, nil
}

// BattleToTargetTask comment
type BattleToTargetTask struct {
	targetid      uint64
	waitMarchBack bool
	minArmy       uint32
	maxArmy       uint32
}

func (t *BattleToTargetTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	taskdata, ok := runtask.GetTaskData(taskindex).(*BattleTaskData)
	if !ok {
		taskdata = &BattleTaskData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	player.PM("addap 1000")
	commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, t.targetid, 0, 0)
	if taskdata.marchEntityID == 0 || !player.CheckEntityExist(taskdata.marchEntityID) {
		taskdata.marchIndex = 0
		taskdata.marchEntityID = 0
		taskdata.targetID = 0
		taskdata.skillGroupId = 0

		armyData := player.RandomArmyDataByRange(t.minArmy, t.maxArmy)
		if armyData == nil {
			common.GStdout.Debug("BattleToTargetTask false 1")
			return false
		}

		marchIndex, marchEntityID, ok := player.CreateMarch(armyData, commandData)
		if !ok {
			common.GStdout.Debug("BattleToTargetTask false 2")
			return false
		}

		taskdata.marchEntityID = marchEntityID
		taskdata.marchIndex = marchIndex
		taskdata.targetID = commandData.TargetId

		if skillGroupId, ok := data.HeroSkillGroup[armyData.Hero1]; ok {
			taskdata.skillGroupId = skillGroupId
		}
	} else {
		if entityData, ok := player.GetEntity(taskdata.marchEntityID); ok {
			if property, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_MarchCommand)); ok {
				marchCommand := property.(*protomsg.MarchCommand)
				if marchCommand.TargetType != protomsg.MarchCommandTarget_kMarchCommandTarget_ForceMoveBack &&
					marchCommand.TargetType != protomsg.MarchCommandTarget_kMarchCommandTarget_Battle {
					if !player.MarchCommand(taskdata.marchIndex, commandData) {
						common.GStdout.Debug("BattleToTargetTask false 3")
						return false
					}
				}
			}
		}
	}

	if taskdata.skillGroupId != 0 {
		if entityData, ok := player.GetEntity(taskdata.marchEntityID); ok {
			if entityproperty, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus)); ok {
				battleStatus := entityproperty.(*protomsg.BattleStatusData)
				if battleStatus.IsBattle > 0 {
					player.CastSkill(taskdata.marchIndex, taskdata.skillGroupId, taskdata.targetID)
				}
			}
		}
	}

	common.GStdout.Debug("BattleToTargetTask suc ")

	if t.waitMarchBack {
		player.WaitMarchBack(taskdata.marchEntityID, 0)
	}
	return true
}

func CreateBattleToTargetTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 4 {
		return nil, fmt.Errorf("battle task params error!")
	}

	target_id, error := strconv.ParseUint(params[0], 10, 64)
	if error != nil {
		return nil, error
	}

	waitMarchBack, error := strconv.ParseBool(params[1])
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

	task := &BattleToTargetTask{targetid: target_id, waitMarchBack: waitMarchBack, minArmy: (uint32)(minArmy), maxArmy: (uint32)(maxArmy)}
	return task, nil
}

// CollectTask 采集
type CollectTask struct {
	searchType    protomsg.MapSearchType
	searchLevel   uint32
	waitMarchBack bool
}

// CollectTaskData 采集
type CollectTaskData struct {
	marchEntityID uint64
	marchIndex    uint32
}

// DoTask comment
func (t *CollectTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	taskdata, ok := runtask.GetTaskData(taskindex).(*CollectTaskData)
	if !ok {
		taskdata = &CollectTaskData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	common.GStdout.Debug("collect %v %v %v", taskindex, taskdata.marchIndex, player.CheckEntityExist(taskdata.marchEntityID))

	if taskdata.marchEntityID > 0 && player.CheckEntityExist(taskdata.marchEntityID) {
		common.GStdout.Debug("collect moveback %v", taskdata.marchIndex)
		commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_MoveBack, 0, 0, 0)
		player.MarchCommand(taskdata.marchIndex, commandData)
		ok := player.WaitMarchBack(taskdata.marchEntityID, 100)

		if !ok {
			common.GStdout.Debug("collect false 1")
			return false
		}
		common.GStdout.Debug("collect moveback %v success.", taskdata.marchIndex)
		taskdata.marchIndex = 0
		taskdata.marchEntityID = 0
	}

	targetID, _, ok := player.SearchTarget(t.searchType, t.searchLevel)
	if !ok {
		common.GStdout.Debug("collect false 2")
		return false
	}
	// 随机军队数量
	armyData := player.RandomArmyData()
	if armyData == nil {
		common.GStdout.Debug("collect false 3")
		return false
	}

	commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Collect, targetID, 0, 0)
	marchIndex, marchEntityID, ok := player.CreateMarch(armyData, commandData)
	if !ok {
		common.GStdout.Debug("collect false 4")
		return false
	}

	taskdata.marchEntityID = marchEntityID
	taskdata.marchIndex = marchIndex
	common.GStdout.Debug("collect start %v %v success.", marchIndex, taskdata.marchEntityID)

	// 等待部队返回
	if t.waitMarchBack {
		player.WaitMarchBack(marchEntityID, 0)
		taskdata.marchEntityID = 0
	}
	return true
}

// CreateCollectTask 创建采集任务
func CreateCollectTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("collect task params error!")
	}

	searchType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	searchLevel, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	if searchLevel == 0 {
		searchLevel = 1
	}

	waitMarchBack, error := strconv.ParseBool(params[2])
	if error != nil {
		return nil, error
	}

	task := &CollectTask{searchType: (protomsg.MapSearchType)(searchType), searchLevel: (uint32)(searchLevel), waitMarchBack: waitMarchBack}
	return task, nil
}

// QuarterMarchTask 驻扎
type QuarterMarchTask struct {
}

// DoTask comment
func (t *QuarterMarchTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	marchFirst := CreateMarchFirst(player)
	if !marchFirst {
		return false
	}
	return true
}

// CreateQuarterMarchTask 创建驻扎任务
func CreateQuarterMarchTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &QuarterMarchTask{}
	return task, nil
}

// CreateMarchFirst 拉出一支部队
func CreateMarchFirst(player *Player) bool {
	common.GStdout.Debug("=========createMarch first.")
	armyData := player.RandomArmyData()
	if armyData == nil {
		return false
	}
	commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Station, 0, 0, 0)
	commandData.IsQuarter = 1

	player.PM("addap 1000")
	_, _, ok := player.CreateMarch(armyData, commandData)
	if !ok {
		return false
	}
	return true
}

// MarchMoveTask comment
type MarchMoveTask struct {
	MoveType    uint32 // 1 主堡为中心 2 指定坐标 3 随机中心坐标
	RangeType   uint32 // 1 以中心随机范围内查看 2 已之前的查看位置按随机范围内查看
	MoveRange   int64
	SrcPosition *protomsg.Vector2D
}

// MarchMoveTaskData comment
type MarchMoveTaskData struct {
	SrcPosition   *protomsg.Vector2D
	MovePosition  *protomsg.Vector2D
	marchIndex    uint32
	marchEntityID uint64
}

func (t *MarchMoveTask) getTaskData(runtask *BTree.BTreeRunTask, taskindex uint32) *MarchMoveTaskData {
	data := runtask.GetTaskData(taskindex)
	if data != nil {
		return data.(*MarchMoveTaskData)
	}

	task_data := &MarchMoveTaskData{SrcPosition: &protomsg.Vector2D{}, MovePosition: &protomsg.Vector2D{}}

	switch t.MoveType {
	case 1:
		player := runtask.TaskObj.(*Player)
		castlePos := player.GetCastlePos()
		common.GStdout.Debug("get castle pos(%v,%v)", castlePos.X, castlePos.Y)
		task_data.SrcPosition.X = castlePos.X
		task_data.SrcPosition.Y = castlePos.Y
		task_data.MovePosition = &protomsg.Vector2D{castlePos.X, castlePos.Y}
	case 2:
		task_data.SrcPosition.X = t.SrcPosition.X
		task_data.SrcPosition.Y = t.SrcPosition.Y
	case 3:
		task_data.SrcPosition.X = rand.Int63n(mapMaxPos)
		task_data.SrcPosition.Y = rand.Int63n(mapMaxPos)
	}

	task_data.MovePosition.X = task_data.SrcPosition.X
	task_data.MovePosition.Y = task_data.SrcPosition.Y
	runtask.SetTaskData(taskindex, task_data)

	return task_data
}

func (t *MarchMoveTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata := t.getTaskData(runtask, taskindex)

	switch t.RangeType {
	case 1:
		taskdata.MovePosition.X = taskdata.SrcPosition.X + rand.Int63n(t.MoveRange*2) - t.MoveRange
		taskdata.MovePosition.Y = taskdata.SrcPosition.Y + rand.Int63n(t.MoveRange*2) - t.MoveRange
	case 2:
		taskdata.MovePosition.X = taskdata.MovePosition.X + rand.Int63n(t.MoveRange*2) - t.MoveRange
		taskdata.MovePosition.Y = taskdata.MovePosition.Y + rand.Int63n(t.MoveRange*2) - t.MoveRange
	}

	if taskdata.MovePosition.X < 0 {
		taskdata.MovePosition.X = 0
	}

	if taskdata.MovePosition.X > mapMaxPos {
		taskdata.MovePosition.X = mapMaxPos
	}

	if taskdata.MovePosition.Y < 0 {
		taskdata.MovePosition.Y = 0
	}

	if taskdata.MovePosition.Y > mapMaxPos {
		taskdata.MovePosition.Y = mapMaxPos
	}

	common.GStdout.Debug("march move type (%v,%v,%v) src(%v,%v) view(%v,%v)", t.MoveType, t.RangeType, t.MoveRange, taskdata.SrcPosition.X, taskdata.SrcPosition.Y, taskdata.MovePosition.X, taskdata.MovePosition.Y)

	command_data := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Position, 0, taskdata.MovePosition.X, taskdata.MovePosition.Y)
	if taskdata.marchEntityID == 0 || !player.CheckEntityExist(taskdata.marchEntityID) {
		taskdata.marchIndex = 0
		taskdata.marchEntityID = 0
		army_data := player.RandomArmyData()
		if army_data == nil {
			common.GStdout.Debug(" march move 1 ")
			return false
		}

		create_command_data := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Station, 0, 0, 0)
		marchIndex, marchEntityID, ok := player.CreateMarch(army_data, create_command_data)
		if !ok {
			common.GStdout.Debug(" march move 2 ")
			return false
		}
		taskdata.marchEntityID = marchEntityID
		taskdata.marchIndex = marchIndex
		return true
	}
	return player.MarchCommand(taskdata.marchIndex, command_data)
}

// CreateMarchMoveTask comment
func CreateMarchMoveTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("create march move task params error")
	}

	moveType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	rangeType, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	moveRange, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	srcPosition := &protomsg.Vector2D{}
	switch moveType {
	case 2:
		if len(params) < 5 {
			return nil, fmt.Errorf("view map task params error!")
		}

		x, error := strconv.ParseInt(params[3], 10, 64)
		if error != nil {
			return nil, error
		}

		y, error := strconv.ParseInt(params[4], 10, 64)
		if error != nil {
			return nil, error
		}
		srcPosition.X = x * 1000 * 6
		srcPosition.Y = y * 1000 * 6
	case 3:
		srcPosition.X = rand.Int63n(mapMaxPos)
		srcPosition.Y = rand.Int63n(mapMaxPos)
	}

	task := &MarchMoveTask{MoveType: (uint32)(moveType), RangeType: (uint32)(rangeType), MoveRange: ((int64)(moveRange)) * 6 * 1000, SrcPosition: srcPosition}
	return task, nil
}

// CreateMarchMoveBackTask comment
func CreateMarchMoveBackTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &MarchMoveBackTask{}
	return task, nil
}

// MarchMoveBackTask comment
type MarchMoveBackTask struct {
}

func (t *MarchMoveBackTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	marchCount := player.GetMarchCount()
	common.GStdout.Debug(" march move back %v ", marchCount)
	commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_MoveBack, 0, 0, 0)
	common.GStdout.Debug(" %v", commandData)
	for i := 0; i < 10; i++ {
		if player.GetMarchCount() == 0 {
			return true
		}

		player.RangeEntities(func(entityID uint64, entityData *data.SyncEntityData) bool {
			if entityData.EntityType == (uint32)(protomsg.EntityType_kEntityType_March) {
				if propertype, ok := entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_March)); ok {
					player.MarchCommand(propertype.(*protomsg.EntityMarchData).MarchIndex, commandData)
				}
			}
			return true
		})
		time.Sleep(10 * time.Second)
	}
	return true
}

// CreateMoveCastleTask comment
func CreateMoveCastleTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("move castle task params error")
	}

	x, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	y, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	SearchRange, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	common.GStdout.Console("move castle x %v y %v", x, y)
	srcPosition := &protomsg.Vector2D{X: ((int64)(x)) * 6000, Y: ((int64)(y)) * 6000}

	task := &MoveCastleTask{CenterPos: srcPosition, SearchRange: (uint32)(SearchRange)}
	return task, nil
}

// MoveCastleTask 采集
type MoveCastleTask struct {
	CenterPos   *protomsg.Vector2D
	SearchRange uint32
}

// DoTask comment
func (t *MoveCastleTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	request := &protomsg.MsgCL2GSSearchEmptyPosRequest{}
	request.EntityType = protomsg.EntityType_kEntityType_Castle
	request.SearchLevel = 1
	request.SearchRange = t.SearchRange
	request.CenterPos = &protomsg.Vector2D{}
	request.CenterPos.X = t.CenterPos.X
	request.CenterPos.Y = t.CenterPos.Y

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSearchEmptyPosReply).(*protomsg.MsgGS2CLSearchEmptyPosReply)
	if !ok {
		return false
	}

	player.PM("movecastle 99 " + strconv.FormatInt(response.Position.X/6000, 10) + " " + strconv.FormatInt(response.Position.Y/6000, 10))
	return true
}

// ScoutTask 侦察
type ScoutTask struct {
	searchType  protomsg.MapSearchType
	searchLevel uint32
	scoutIndex  uint32
}

// DoTask comment
func (t *ScoutTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	targetID, _, ok := player.SearchTarget(t.searchType, t.searchLevel)
	if !ok {
		common.GStdout.Debug("do scout error , not found target")
		return false
	}

	player.PM("rich 100000")

	request := &protomsg.MsgCL2GSScoutRequest{}
	request.ScoutIndex = t.scoutIndex
	request.Command = &protomsg.ScoutCommand{}
	request.Command.TargetType = protomsg.ScoutCommandTarget_kScoutCommandTarget_DoScout
	request.Command.TargetId = targetID
	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLScoutReply).(*protomsg.MsgGS2CLScoutReply)
	if !ok {
		return false
	}

	return int32(response.ErrorCode) == 0
}

// CreateScoutTask 创建侦察任务
func CreateScoutTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("scout task params error!")
	}

	scoutIndex, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	searchType, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	searchLevel, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	if searchLevel == 0 {
		searchLevel = 1
	}

	task := &ScoutTask{searchType: (protomsg.MapSearchType)(searchType), searchLevel: uint32(searchLevel), scoutIndex: uint32(scoutIndex)}
	return task, nil
}

// ScoutTargetTask 侦察
type ScoutTargetTask struct {
	scoutIndex    uint32
	scoutTargetID uint64
}

// DoTask comment
func (t *ScoutTargetTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	player.PM("rich 10000")

	request := &protomsg.MsgCL2GSScoutRequest{}
	request.ScoutIndex = t.scoutIndex
	request.Command = &protomsg.ScoutCommand{}
	request.Command.TargetType = protomsg.ScoutCommandTarget_kScoutCommandTarget_DoScout
	request.Command.TargetId = t.scoutTargetID
	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLScoutReply).(*protomsg.MsgGS2CLScoutReply)
	if !ok {
		return false
	}

	return int32(response.ErrorCode) == 0
}

// CreateScoutTargetTask 创建侦察任务
func CreateScoutTargetTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("scout task target params error")
	}

	scoutIndex, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	targetID, error := strconv.ParseUint(params[1], 10, 64)
	if error != nil {
		return nil, error
	}

	task := &ScoutTargetTask{scoutIndex: (uint32)(scoutIndex), scoutTargetID: targetID}
	return task, nil
}

// ExploreMistTask 侦察
type ExploreMistTask struct {
	scoutIndex  uint32
	MoveType    uint32 // 1 主堡为中心 2 指定坐标 3 随机中心坐标
	RangeType   uint32 // 1 以中心随机范围内查看 2 已之前的查看位置按随机范围内查看
	MoveRange   int64
	SrcPosition *protomsg.Vector2D
}

// ExploreMistTaskData comment
type ExploreMistTaskData struct {
	SrcPosition    *protomsg.Vector2D
	targetPosition *protomsg.Vector2D
}

func (t *ExploreMistTask) getTaskData(runtask *BTree.BTreeRunTask, taskindex uint32) *ExploreMistTaskData {
	data := runtask.GetTaskData(taskindex)
	if data != nil {
		return data.(*ExploreMistTaskData)

	}

	taskData := &ExploreMistTaskData{SrcPosition: &protomsg.Vector2D{}, targetPosition: &protomsg.Vector2D{}}

	switch t.MoveType {
	case 1:
		player := runtask.TaskObj.(*Player)
		castlePos := player.GetCastlePos()
		common.GStdout.Debug("get castle pos(%v,%v)", castlePos.X, castlePos.Y)
		taskData.SrcPosition.X = castlePos.X
		taskData.SrcPosition.Y = castlePos.Y
		taskData.targetPosition = &protomsg.Vector2D{castlePos.X, castlePos.Y}
	case 2:
		taskData.SrcPosition.X = t.SrcPosition.X
		taskData.SrcPosition.Y = t.SrcPosition.Y
	case 3:
		taskData.SrcPosition.X = rand.Int63n(mapMaxPos)
		taskData.SrcPosition.Y = rand.Int63n(mapMaxPos)
	}

	taskData.targetPosition.X = taskData.SrcPosition.X
	taskData.targetPosition.Y = taskData.SrcPosition.Y
	runtask.SetTaskData(taskindex, taskData)

	return taskData
}

func (t *ExploreMistTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata := t.getTaskData(runtask, taskindex)

	player.PM("rich 100000")

	switch t.RangeType {
	case 1:
		taskdata.targetPosition.X = taskdata.SrcPosition.X + rand.Int63n(t.MoveRange*2) - t.MoveRange
		taskdata.targetPosition.Y = taskdata.SrcPosition.Y + rand.Int63n(t.MoveRange*2) - t.MoveRange
	case 2:
		taskdata.targetPosition.X = taskdata.targetPosition.X + rand.Int63n(t.MoveRange*2) - t.MoveRange
		taskdata.targetPosition.Y = taskdata.targetPosition.Y + rand.Int63n(t.MoveRange*2) - t.MoveRange
	}

	if taskdata.targetPosition.X < 0 {
		taskdata.targetPosition.X = 0
	}

	if taskdata.targetPosition.X > mapMaxPos {
		taskdata.targetPosition.X = mapMaxPos
	}

	if taskdata.targetPosition.Y < 0 {
		taskdata.targetPosition.Y = 0
	}

	if taskdata.targetPosition.Y > mapMaxPos {
		taskdata.targetPosition.Y = mapMaxPos
	}

	common.GStdout.Debug("explore mist type (%v,%v,%v) src(%v,%v) view(%v,%v)", t.MoveType, t.RangeType, t.MoveRange, taskdata.SrcPosition.X, taskdata.SrcPosition.Y, taskdata.targetPosition.X, taskdata.targetPosition.Y)

	var x int64 = taskdata.targetPosition.X / 1000
	var y int64 = taskdata.targetPosition.Y / 1000

	x = (x/270)*270 + 1
	y = (y/270)*270 + 1

	request := &protomsg.MsgCL2GSMistExploreRequest{}
	request.Route = make([]*protomsg.Vector2D, 0)
	var i int64 = 0
	var j int64 = 0
	for i = 0; i < 9; i++ {
		for j = 0; j < 9; j++ {
			pos := &protomsg.Vector2D{}
			pos.X = (x + i*30) * 1000
			pos.Y = (y + j*30) * 1000

			request.Route = append(request.Route, pos)
		}
	}
	request.ScoutIndex = t.scoutIndex
	if response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMistExploreReply).(*protomsg.MsgGS2CLMistExploreReply); ok {
		common.GStdout.Debug("explore mist result %v", response.ErrorCode)
		return response.ErrorCode == 0
	}
	return false
}

// CreateExploreMistTask comment
func CreateExploreMistTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 4 {
		return nil, fmt.Errorf("create explore mist task params error")
	}

	scoutIndex, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	moveType, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	rangeType, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	moveRange, error := strconv.Atoi(params[3])
	if error != nil {
		return nil, error
	}

	srcPosition := &protomsg.Vector2D{}
	switch moveType {
	case 2:
		if len(params) < 6 {
			return nil, fmt.Errorf("explore mist task params error")
		}

		x, error := strconv.ParseInt(params[4], 10, 64)
		if error != nil {
			return nil, error
		}

		y, error := strconv.ParseInt(params[5], 10, 64)
		if error != nil {
			return nil, error
		}
		srcPosition.X = x * 1000 * 6
		srcPosition.Y = y * 1000 * 6
	case 3:
		srcPosition.X = rand.Int63n(mapMaxPos)
		srcPosition.Y = rand.Int63n(mapMaxPos)
	}

	task := &ExploreMistTask{scoutIndex: (uint32)(scoutIndex), MoveType: (uint32)(moveType), RangeType: (uint32)(rangeType), MoveRange: ((int64)(moveRange)) * 6 * 1000, SrcPosition: srcPosition}
	return task, nil
}
