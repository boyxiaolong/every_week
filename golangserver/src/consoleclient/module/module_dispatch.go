package module

import (
	"consoleclient/handle"
	"public/connect"
	"public/message/msgtype"

	"github.com/golang/protobuf/proto"
)

//Dispatch comment
type Dispatch struct {
	connect.BaseDispatch
}

func init() {
	moduledipatch := &Dispatch{}
	handle.GhandleMsg.RegDispatch(handle.ModuleType, moduledipatch)
}

//Dispatch comment
func (m *Dispatch) Dispatch(msgtype msgtype.MsgType, msg proto.Message) {
	go GGameInfo.WaitMessage.Done(uint16(msgtype), msg)
}
