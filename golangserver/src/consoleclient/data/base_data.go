package data

import (
	"public/common"
	"public/message/protomsg"

	"github.com/golang/protobuf/proto"
)

// GBaseData xxx
var GBaseData *BaseData

func init() {
	GBaseData = &BaseData{}
}

// BaseData struct
type BaseData struct {
	playerID   uint64
	playerName string
	dailyTasks []uint32
	castlePos  protomsg.Vector2D
}

// InitBasedData 初始化
func InitBasedData(msg proto.Message) {
	newMsg := msg.(*protomsg.MsgGS2CLPlayerBaseNotice)

	GBaseData.playerID = newMsg.PlayerId
	GBaseData.playerName = newMsg.Name

	common.GStdout.Success("init player base success")
}

// OnMsgGS2CLTaskDailyTaskInfoNotice  每日任务通知
func OnMsgGS2CLTaskDailyTaskInfoNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLTaskDailyTaskInfoNotice)
	GBaseData.dailyTasks = make([]uint32, len(msg.Data.Tasks))
	for i, task := range msg.Data.Tasks {
		GBaseData.dailyTasks[i] = task.TaskId
	}
}

// OnMsgGS2CLMapDataReply   玩家地图基础信息
func OnMsgGS2CLMapDataReply(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLMapDataReply)
	if msg.Castle != nil {
		GBaseData.castlePos = *msg.Castle.Position
	}

}

// OnMsgGS2CLMoveCastleReply  迁城后信息
func OnMsgGS2CLMoveCastleReply(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLMoveCastleReply)
	if msg.ErrorCode == 0 {
		GBaseData.castlePos = *msg.Position
	}
}

func doGetPlayerName() interface{} {
	return GBaseData.playerName
}

//SetPlayerName 设置名字
func SetPlayerName(name string) {
	GBaseData.playerName = name
	common.GStdout.Success("set player base success %v", GBaseData.playerName)
}

func doGetPlayerID() interface{} {
	return GBaseData.playerID
}

//GetPlayerID 获取玩家ID
func GetPlayerID() uint64 {
	return GDataCenter.GetUint64Value(doGetPlayerID)
}

// GetPlayerName 获取名字
func GetPlayerName() string {
	return GDataCenter.GetStringValue(doGetPlayerName)
}

func doGetFirstDaiylTaskID() interface{} {
	if len(GBaseData.dailyTasks) == 0 {
		return 0
	}

	return GBaseData.dailyTasks[0]
}

// GetFirstDaiylTaskID 获取每日任务第一个任务ID
func GetFirstDaiylTaskID() uint32 {
	return GDataCenter.GetUintValue(doGetFirstDaiylTaskID)
}

func doGetCastlePos() interface{} {
	return GBaseData.castlePos
}

// GetCastlePosX 获取主堡位置
func GetCastlePosX() int64 {
	return GDataCenter.GetData(doGetCastlePos).(protomsg.Vector2D).X / 1000
}

// GetCastlePosY 获取主堡位置
func GetCastlePosY() int64 {
	return GDataCenter.GetData(doGetCastlePos).(protomsg.Vector2D).Y / 1000
}
