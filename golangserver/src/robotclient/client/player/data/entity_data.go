package data

import (
	"bytes"
	"encoding/binary"
	"public/common"
	"public/message/protomsg"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
)

//GPropertyTypeName 属性名
var GPropertyTypeName = map[protomsg.EntityPropertyType]string{
	protomsg.EntityPropertyType_kEntityPropertyType_Map:              "protomsg.MapData",
	protomsg.EntityPropertyType_kEntityPropertyType_Owner:            "protomsg.EntityOwnerData",
	protomsg.EntityPropertyType_kEntityPropertyType_Move:             "protomsg.MoveData",
	protomsg.EntityPropertyType_kEntityPropertyType_GuildBuilding:    "protomsg.GuildBuildingData",
	protomsg.EntityPropertyType_kEntityPropertyType_Army:             "protomsg.ArmyData",
	protomsg.EntityPropertyType_kEntityPropertyType_Pass:             "protomsg.PassData",
	protomsg.EntityPropertyType_kEntityPropertyType_Temple:           "protomsg.TempleData",
	protomsg.EntityPropertyType_kEntityPropertyType_Collect:          "protomsg.CollectData",
	protomsg.EntityPropertyType_kEntityPropertyType_Reenforce:        "protomsg.DBReenforceData",
	protomsg.EntityPropertyType_kEntityPropertyType_BarbarianFort:    "protomsg.BarbarianFortData",
	protomsg.EntityPropertyType_kEntityPropertyType_Barbarian:        "protomsg.BarbariansData",
	protomsg.EntityPropertyType_kEntityPropertyType_CityWall:         "protomsg.CityWallData",
	protomsg.EntityPropertyType_kEntityPropertyType_EntityType:       "protomsg.EntityTypeData",
	protomsg.EntityPropertyType_kEntityPropertyType_MarchCommand:     "protomsg.MarchCommand",
	protomsg.EntityPropertyType_kEntityPropertyType_March:            "protomsg.EntityMarchData",
	protomsg.EntityPropertyType_kEntityPropertyType_Resource:         "protomsg.ResourceData",
	protomsg.EntityPropertyType_kEntityPropertyType_GuildMarch:       "protomsg.GuildMarchData",
	protomsg.EntityPropertyType_kEntityPropertyType_CastleExtendData: "protomsg.CastleExtendData",
	protomsg.EntityPropertyType_kEntityPropertyType_AttackTarget:     "protomsg.AttackTargetData",
	protomsg.EntityPropertyType_kEntityPropertyType_MoveSpeed:        "protomsg.MoveSpeedData",
	protomsg.EntityPropertyType_kEntityPropertyType_Alarm:            "protomsg.MapAlarmSyncDatas",
	protomsg.EntityPropertyType_kEntityPropertyType_SkillStatus:      "protomsg.SkillStatusData",
	protomsg.EntityPropertyType_kEntityPropertyType_MultiStrike:      "protomsg.MultiStrikeData",
	protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus:     "protomsg.BattleStatusData",
	protomsg.EntityPropertyType_kEntityPropertyType_ScoutCommand:     "protomsg.ScoutCommand",
	protomsg.EntityPropertyType_kEntityPropertyType_Scout:            "protomsg.EntityScoutData",
	protomsg.EntityPropertyType_kEntityPropertyType_Summon:           "protomsg.SummonData",
	protomsg.EntityPropertyType_kEntityPropertyType_Trap:             "protomsg.TrapData",
	protomsg.EntityPropertyType_kEntityPropertyType_GuildCollect:     "protomsg.GuildCollectBuildingData",
	protomsg.EntityPropertyType_kEntityPropertyType_GenericMonster:   "protomsg.GenericMonsterData",
	protomsg.EntityPropertyType_kEntityPropertyType_SummonMonster:    "protomsg.SummonMonsterData",
	protomsg.EntityPropertyType_kEntityPropertyType_ExpeStronghold:   "protomsg.ExpeStrongholdData",
	protomsg.EntityPropertyType_kEntityPropertyType_ReenforceData:    "protomsg.ReenforceData",
	protomsg.EntityPropertyType_kEntityPropertyType_ExpeEnemyFort:    "protomsg.ExpeEnemyFortData",
	protomsg.EntityPropertyType_kEntityPropertyType_MovePath:         "protomsg.MovePathData",
	protomsg.EntityPropertyType_kEntityPropertyType_PompeiiBuild:     "protomsg.PompeiiBuildData",
	protomsg.EntityPropertyType_kEntityPropertyType_Static_Enemy:     "protomsg.AttackStickEnemyData",
	protomsg.EntityPropertyType_kEntityPropertyType_ExpeGroupMarch:   "protomsg.ExpeGroupMarchData",
	protomsg.EntityPropertyType_kEntityPropertyType_TrafficUnit:      "protomsg.TrafficUnitData",
	protomsg.EntityPropertyType_kEntityPropertyType_CarriageData:     "protomsg.EntityCarriageData",
	protomsg.EntityPropertyType_kEntityPropertyType_GarrisonHero:     "protomsg.GarrisonHeroData",
	protomsg.EntityPropertyType_kEntityPropertyType_PompeiiMonster:   "protomsg.EntityPompeiiMonsterData",
	protomsg.EntityPropertyType_kEntityPropertyType_PompeiiBoss:      "protomsg.EntityPompeiiBossData",
	protomsg.EntityPropertyType_kEntityPropertyType_GuildMarchLimit:  "protomsg.GuildMarchLimit",
	protomsg.EntityPropertyType_kEntityPropertyType_KingdomTitle:     "protomsg.EntityKingdomTitle",
	protomsg.EntityPropertyType_kEntityPropertyType_ScoutExtendData:  "protomsg.ScoutExtendData",
	protomsg.EntityPropertyType_kEntityPropertyType_BaseArmy:         "protomsg.ArmyData",
}

