package data

import (
	"public/message/protomsg"

	"github.com/golang/protobuf/proto"
)

//GuildData 公会数据
type GuildData struct {
	// guildID 公会ID
	guildID uint64

	// guildName 公会名字
	guildName string

	// guildPosition 公会职务
	guildPosition protomsg.GuildPosition
}

// GGuildData xxx
var GGuildData *GuildData

func init() {
	GGuildData = &GuildData{}
}

// OnMsgGS2CLGuildBaseInfoReply  公会基础信息
func OnMsgGS2CLGuildBaseInfoReply(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLGuildBaseInfoReply)
	GGuildData.guildID = msg.Info.GuildId
	GGuildData.guildName = msg.Info.GuildName
	GGuildData.guildPosition = msg.Info.Position
}

// OnMsgGS2CLGuildJoinNotice  加入公会通知
func OnMsgGS2CLGuildJoinNotice(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLGuildJoinNotice)
	GGuildData.guildID = msg.Info.GuildId
	GGuildData.guildName = msg.Info.GuildName
	GGuildData.guildPosition = msg.Info.Position
}

// OnMsgGS2CLGuildQuitNotice  退出公会通知
func OnMsgGS2CLGuildQuitNotice(msginterface proto.Message) {
	ClearGuildInfo(0)
}

//OnMsgGS2CLGuildQuitReply 退出公会
func OnMsgGS2CLGuildQuitReply(msginterface proto.Message) {
	msg := msginterface.(*protomsg.MsgGS2CLGuildQuitReply)
	if msg.ErrorCode == 0 {
		ClearGuildInfo(0)
	}
}

func doGetGuildID() interface{} {
	return GGuildData.guildID
}

// GetGuildID 获取公会ID
func GetGuildID() uint64 {
	return GDataCenter.GetUint64Value(doGetGuildID)
}

func doGetGuildName() interface{} {
	return GGuildData.guildName
}

// GetGuildName 获取公会名字
func GetGuildName() string {
	return GDataCenter.GetStringValue(doGetGuildName)
}

// SetGuildName comment
func SetGuildName(param string) {

	GGuildData.guildName = param
}

// SetGuildID 设置公会ID
func SetGuildID(id uint64) {
	GGuildData.guildID = id
}

// SetGuildPosition 设置公会职务
func SetGuildPosition(param protomsg.GuildPosition) {
	GGuildData.guildPosition = param
}

// ClearGuildInfo comment
func ClearGuildInfo(param uint64) {
	GGuildData.guildID = 0
	GGuildData.guildName = ""
	GGuildData.guildPosition = protomsg.GuildPosition_kGuildPositionNone
}
