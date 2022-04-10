package module

import (
	"public/command"
	"public/common"
	"public/message/protomsg"
	"strings"
	"public/message/msgtype"
)

func init() {
	command.GCommand.RegCommand("pm", Pm, "Execute PM Command")
}

//Pm 执行PM命令
func Pm(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("Invalid PM Command")
		return
	}

	if GGameInfo.Online == false {
		common.GStdout.Error("user not online")
		return
	}

	params := make([]string, 0)
	for i, v := range str.Strs {
		if i > 0 {
			params = append(params, v)
		}
	}
	cmd := strings.Join(params, " ")

	common.GStdout.Info("pm:%v", cmd)
	msgCL2GSPMCommandRequest := &protomsg.MsgCL2GSPMCommandRequest{}
	msgCL2GSPMCommandRequest.Command = cmd

	GGameInfo.SendAndWait(msgCL2GSPMCommandRequest,msgtype.MsgType_kMsgGS2CLPMCommandReply)

	return
}