var GEntityPropertyTypeSet = map[protomsg.EntityType]map[protomsg.EntityPropertyType]bool{
	protomsg.EntityType_kEntityType_StrongHold: {
		protomsg.EntityPropertyType_kEntityPropertyType_Owner: true,
	},

	protomsg.EntityType_kEntityType_Castle: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
	},
	protomsg.EntityType_kEntityType_March: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_March:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_Army:         true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
		protomsg.EntityPropertyType_kEntityPropertyType_MarchCommand: true,
	},
	protomsg.EntityType_kEntityType_Scout: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:   true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner: true,
		protomsg.EntityPropertyType_kEntityPropertyType_Scout: true,
	},
	protomsg.EntityType_kEntityType_Carriage: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_CarriageData: true,
	},
	protomsg.EntityType_kEntityType_GuildMarch: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_GuildMarch:   true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
		protomsg.EntityPropertyType_kEntityPropertyType_MarchCommand: true,
	},
	protomsg.EntityType_kEntityType_GuildFort: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
	},
	protomsg.EntityType_kEntityType_GuildFlag: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
	},
	protomsg.EntityType_kEntityType_Temple: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
	},
	protomsg.EntityType_kEntityType_Pass: {
		protomsg.EntityPropertyType_kEntityPropertyType_Map:          true,
		protomsg.EntityPropertyType_kEntityPropertyType_Owner:        true,
		protomsg.EntityPropertyType_kEntityPropertyType_BattleStatus: true,
	},
}

type SyncEntityData struct {
	EntityId   uint64
	EntityType uint32
	properties *sync.Map
}

func CreateSyncEntityData() *SyncEntityData {
	data := &SyncEntityData{
		EntityId:   0,
		EntityType: 0,
		properties: &sync.Map{},
	}
	return data
}

func (s *SyncEntityData) GetProperty(propertyType uint32) (proto.Message, bool) {
	if property, ok := s.properties.Load(propertyType); ok {
		return property.(proto.Message), true
	}
	return nil, false
}

func (s *SyncEntityData) GetPosition() *protomsg.Vector2D {
	if property, ok := s.properties.Load((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Map)); ok {
		mapData := property.(*protomsg.MapData)
		return mapData.Position
	}
	return &protomsg.Vector2D{}
}

type CenterManagerInterface interface {
	OnNewEntity(centerType uint32, s *SyncEntityData)
	OnRemoveEntity(centerType uint32, s *SyncEntityData)
}

