package data

import (
	"public/message/msgtype"
	"public/message/protomsg"

	"github.com/golang/protobuf/proto"
)

type BaseData struct {
	Player_id     uint64
	Player_name   string
	Player_Emoney uint64
}

func init() {
	RegisterDataCreator(CreateBaseData)
}

func CreateBaseData(data_center *DataCenter) {
	data := &BaseData{}
	data_center.DataRegister(data)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLPlayerBaseNotice, data.InitBasedData)
	data_center.RegDispatch(msgtype.MsgType_kMsgGS2CLPlayerAttribute, data.UpdatePlayerAttribute)
}

func (data *BaseData) InitBasedData(msg proto.Message) {
	new_msg := msg.(*protomsg.MsgGS2CLPlayerBaseNotice)

	data.Player_id = new_msg.PlayerId
	data.Player_name = new_msg.Name
	data.Player_Emoney = new_msg.Currencies[1]

	// common.GStdout.Info("init player base success")
}

// UpdatePlayerAttribute comment
func (data *BaseData) UpdatePlayerAttribute(msg proto.Message) {
	newMsg := msg.(*protomsg.MsgGS2CLPlayerAttribute)

	if len(newMsg.Attributes) == 0 {
		return
	}
	if newMsg.Attributes[0].Type == uint32(protomsg.PlayerAttributeType_kPlayerAttrCurrency) && newMsg.Attributes[0].SubType == uint32(protomsg.CurrencyType_kCurrencyTypeEmoney) {
		data.Player_Emoney = newMsg.Attributes[0].Value
	}

	// common.GStdout.Info("update player attribute success")
}
