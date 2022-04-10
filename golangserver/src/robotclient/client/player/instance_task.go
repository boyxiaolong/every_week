package player

import (
	"fmt"
	BTree "public/behaviortree"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
	"public/common"
	"robotclient/client/player/data"
	"math/rand"
	"time"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("IntoExpr", CreateIntoExprTask)
	BTree.RegisterTaskCreator("RandomExpr", CreateRandomExprTask)
}

func (p *Player) InstanceRandomArmyData( intance_id uint32 , is_hard bool ) *protomsg.ArmyData {
	if is_hard {
		return p.RandomArmyData()
	}

	army_settting, ok := data.InstancesArmy[intance_id]
	if !ok {
		return nil
	}

	var hero1 uint32
	var hero2 uint32
	regionData,ok:= p.GetMainRegion()
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

	if army_settting[0] == 0{
		return nil
	}

	rank_troops := data.ValidTroops[army_settting[0] - 1]
	if len(rank_troops) == 0 {
		return nil
	}
	
	troopid := rank_troops[rand.Intn(len(rank_troops))]

	troop_uints := army_settting[1]
	if troop_uints > p.MaxTroopCapacity() {
		troop_uints = p.MaxTroopCapacity()
	}

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

	time.Sleep(1 * time.Second)
	return MakeArmyData(hero1, hero2, troopid, troop_uints)
}

// IntoExprTask comment
type IntoExprTask struct {
	instanceid uint32
	Ishard   uint32
}

// IntoExprData comment
type IntoExprData struct {
	instanceId 		uint64
}

func (p *Player) InstanceMarchCommand(instanceId uint64,marchIndex uint32, command *protomsg.MarchCommand) bool {
	request := &protomsg.MsgCL2GSMarchCommandRequest{}
	request.Command = command
	request.MarchIndex = marchIndex
	request.RegionId  = instanceId
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLMarchCommandReply).(*protomsg.MsgGS2CLMarchCommandReply)
	if !ok {
		return false
	}
	common.GStdout.Info("=================== MarchCommand response %v", response.ErrorCode)
	return int32(response.ErrorCode) == 0
}

func (p *Player) GetInstance(instanceId uint64) (*data.RegionData,bool) {
	regionDataCenter := p.GetData("RegionDataCenter").(*data.RegionDataCenter)
	return regionDataCenter.GetRegion(instanceId)
}

func (p *Player) GetViewRegionId() uint64 {
	regionDataCenter := p.GetData("RegionDataCenter").(*data.RegionDataCenter)
	return regionDataCenter.ViewRegionId
}

// RangeEntities  通知地图数据
func (p *Player) RangeInstancePlayerEntities( instanceId uint64,fn func(entityID uint64, entityData *data.SyncEntityData) bool) {
	if instanceRegion, ok := p.GetInstance(instanceId); ok {
		instanceRegion.GetPlayerCenter().RangeEntities(fn)
	}
}

// RangeEntities  通知地图数据
func (p *Player) RangeInstanceViewEntities( instanceId uint64,fn func(entityID uint64, entityData *data.SyncEntityData) bool) {
	if instanceRegion, ok := p.GetInstance(instanceId); ok {
		instanceRegion.GetViewCenter().RangeEntities(fn)
	}
}