type EntityCenter struct {
	CenterType uint32
	Entities   *sync.Map
	CenterMgr  CenterManagerInterface
}

func (d *EntityCenter) GetEntityOrNew(entityId uint64) (*SyncEntityData, bool) {
	if v, ok := d.Entities.Load(entityId); ok {
		return v.(*SyncEntityData), false
	} else {
		syncEntityData := CreateSyncEntityData()
		d.Entities.Store(entityId, syncEntityData)
		return syncEntityData, true
	}
}

func (d *EntityCenter) GetEntity(entityId uint64) (*SyncEntityData, bool) {
	if v, ok := d.Entities.Load(entityId); ok {
		return v.(*SyncEntityData), true
	}
	return nil, false
}

func (d *EntityCenter) RemoveEntity(entityId uint64) (*SyncEntityData, bool) {
	if v, ok := d.Entities.Load(entityId); ok {
		common.GStdout.Debug("EntityCenter %v remove entity {}", d.CenterType, entityId)
		d.Entities.Delete(entityId)
		return v.(*SyncEntityData), true
	}
	return nil, false
}

func (d *EntityCenter) GetEntityCount(entityType uint32) uint32 {
	var count uint32
	d.Entities.Range(func(_, v interface{}) bool {
		entity := v.(*SyncEntityData)
		if entity.EntityType == entityType {
			count++
		}
		return true
	})
	return count
}

func (d *EntityCenter) RangeEntities(fn func(entityID uint64, entityData *SyncEntityData) bool) {
	d.Entities.Range(func(k, v interface{}) bool {
		return fn((k).(uint64), v.(*SyncEntityData))
	})
	return
}

func (center *EntityCenter) UpdateSyncEntitiesData(syncData []byte) {
	buf := bytes.NewReader(syncData)
	var entityCount uint32
	err := binary.Read(buf, binary.LittleEndian, &entityCount)
	if err != nil {
		common.GStdout.Debug("%v\n", err)
		return
	}

	//common.GStdout.Debug(" region id %v, type %v , entity count %v \n", msg.RegionId , dataType, entityCount)
	for i := 0; i < int(entityCount); i++ {
		var entityId uint64
		var entityType uint32
		var version uint32
		err = binary.Read(buf, binary.LittleEndian, &entityType)
		if err != nil {
			common.GStdout.Debug("%v\n", err)
			return
		}

		err = binary.Read(buf, binary.LittleEndian, &entityId)
		if err != nil {
			common.GStdout.Debug("%v\n", err)
			return
		}

		err = binary.Read(buf, binary.LittleEndian, &version)
		if err != nil {
			common.GStdout.Debug("%v\n", err)
			return
		}

		var propertyCount uint32
		err = binary.Read(buf, binary.LittleEndian, &propertyCount)
		if err != nil {
			common.GStdout.Debug("%v\n", err)
			return
		}

		// common.GStdout.Debug("entity id %v , entity type %v , property count %v\n", entityId, entityType, propertyCount)
		var syncEntityData *SyncEntityData
		var isnew bool
		properties, ok := GEntityPropertyTypeSet[(protomsg.EntityType)(entityType)]
		if ok {
			syncEntityData, isnew = center.GetEntityOrNew(entityId)
		}

		// if (protomsg.EntityType)(entityType) == protomsg.EntityType_kEntityType_StrongHold {
		// 	common.GStdout.Debug("entity id %v , entity type %v , property count %v\n", entityId, entityType, propertyCount)
		// }

		for j := 0; j < int(propertyCount); j++ {
			var propertyType uint32
			var propertyLength uint32
			err = binary.Read(buf, binary.LittleEndian, &propertyType)
			if err != nil {
				common.GStdout.Debug("%v\n", err)
				return
			}

			err = binary.Read(buf, binary.LittleEndian, &propertyLength)
			if err != nil {
				common.GStdout.Debug("%v\n", err)
				return
			}

			propertyData := make([]byte, propertyLength)

			err = binary.Read(buf, binary.LittleEndian, propertyData)
			if err != nil {
				common.GStdout.Debug("%v\n", err)
				return
			}

			if syncEntityData == nil {
				continue
			}

			syncEntityData.EntityId = entityId
			syncEntityData.EntityType = entityType
			if _, ok := properties[(protomsg.EntityPropertyType)(propertyType)]; ok {
				//common.GStdout.Debug("property type %v , property length %v , data[%v]\n", propertyType, propertyLength, propertyData)
				propertyMsg := MakeProperty((protomsg.EntityPropertyType)(propertyType), propertyData)
				// if propertyMsg != nil {
				// 	PrintPropertyDetail(entityID, entitytype, propertyType, propertyMsg)
				// }
				syncEntityData.properties.Store(propertyType, propertyMsg)
			}
		}

		if isnew && center.CenterMgr != nil {
			center.CenterMgr.OnNewEntity(center.CenterType, syncEntityData)
		}
	}
	// common.GStdout.Debug("========================================= recv UpdateSyncEntitieNotice end\n")
}

