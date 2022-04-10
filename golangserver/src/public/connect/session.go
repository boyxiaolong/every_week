package connect

import (

	//"fmt"
	"public/common"
	"public/message/msgtype"
	"reflect"
	"runtime/debug"
	"sync"

	"public/link"

	"github.com/golang/protobuf/proto"
)

type Session struct {
	Lsession     *link.Session
	stopWait     sync.WaitGroup
	MsgRecv      chan MessageResponse
	MsgSend      chan proto.Message
	stopRecvChan chan struct{}
	stopSendChan chan struct{}
	bStart       bool
	once         sync.Once
	ServerType   uint64
	HandlerMsg   MessageHandler
}

type MessageResponse struct {
	msg_buf  []byte
	msg_type uint16
}

func NewSessionByMsgHander(lsession *link.Session, server_type uint64, handlerMsg MessageHandler) *Session {
	session := &Session{
		Lsession:     lsession,
		MsgRecv:      make(chan MessageResponse, 10),
		MsgSend:      make(chan proto.Message, 10),
		stopRecvChan: make(chan struct{}),
		stopSendChan: make(chan struct{}),
		bStart:       true,
		ServerType:   server_type,
		HandlerMsg:   handlerMsg,
	}

	GSessionManager.Put(server_type, session)

	return session
}

func (m *Session) ReciveMessage() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	for {
		var msg []byte
		var msg_type uint16

		if err := m.Lsession.Receive(&msg, &msg_type); err != nil {
			common.GStdout.Error("recive fail server_type %d, error %v ", m.ServerType, err)
			break
		}

		msg_package := MessageResponse{}
		msg_package.msg_buf = msg
		msg_package.msg_type = msg_type
		m.AddRecv(msg_package)
	}

	m.StopOnce()
	m.Lsession = nil
}

func (m *Session) StopOnce() {
	m.once.Do(m.Stop)
}

func (m *Session) Stop() {
	m.bStart = false
	m.stopWait.Add(2)
	close(m.stopRecvChan)
	close(m.stopSendChan)

	m.stopWait.Wait()
	m.Lsession.Close()

	common.GStdout.Error("Close server [%d] session", m.ServerType)
	GSessionManager.Remove(m.ServerType)
}

func (m *Session) Start() {
	go m.ReciveMessage()
	go m.SendMessage()
	go m.DoMessage()
}

func (m *Session) AddRecv(msg interface{}) error {
	if !m.bStart {
		return nil
	}
	select {
	case m.MsgRecv <- msg.(MessageResponse):
	}
	return nil
}

func (m *Session) AddSend(msg proto.Message) error {
	if !m.bStart {
		return nil
	}

	select {
	case m.MsgSend <- msg:
	}
	return nil
}

func (m *Session) SendMessage() {
	for {
		select {
		case msg := <-m.MsgSend:
			if m.Lsession == nil {
				return
			}
			m.Write(msg)
		case <-m.stopSendChan:
			m.stopWait.Done()
			return
		}
	}
}

func (m *Session) Write(msg proto.Message) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	data, err := proto.Marshal(msg)
	if err != nil {
		common.GStdout.Error("protobuf  编码出错:%v", err.Error())
		return
	}

	if m.Lsession == nil {
		return
	}

	if m.Lsession.IsClosed() {
		return
	}

	//GStdout.Error("write:%#v,%v", msg, data)

	msg_type := m.HandlerMsg.GetMsgID(msg)

	if msg_type == 0 {
		common.GStdout.Error("write message is zero：%#v", msg)
		return
	}

	err = m.Lsession.Send(data, int(msg_type))
}

func (m *Session) DoMessage() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	for {
		select {
		case msg := <-m.MsgRecv:
			m.InitMsg(msg.msg_type, msg.msg_buf)
		case <-m.stopRecvChan:
			m.stopWait.Done()
			return
		}
	}
}

func (m *Session) InitMsg(msg_type uint16, data []byte) (err error) {
	ok, t := m.HandlerMsg.GetMsgType(msg_type)

	if !ok {
		return
	}

	if module, ok := reflect.New(t.Elem()).Interface().(proto.Message); ok {
		err = proto.Unmarshal(data, module)
		if err != nil {
			common.GStdout.Error("protobuf 解码出错:%v", err.Error())
			return
		}

		if msg_type != uint16(msgtype.MsgType_kMsgGS2CLServerTimeNotice) {
			GMsgPrint.OnMessage(msg_type, module, len(data))
		}

		m.HandlerMsg.DoDispatch(module, m)
	}

	return
}

func (m *Session) Send(pb proto.Message) {
	m.AddSend(pb)
}