func (t *IntoExprTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata, ok := runtask.GetTaskData(taskindex).(*IntoExprData)
	if !ok {
		taskdata = &IntoExprData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	if _,ok:= player.GetInstance(taskdata.instanceId);!ok {
		common.GStdout.Info("IntoExprTask 1")
		taskdata.instanceId = 0
		army_data := player.InstanceRandomArmyData(t.instanceid,t.Ishard > 0)
		if army_data == nil {
			return false
		}
	
		request := &protomsg.MsgCL2GSPlayerExpeChallengeLevelRequest{}
		request.LevelId = t.instanceid
		request.HardMode = t.Ishard > 0
		request.ArmyArray = append(request.ArmyArray, army_data)
	
		response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeChallengeLevelReply).(*protomsg.MsgGS2CLPlayerExpeChallengeLevelReply)
		if !ok {
			common.GStdout.Info("IntoExprTask 2")
			return false
		}
		taskdata.instanceId = response.RegionId
		common.GStdout.Info("IntoExprTask 3")
	}

	if player.GetViewRegionId() != taskdata.instanceId {
		common.GStdout.Info("IntoExprTask 4")
		var marchEntity *data.SyncEntityData
		player.RangeInstancePlayerEntities(taskdata.instanceId,func(entityID uint64, entityData *data.SyncEntityData)bool{
			if entityData.EntityType != (uint32)(protomsg.EntityType_kEntityType_March) {
				return true
			}
			if property,ok:= entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner));ok{
				ownerData := property.(*protomsg.EntityOwnerData)
				if ownerData.PlayerId == player.PlayerID {
					marchEntity = entityData
					return false
				}
			}
			return true
		})	
			
		if marchEntity == nil {
			common.GStdout.Info("IntoExprTask 5")
			return false
		}

		common.GStdout.Info("IntoExprTask 5 5")

		if property,ok:= marchEntity.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Map));ok{
			mapData := property.(*protomsg.MapData)

			request := &protomsg.MsgCL2GSViewMapRequest{}
			request.ViewLod = 0
			request.Position = mapData.Position
			request.RegionId = taskdata.instanceId
		
			response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLViewMapReply).(*protomsg.MsgGS2CLViewMapReply)
			if !ok || response.ErrorCode != 0 {
				common.GStdout.Info("IntoExprTask 6")
				return false
			}
			common.GStdout.Info("IntoExprTask 6 6")
		} else {
			common.GStdout.Info("IntoExprTask 7")
			return false
		}
	}

	marchs := make(map[uint32]*data.SyncEntityData)
	player.RangeInstancePlayerEntities(taskdata.instanceId,func(entityID uint64, entityData *data.SyncEntityData)bool{
		if entityData.EntityType != (uint32)(protomsg.EntityType_kEntityType_March) {
			return true
		}

		if property,ok:= entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_March));ok{
			marchData := property.(*protomsg.EntityMarchData)
			marchs[marchData.MarchIndex] = entityData
		}
		return true
	})	

	enemies := make(map[uint64]*data.SyncEntityData)
	player.RangeInstanceViewEntities(taskdata.instanceId,func(entityID uint64, entityData *data.SyncEntityData)bool{
		if property,ok:= entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner));ok{
			ownerData := property.(*protomsg.EntityOwnerData)
			if ownerData.PlayerId != player.PlayerID {
				enemies[entityID] = entityData
				return false
			}
		}
		return true
	})

	if len(enemies) > 0{
		common.GStdout.Info("IntoExprTask 8")
		for k := range enemies{
			commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, k, 0,0)
			for marchIndex := range marchs {
				player.InstanceMarchCommand(taskdata.instanceId,marchIndex, commandData)
			}
			break
		}
	}else{
		common.GStdout.Info("IntoExprTask 9")
	}

	return true
}

func CreateIntoExprTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("battle task params error!")
	}

	instanceid, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	Ishard, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	task := &IntoExprTask{instanceid: (uint32)(instanceid), Ishard: (uint32)(Ishard)}
	return task, nil
}

// RandomExprTask comment
type RandomExprTask struct {
	BeginInstanceid uint32
	EndInstanceid 	uint32
	Ishard   				uint32
}

// RandomExprData comment
type RandomExprData struct {
	instanceId 		uint64
}

func (t *RandomExprTask) RandomInstanceData()(uint32,bool){
	randomRange := t.EndInstanceid - t.BeginInstanceid + 1
	if randomRange == 0 {
		return t.BeginInstanceid,t.Ishard > 0
	}
	return t.BeginInstanceid + (uint32)(rand.Intn((int)(randomRange))),t.Ishard > 0
}

