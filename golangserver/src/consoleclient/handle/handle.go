package handle

import (
	"public/connect"
)

var ModuleType uint16 = 1
var Datatype uint16 = 2

var GhandleMsg *HandleMsg

func init() {
	GhandleMsg  = MakeHandleMsg()
	GhandleMsg.Init()
}

func MakeHandleMsg() *HandleMsg {
	return &HandleMsg{}
}

type HandleMsg struct{
	connect.HandlerMsgBase
}
