package connect

import (
	"public/message/msgtype"
	"runtime/debug"

	"github.com/golang/protobuf/proto"
)

type InitCall func(msg proto.Message)

type MsgData struct {
	msgtype msgtype.MsgType
	msg     proto.Message
}

type DispatchInterface interface {
	AddMsg(msgtype msgtype.MsgType, addmsg proto.Message)
	RegDispatch(msgtype msgtype.MsgType, call InitCall)
	Dispatch(msgtype msgtype.MsgType, msg proto.Message)
	DoDispatch()
}

type BaseDispatch struct {
	CallBacks map[msgtype.MsgType]InitCall
	MsgChan   chan *MsgData
}

func (m *BaseDispatch) Init() {
	m.CallBacks = make(map[msgtype.MsgType]InitCall)
	m.MsgChan = make(chan *MsgData, 100)
}

func (m *BaseDispatch) AddMsg(msgtype msgtype.MsgType, addmsg proto.Message) {
	MsgData := &MsgData{
		msgtype,
		addmsg,
	}

	m.MsgChan <- MsgData
}

func (m *BaseDispatch) RegDispatch(msgtype msgtype.MsgType, call InitCall) {
	m.CallBacks[msgtype] = call
	//common.GStdout.Console("msg_response handler[msg_id]:%v", msgtype)
}

func (m *BaseDispatch) Dispatch(msgtype msgtype.MsgType, msg proto.Message) {
	if call, ok := m.CallBacks[msgtype]; ok {
		call(msg)
	}
}

func (m *BaseDispatch) DoDispatch() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	for {
		select {
		case msgInfo := <-m.MsgChan:
			m.Dispatch(msgInfo.msgtype, msgInfo.msg)
		}
	}
}
