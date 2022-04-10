package module

import (
	"consoleclient/loadconfig"
	"public/command"
	"public/common"
	"public/connect"
	"public/message/error_code"
	"public/message/protomsg"

	//. "game/wait"
	"consoleclient/handle"
	"public/message/msgtype"
)

func init() {
	command.GCommand.RegCommand("login", Login, "Login")

	//RegResHandler(uint16(MsgType_kMsgLS2CLLoginReply), LoginResponse)
	//RegResHandler(uint16(MsgType_kMsgGS2CLLoginReply), GameLoginResponse)
	//RegResHandler(uint16(MsgType_kMsgGS2CLEnterGameReply), GameEnterResponse)
}

//Login Login
func Login(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("login param error")
		return
	}

	account := uint64(str.GetInt(1))
	kingdomID := uint32(str.GetInt(2))
	if kingdomID == 0 {
		kingdomID = 1
	}
	LoginToServer(account)

	return
}

//LoginToServer LoginToServer
func LoginToServer(account uint64) {
	session := GetNewSession()

	if session == nil {
		return
	}

	session.Start()

	msgCL2LSLoginRequest := &protomsg.MsgCL2LSLoginRequest{}
	msgCL2LSLoginRequest.Account = account
	msgCL2LSLoginRequest.StrMacAddress = "00:00:00:00:00:00"
	msgCL2LSLoginRequest.WebSession = "sesskonkey"
	msgCL2LSLoginRequest.AreaId = 1
	msgCL2LSLoginRequest.GameId = "1234"
	msgCL2LSLoginRequest.CampId = 1
	msgCL2LSLoginRequest.Language = 1

	msgLS2CLLoginReply, ok := GGameInfo.SendAndWaitBySesision(session,
		msgCL2LSLoginRequest,
		msgtype.MsgType_kMsgLS2CLLoginReply).(*protomsg.MsgLS2CLLoginReply)

	if !ok {
		session.StopOnce()
		return
	}

	if msgLS2CLLoginReply.ErrorCode != 0 {
		common.GStdout.Error("login fail errorcode:%v", msgLS2CLLoginReply)
		return
	}

	GGameInfo.SetNetAddress(msgLS2CLLoginReply.NetAddress.Ip, msgLS2CLLoginReply.NetAddress.Port)
	GGameInfo.SetPlayerID(msgLS2CLLoginReply.PlayerId)
	GGameInfo.SetLoginSession(msgLS2CLLoginReply.LoginSession)
	session.StopOnce()
	GGameInfo.SetAccount(account)

	err := GGameInfo.NewSession()

	if err != nil {
		return
	}

	GGameInfo.Login()
}

//GetNewSession GetNewSession
func GetNewSession() *connect.Session {
	session, err := connect.GetConnectSessionByMsgHander(loadconfig.GetIp(), loadconfig.GetPort(), common.CLIENT_TYPE_LOGIN_SERVER, handle.GhandleMsg)

	if err != nil {
		common.GStdout.Error("login new session error:%v", err)
		return nil
	}

	return session
}

//LoginResponse LoginResponse
func LoginResponse(msg interface{}) (err error) {
	msgReply := msg.(*protomsg.MsgLS2CLLoginReply)

	GGameInfo.SetNetAddress(msgReply.NetAddress.Ip, msgReply.NetAddress.Port)
	GGameInfo.SetPlayerID(msgReply.PlayerId)
	GGameInfo.SetLoginSession(msgReply.LoginSession)

	return
}

//GameLoginResponse GameLoginResponse
func GameLoginResponse(msg interface{}) (err error) {
	msgReply := msg.(*protomsg.MsgGS2CLLoginReply)
	common.GStdout.Success("game login result:%v", msgReply)

	if msgReply.ErrorCode == int32(error_code.ErrorCode_kECSuccess) {
		common.GStdout.Success("game login result:%v", msgReply)
	} else {
		common.GStdout.Error("game login result:%v", msgReply)
	}

	return
}

//GameEnterResponse GameEnterResponse
func GameEnterResponse(msg interface{}) (err error) {
	msgReply := msg.(*protomsg.MsgGS2CLEnterGameReply)
	common.GStdout.Info("game enter result:%v", msgReply)

	if msgReply.ErrorCode == int32(error_code.ErrorCode_kECSuccess) {
		common.GStdout.Success("game enter result:%v", msgReply.ErrorCode)
	} else {
		common.GStdout.Error("game enter result:%v", msgReply.ErrorCode)
	}

	return
}

//NewPlayerNotice NewPlayerNotice
func NewPlayerNotice(msg interface{}) (err error) {
	common.GStdout.Success("I am a new player")
	return
}