func (center *EntityCenter) RemoveEntitiesData(EntityIds []uint64) {
	for _, v := range EntityIds {
		if entity, ok := center.RemoveEntity(v); ok {
			if center.CenterMgr != nil {
				center.CenterMgr.OnRemoveEntity(center.CenterType, entity)
			}
		}
	}
}

func (center *EntityCenter) Clear() {
	center.Entities = &sync.Map{}
}

type RegionData struct {
	RegionId      uint64
	EntityCenters *sync.Map
	UsedHeros     *sync.Map
}

func CreateRegionData() *RegionData {
	data := &RegionData{
		RegionId:      0,
		EntityCenters: &sync.Map{},
		UsedHeros:     &sync.Map{},
	}

	viewCenter := &EntityCenter{CenterType: 1, Entities: &sync.Map{}, CenterMgr: data}
	worldCenter := &EntityCenter{CenterType: 2, Entities: &sync.Map{}, CenterMgr: data}
	playerCenter := &EntityCenter{CenterType: 3, Entities: &sync.Map{}, CenterMgr: data}
	guildCenter := &EntityCenter{CenterType: 4, Entities: &sync.Map{}, CenterMgr: data}

	data.EntityCenters.Store(uint32(1), viewCenter)
	data.EntityCenters.Store(uint32(2), worldCenter)
	data.EntityCenters.Store(uint32(3), playerCenter)
	data.EntityCenters.Store(uint32(4), guildCenter)
	return data
}

func (d *RegionData) GetEntityCenter(centerType uint32) *EntityCenter {
	if center, ok := d.EntityCenters.Load(centerType); ok {
		return center.(*EntityCenter)
	}
	return nil
}

func (d *RegionData) GetViewCenter() *EntityCenter {
	return d.GetEntityCenter(1)
}

func (d *RegionData) GetWorldCenter() *EntityCenter {
	return d.GetEntityCenter(2)
}

func (d *RegionData) GetPlayerCenter() *EntityCenter {
	return d.GetEntityCenter(3)
}

func (d *RegionData) GetGuildCenter() *EntityCenter {
	return d.GetEntityCenter(4)
}

func (d *RegionData) Clear() {
	d.RegionId = 0
	d.EntityCenters = &sync.Map{}
	d.UsedHeros = &sync.Map{}

	viewCenter := &EntityCenter{CenterType: 1, Entities: &sync.Map{}, CenterMgr: d}
	worldCenter := &EntityCenter{CenterType: 2, Entities: &sync.Map{}, CenterMgr: d}
	playerCenter := &EntityCenter{CenterType: 3, Entities: &sync.Map{}, CenterMgr: d}
	guildCenter := &EntityCenter{CenterType: 4, Entities: &sync.Map{}, CenterMgr: d}
	d.EntityCenters.Store(uint32(1), viewCenter)
	d.EntityCenters.Store(uint32(2), worldCenter)
	d.EntityCenters.Store(uint32(3), playerCenter)
	d.EntityCenters.Store(uint32(4), guildCenter)
}

type RegionDataCenter struct {
	MainRegionId uint64
	RegionDatas  *sync.Map
	ViewRegionId uint64
}

