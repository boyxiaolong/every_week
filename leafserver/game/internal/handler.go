package internal

import (
	"reflect"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	handler(&msg.ReqLogin{}, handleReqlogin)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleReqlogin(args []interface{}) {
	m := args[0].(*msg.ReqLogin)
	a := args[1].(gate.Agent)

	log.Debug("reqlogin %v", m)
	resp := msg.RespLogin{}
	a.WriteMsg(&resp)
}
