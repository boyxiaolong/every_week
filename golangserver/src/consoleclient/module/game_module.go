package module

import (
	"fmt"
	"public/common"
	"public/config"
	"public/connect"
	"public/message/msgtype"
	"public/message/protomsg"

	"github.com/golang/protobuf/proto"

	//"public/wait"
	//"sync"
	"consoleclient/handle"
	"public/wait"
	"sync"
	"time"
)

//GGameInfo GGameInfo
var GGameInfo *GameInfo

func init() {
	GGameInfo = &GameInfo{
		Online: false,
	}
	GGameInfo.WaitMessage.Init(config.Mode)
}

//GameInfo GameInfo
type GameInfo struct {
	GSession     *connect.Session
	IP           string
	Port         uint32
	PlayerID     uint64
	LoginSession string
	Account      uint64
	Online       bool
	WaitMessage  wait.Wait
}

//SetNetAddress SetNetAddress
func (m *GameInfo) SetNetAddress(ip string, port uint32) {
	m.IP = ip
	m.Port = port
}

//SetPlayerID SetPlayerID
func (m *GameInfo) SetPlayerID(playerID uint64) {
	m.PlayerID = playerID
}

//SetLoginSession SetLoginSession
func (m *GameInfo) SetLoginSession(loginSession string) {
	m.LoginSession = loginSession
}

//SetAccount SetAccount
func (m *GameInfo) SetAccount(account uint64) {
	m.Account = account
}

//NewSession NewSession
func (m *GameInfo) NewSession() error {
	port := fmt.Sprintf("%v", m.Port)
	session, err := connect.GetConnectSessionByMsgHander(m.IP, port, common.CLIENT_TYPE_GAME_SERVER, handle.GhandleMsg)

	if err != nil {
		common.GStdout.Error("game new session error:%v", err)
		return err
	}

	m.GSession = session
	m.GSession.Start()

	return nil
}

//Login Login
func (m *GameInfo) Login() {
	msgCL2GSLoginRequest := &protomsg.MsgCL2GSLoginRequest{}
	msgCL2GSLoginRequest.PlayerId = m.PlayerID
	msgCL2GSLoginRequest.Account = m.Account
	msgCL2GSLoginRequest.LoginSession = m.LoginSession
	msgCL2GSLoginRequest.Udid = "fake udid"
	msgCL2GSLoginRequest.ClientVersion = "0.0.0"

	msgGS2CLLoginReply := m.SendAndWait(
		msgCL2GSLoginRequest,
		msgtype.MsgType_kMsgGS2CLLoginReply).(*protomsg.MsgGS2CLLoginReply)

	common.CHECK_ERROR(int32(msgGS2CLLoginReply.ErrorCode), uint16(msgtype.MsgType_kMsgGS2CLLoginReply))

	m.LoginReady()
}

//LoginReady LoginReady
func (m *GameInfo) LoginReady() {
	msgCL2GSEnterGameRequest := &protomsg.MsgCL2GSEnterGameRequest{}

	msgGS2CLEnterGameReply := m.SendAndWait(
		msgCL2GSEnterGameRequest,
		msgtype.MsgType_kMsgGS2CLEnterGameReply).(*protomsg.MsgGS2CLEnterGameReply)

	m.Online = true
	go m.KeepLive()
	m.GSession.Lsession.AddCloseCallback(GGameInfo.GSession, func() {
		m.Online = false
	})

	common.CHECK_ERROR(int32(msgGS2CLEnterGameReply.ErrorCode), uint16(msgtype.MsgType_kMsgGS2CLEnterGameReply))

	if config.Mode == common.MODE_EXIT {
		go GTestMgr.TestAllModule()
	}

}

//KeepLive 心跳
func (m *GameInfo) KeepLive() {
	if m.Online {
		//common.GStdout.Info("Keep live")
		request := &protomsg.MsgCL2GSKeepLiveRequest{}
		GGameInfo.Send(request)
		time.AfterFunc(time.Second*15, m.KeepLive)
	}
}

//Send comment
func (m *GameInfo) Send(pb proto.Message) {
	m.GSession.AddSend(pb)
}

//SendAndWait comment
func (m *GameInfo) SendAndWait(pb proto.Message, msgType msgtype.MsgType) proto.Message {
	m.GSession.AddSend(pb)
	waitMessage := m.WaitMessage.DoWait(uint16(msgType), wait.DefaultWaitSeconds)
	return waitMessage
}

//SendAndWaitBySesision comment
func (m *GameInfo) SendAndWaitBySesision(session *connect.Session, pb proto.Message, msgType msgtype.MsgType) proto.Message {
	session.AddSend(pb)
	waitMessage := m.WaitMessage.DoWait(uint16(msgType), wait.DefaultWaitSeconds)
	return waitMessage
}

//SendAndWaitMultiple comment
func (m *GameInfo) SendAndWaitMultiple(pb proto.Message, args ...msgtype.MsgType) []proto.Message {
	m.GSession.AddSend(pb)
	return m.WaitMessage.WaitMultiple(args...)
}

//WaitMultiple comment
func (m *GameInfo) WaitMultiple(args ...msgtype.MsgType) []proto.Message {
	var syncWait sync.WaitGroup
	messages := make([]proto.Message, len(args))
	syncWait.Add(len(args))
	for k, v := range args {
		go func(i int, msg_type msgtype.MsgType) {
			messages[i] = m.WaitMessage.DoWait(uint16(msg_type), wait.DefaultWaitSeconds)
			syncWait.Done()
		}(k, v)
	}

	syncWait.Wait()
	return messages
}

//Wait comment
func (m *GameInfo) Wait(msgType msgtype.MsgType) proto.Message {
	return m.WaitMessage.DoWait(uint16(msgType), wait.DefaultWaitSeconds)
}

//WaitSeconds comment
func (m *GameInfo) WaitSeconds(msgType msgtype.MsgType, seconds uint16) proto.Message {
	return m.WaitMessage.DoWait(uint16(msgType), seconds)
}

//Break comment
func (m *GameInfo) Break() {
	m.WaitMessage.SetBreakStatus(1)
}