func (t *RandomExprTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	taskdata, ok := runtask.GetTaskData(taskindex).(*RandomExprData)
	if !ok {
		taskdata = &RandomExprData{}
		runtask.SetTaskData(taskindex, taskdata)
	}

	if _,ok:= player.GetInstance(taskdata.instanceId);!ok {
		common.GStdout.Info("RandomExprTask 1")
		taskdata.instanceId = 0

		request := &protomsg.MsgCL2GSPlayerExpeChallengeLevelRequest{}
		id,hard := t.RandomInstanceData()
		if id == 0 {
			return false
		}

		army_data := player.InstanceRandomArmyData(id,hard)
		if army_data == nil {
			return false
		}
	
		request.LevelId = id
		request.HardMode = hard
		request.ArmyArray = append(request.ArmyArray, army_data)
	
		response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLPlayerExpeChallengeLevelReply).(*protomsg.MsgGS2CLPlayerExpeChallengeLevelReply)
		if !ok {
			common.GStdout.Info("RandomExprTask 2")
			return false
		}
		taskdata.instanceId = response.RegionId
		common.GStdout.Info("RandomExprTask 3")
	}

	if player.GetViewRegionId() != taskdata.instanceId {
		common.GStdout.Info("RandomExprTask 4")
		var marchEntity *data.SyncEntityData
		player.RangeInstancePlayerEntities(taskdata.instanceId,func(entityID uint64, entityData *data.SyncEntityData)bool{
			if entityData.EntityType != (uint32)(protomsg.EntityType_kEntityType_March) {
				return true
			}
			if property,ok:= entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner));ok{
				ownerData := property.(*protomsg.EntityOwnerData)
				if ownerData.PlayerId == player.PlayerID {
					marchEntity = entityData
					return false
				}
			}
			return true
		})	
			
		if marchEntity == nil {
			common.GStdout.Info("RandomExprTask 5")
			return false
		}

		common.GStdout.Info("RandomExprTask 5 5")

		if property,ok:= marchEntity.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Map));ok{
			mapData := property.(*protomsg.MapData)

			request := &protomsg.MsgCL2GSViewMapRequest{}
			request.ViewLod = 0
			request.Position = mapData.Position
			request.RegionId = taskdata.instanceId
		
			response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLViewMapReply).(*protomsg.MsgGS2CLViewMapReply)
			if !ok || response.ErrorCode != 0 {
				common.GStdout.Info("RandomExprTask 6")
				return false
			}
			common.GStdout.Info("RandomExprTask 6 6")
		} else {
			common.GStdout.Info("RandomExprTask 7")
			return false
		}
	}

	marchs := make(map[uint32]*data.SyncEntityData)
	player.RangeInstancePlayerEntities(taskdata.instanceId,func(entityID uint64, entityData *data.SyncEntityData)bool{
		if entityData.EntityType != (uint32)(protomsg.EntityType_kEntityType_March) {
			return true
		}

		if property,ok:= entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_March));ok{
			marchData := property.(*protomsg.EntityMarchData)
			marchs[marchData.MarchIndex] = entityData
		}
		return true
	})	

	enemies := make(map[uint64]*data.SyncEntityData)
	player.RangeInstanceViewEntities(taskdata.instanceId,func(entityID uint64, entityData *data.SyncEntityData)bool{
		if property,ok:= entityData.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Owner));ok{
			ownerData := property.(*protomsg.EntityOwnerData)
			if ownerData.PlayerId != player.PlayerID {
				enemies[entityID] = entityData
				return false
			}
		}
		return true
	})

	if len(enemies) > 0{
		common.GStdout.Info("RandomExprTask 8")
		for k := range enemies{
			commandData := MakeCommand(protomsg.MarchCommandTarget_kMarchCommandTarget_Battle, k, 0,0)
			for marchIndex := range marchs {
				player.InstanceMarchCommand(taskdata.instanceId,marchIndex, commandData)
			}
			break
		}
	}else{
		common.GStdout.Info("RandomExprTask 9")
	}

	return true
}

func CreateRandomExprTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("battle task params error!")
	}

	BeginInstanceid, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	EndInstanceid, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}


	Ishard, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	task := &RandomExprTask{BeginInstanceid: (uint32)(BeginInstanceid), EndInstanceid: (uint32)(EndInstanceid), Ishard: (uint32)(Ishard)}
	return task, nil
}
