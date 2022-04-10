package data

import (
	"public/message/msgtype"
	"public/message/protomsg"
	"github.com/golang/protobuf/proto"
)

type GuildData struct {
	// GuildID 公会ID
	GuildID uint64

	// GuildName 公会名字
	GuildName string

	// GuildPosition 公会职务
	GuildPosition protomsg.GuildPosition

	// GuildMarchs *sync.Map
}

func init() {
	RegisterDataCreator(CreateGuildData)
}

func CreateGuildData(data_center *DataCenter) {
	data := &GuildData{}
	data_center.DataRegister(data)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildBaseInfoReply, data.OnMsgGS2CLGuildBaseInfoReply)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildJoinNotice, data.OnMsgGS2CLGuildJoinNotice)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildQuitNotice, data.OnMsgGS2CLGuildQuitNotice)
	//data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLGetGuildMarchListReply, data.OnMsgGS2CLGetGuildMarchListReply)
	//data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLAddGuildMarchNotice, data.OnMsgGS2CLAddGuildMarchNotice)
	// data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLRemoveGuildMarchNotice, data.OnMsgGS2CLRemoveGuildMarchNotice)

}

// OnMsgGS2CLGuildBaseInfoReply  公会基础信息
func (data *GuildData) OnMsgGS2CLGuildBaseInfoReply(msg proto.Message) {
	msg_data := msg.(*protomsg.MsgGS2CLGuildBaseInfoReply)
	data.GuildID = msg_data.Info.GuildId
	data.GuildName = msg_data.Info.GuildName
	data.GuildPosition = msg_data.Info.Position
}

// OnMsgGS2CLGuildJoinNotice  加入公会通知
func (data *GuildData) OnMsgGS2CLGuildJoinNotice(msg proto.Message) {
	msg_data := msg.(*protomsg.MsgGS2CLGuildJoinNotice)
	data.GuildID = msg_data.Info.GuildId
	data.GuildName = msg_data.Info.GuildName
	data.GuildPosition = msg_data.Info.Position
}

// OnMsgGS2CLGuildQuitNotice  退出公会通知
func (data *GuildData) OnMsgGS2CLGuildQuitNotice(msg proto.Message) {
	data.ClearGuildInfo(0)
}

// // OnMsgGS2CLGetGuildMarchListReply  退出公会通知
// func (data *GuildData) OnMsgGS2CLGetGuildMarchListReply(msg proto.Message) {
// 	data.GuildMarchs = &sync.Map{}
// 	data_msg := msg.(*protomsg.MsgGS2CLGetGuildMarchListReply)
// 	for _, v := range data_msg.Datas {
// 		data.GuildMarchs.Store(v.EntityId, true)
// 	}
// }

// // OnMsgGS2CLAddGuildMarchNotice  退出公会通知
// func (data *GuildData) OnMsgGS2CLAddGuildMarchNotice(msg proto.Message) {
// 	data_msg := msg.(*protomsg.MsgGS2CLAddGuildMarchNotice)
// 	data.GuildMarchs.Store(data_msg.Data.EntityId, true)
// }

// OnMsgGS2CLRemoveGuildMarchNotice  退出公会通知
// func (data *GuildData) OnMsgGS2CLRemoveGuildMarchNotice(msg proto.Message) {
// 	data_msg := msg.(*protomsg.MsgGS2CLRemoveGuildMarchNotice)
// 	data.GuildMarchs.Delete(data_msg.GuildMarchEntityId)
// }

// ClearGuildInfo comment
func (data *GuildData) ClearGuildInfo(param uint64) {
	data.GuildID = 0
	data.GuildName = ""
	data.GuildPosition = protomsg.GuildPosition_kGuildPositionNone
	// data.GuildMarchs = &sync.Map{}
}