func CreateRegionDataCenter() *RegionDataCenter {
	data := &RegionDataCenter{
		MainRegionId: 0,
		RegionDatas:  &sync.Map{},
		ViewRegionId: 0,
	}
	return data
}

func (c *RegionDataCenter) GetMainRegion() (*RegionData, bool) {
	if v, ok := c.RegionDatas.Load(c.MainRegionId); ok {
		return v.(*RegionData), true
	}
	return nil, false
}

func (c *RegionDataCenter) GetViewRegion() (*RegionData, bool) {
	if v, ok := c.RegionDatas.Load(c.ViewRegionId); ok {
		return v.(*RegionData), true
	}
	return nil, false
}

func (c *RegionDataCenter) GetRegion(regionid uint64) (*RegionData, bool) {
	if v, ok := c.RegionDatas.Load(regionid); ok {
		return v.(*RegionData), true
	}
	return nil, false
}

func (c *RegionDataCenter) CreateRegion(regionid uint64) *RegionData {
	if v, ok := c.RegionDatas.Load(regionid); ok {
		return v.(*RegionData)
	} else {
		regionData := CreateRegionData()
		regionData.RegionId = regionid
		c.RegionDatas.Store(regionid, regionData)
		return regionData
	}
}

func (c *RegionDataCenter) DeleteRegion(regionid uint64) {
	c.RegionDatas.Delete(regionid)
}

func (c *RegionDataCenter) LoadRegionCenter(regionid uint64, centerType uint32) *EntityCenter {
	regionData := c.CreateRegion(regionid)
	if regionData == nil {
		return nil
	}
	center := regionData.GetEntityCenter(centerType)
	return center
}

func (c *RegionDataCenter) GetRegionCenter(regionid uint64, centerType uint32) *EntityCenter {
	if regionData, ok := c.GetRegion(regionid); ok {
		return regionData.GetEntityCenter(centerType)
	}
	return nil
}

func (c *RegionDataCenter) Clear() {
	c.MainRegionId = 0
	c.RegionDatas = &sync.Map{}
}

func (regionData *RegionData) OnNewEntity(centerType uint32, s *SyncEntityData) {
	if centerType == 3 && s.EntityType == 1 { // 行军需要记录下hero id
		if propery, ok := s.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Army)); ok {
			armyData := propery.(*protomsg.ArmyData)
			if armyData.Hero1 != 0 {
				regionData.UsedHeros.Store(armyData.Hero1, true)
			}

			if armyData.Hero2 != 0 {
				regionData.UsedHeros.Store(armyData.Hero2, true)
			}
		}
	}
}

func (regionData *RegionData) OnRemoveEntity(centerType uint32, s *SyncEntityData) {
	if centerType == 3 && s.EntityType == 1 {
		if propery, ok := s.GetProperty((uint32)(protomsg.EntityPropertyType_kEntityPropertyType_Army)); ok {
			armyData := propery.(*protomsg.ArmyData)
			if armyData.Hero1 != 0 {
				regionData.UsedHeros.Delete(armyData.Hero1)
			}

			if armyData.Hero2 != 0 {
				regionData.UsedHeros.Delete(armyData.Hero2)
			}
		}
	}
}

// GetPropertyDataType content
func GetPropertyDataType(propertyType protomsg.EntityPropertyType) (bool, reflect.Type) {
	name, ok := GPropertyTypeName[propertyType]
	if !ok {
		common.GStdout.Error("propertyType not exist:%v", propertyType)
		return false, nil
	}

	t := proto.MessageType(name)
	if t == nil {
		common.GStdout.Error("propertyType message not exist:%v,name:%v", propertyType, name)
		return false, nil
	}
	return true, t
}

// MakeProperty content
func MakeProperty(propertyType protomsg.EntityPropertyType, data []byte) proto.Message {
	ok, t := GetPropertyDataType(propertyType)
	if !ok {
		return nil
	}

	if propertyMsg, ok := reflect.New(t.Elem()).Interface().(proto.Message); ok {
		err := proto.Unmarshal(data, propertyMsg)
		if err != nil {
			common.GStdout.Error("protobuf 解码出错:%v", err.Error())
			return nil
		}
		return propertyMsg
	}
	return nil
}
