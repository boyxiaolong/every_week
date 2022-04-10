package connect

import (
	"encoding/json"
	"public/command"
	"public/common"
	"public/config"
	"public/message/msgtype"

	"github.com/golang/protobuf/proto"
)

var (
	MSG_PRINT_NONE   = 0
	MSG_PRINT_BRIEF  = 1
	MSG_PRINT_DETAIL = 2

	GMsgPrint *MsgPrint
)

func init() {
	GMsgPrint = &MsgPrint{
		level: int32(config.MsgLogLevel),
	}

	command.GCommand.RegCommand("showmsg", ShowMsg, "Set console show msg level,0 is PRINT_NONE, 1 is PRINT_BRIEF,2 is PRINT_DETAIL")
}

func ShowMsg(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("showmsg param error")
		return
	}

	GMsgPrint.SetLevel(int32(str.GetInt(1)))

	return
}

type MsgPrint struct {
	level int32
}

func (m *MsgPrint) SetLevel(level int32) {
	if level < int32(MSG_PRINT_NONE) || level > int32(MSG_PRINT_DETAIL) {
		common.GStdout.Error("print set level param error:%d", level)
		return
	}

	m.level = level
	common.GStdout.Success("Console msg print level : %v", m.level)
}

func (m *MsgPrint) OnMessage(msg_id uint16, pbMessage proto.Message, size int) {
	if msg_id == uint16(msgtype.MsgType_kMsgGS2CLKeepLiveReply) {
		return
	}

	if m.level == int32(MSG_PRINT_BRIEF) {
		m.PrintBrief(msg_id, pbMessage, size)
		return
	}

	if m.level == int32(MSG_PRINT_DETAIL) {
		m.PrintDetail(msg_id, pbMessage, size)
		return
	}
}

func (m *MsgPrint) PrintDetail(msg_id uint16, pbMessage proto.Message, size int) {
	strMsgName := msgtype.MsgType_name[int32(msg_id)]
	data, err := json.Marshal(pbMessage)
	if err != nil {
		common.GStdout.Error("JSON marshaling failed: %s", err)
		return
	}
	common.GStdout.Info("receive msssage success:msg_name:%v,msg_type:%v,msg_size:%v,data:%s", strMsgName, msg_id, size, data)
}

func (m *MsgPrint) PrintBrief(msg_id uint16, pbMessage proto.Message, size int) {
	strMsgName := msgtype.MsgType_name[int32(msg_id)]
	common.GStdout.Info("receive msssage success:msg_name:%v,msg_type:%v,msg_size:%v", strMsgName, msg_id, size)
}
