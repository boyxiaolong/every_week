package data

import (
	"public/message/msgtype"
	"public/message/protomsg"

	"github.com/golang/protobuf/proto"
)

type RadarData struct {
	Radars []*protomsg.RadarData
}

func init() {
	RegisterDataCreator(CreateRadarData)
}

func CreateRadarData(data_center *DataCenter) {
	data := &RadarData{}
	data_center.DataRegister(data)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLRadarAllInfoNotice, data.OnMsgGS2CLRadarAllInfoNotice)
}

//
func (data *RadarData) OnMsgGS2CLRadarAllInfoNotice(msg proto.Message) {
	data_msg := msg.(*protomsg.MsgGS2CLRadarAllInfoNotice)
	data.Radars = data_msg.Radars
}
