package connect

import (
	"public/common"
	"public/message/msgtype"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

type ResFun func(msg interface{}) (err error)

type MsgInfo struct {
	msg     proto.Message
	session *Session
}

type MessageHandler interface {
	GetMsgID(msg proto.Message) uint16
	GetMsgType(msg_type uint16) (bool, reflect.Type)
	DoDispatch(msg proto.Message, session *Session)
	RegDispatch(dispatch_type uint16, dispatch DispatchInterface)
}

type HandlerMsgBase struct {
	Dispatchs map[uint16]DispatchInterface
}

func (m *HandlerMsgBase) Init() {
	m.Dispatchs = make(map[uint16]DispatchInterface)
}

func (m *HandlerMsgBase) ReleaseHandler() {
	m.Dispatchs = nil
}

func (m *HandlerMsgBase) RegDispatch(dispatch_type uint16, dispatch DispatchInterface) {
	m.Dispatchs[dispatch_type] = dispatch
	common.GStdout.Info("RegDispatch type : %v", dispatch_type)
}

func (m *HandlerMsgBase) DoDispatch(msg proto.Message, session *Session) {
	strkey := reflect.TypeOf(msg).String()
	strkey = strings.Replace(strkey, "*protomsg.", "k", 1)

	msg_type, ok := msgtype.MsgType_value[strkey]

	if !ok {
		common.GStdout.Error("msg id not handler:%#v", msg)
		return
	}

	for _, v := range m.Dispatchs {
		v.Dispatch(msgtype.MsgType(msg_type), msg)
	}
}

func (m *HandlerMsgBase) GetMsgID(msg proto.Message) uint16 {
	strkey := reflect.TypeOf(msg).String()
	strkey = strings.Replace(strkey, "*protomsg.", "k", 1)

	return uint16(msgtype.MsgType_value[strkey])
}

func (m *HandlerMsgBase) GetMsgType(msg_type uint16) (bool, reflect.Type) {
	key, ok := msgtype.MsgType_name[int32(msg_type)]

	if !ok {
		common.GStdout.Error("msg_type not exist:%v", msg_type)
		return false, nil
	}

	strkey := strings.Replace(key, "k", "protomsg.", 1)

	t := proto.MessageType(strkey)

	if t == nil {
		common.GStdout.Error("msg_type1 not exist:%v,strkey:%v", msg_type, strkey)
		return false, nil
	}

	return true, t
}
