package data

import (
	"consoleclient/application"
	"consoleclient/handle"
	"fmt"
	"public/connect"
	"public/message/msgtype"

	"github.com/golang/protobuf/proto"
)

// DataDispatch xxx
type DataDispatch struct {
	connect.BaseDispatch
}

func init() {
	moduledipatch := &DataDispatch{}
	moduledipatch.Init()
	moduledipatch.RegMsg()
	handle.GhandleMsg.RegDispatch(handle.Datatype, moduledipatch)
}

//RegMsg 注册消息
func (m *DataDispatch) RegMsg() {
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLPlayerBaseNotice, InitBasedData)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildBaseInfoReply, OnMsgGS2CLGuildBaseInfoReply)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildJoinNotice, OnMsgGS2CLGuildJoinNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildQuitNotice, OnMsgGS2CLGuildQuitNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLGuildQuitReply, OnMsgGS2CLGuildQuitReply)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLSyncEntitiesDataNotice, OnMsgGS2CLUpdateSyncEntitieNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLCityInfoReply, OnMsgGS2CLCityInfoReply)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLCityBuildingUpdateNotice, OnMsgGS2CLCityBuildingUpdateNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLTaskDailyTaskInfoNotice, OnMsgGS2CLTaskDailyTaskInfoNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLMapDataReply, OnMsgGS2CLMapDataReply)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLMoveCastleReply, OnMsgGS2CLMoveCastleReply)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLAllEquipmentNotice, OnMsgGS2CLAllEquipmentNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLAddEquipmentsNotice, OnMsgGS2CLAddEquipmentsNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLEquipmentNotice, OnMsgGS2CLEquipmentNotice)
	m.RegDispatch(msgtype.MsgType_kMsgGS2CLRemoveEquipmentsNotice, OnMsgGS2CLRemoveEquipmentsNotice)
}

//Dispatch 分发消息
func (m *DataDispatch) Dispatch(msgtype msgtype.MsgType, msg proto.Message) {
	application.GetApplication().DataTask.AddTask(func() {
		if call, ok := m.CallBacks[msgtype]; ok {
			fmt.Printf("Dispatch %v \n", msgtype)
			call(msg)
		}
	})
}
