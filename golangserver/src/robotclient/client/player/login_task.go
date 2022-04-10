package player

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"math/rand"
	BTree "public/behaviortree"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("Login", CreateLoginTask)
	BTree.RegisterTaskCreator("Logout", CreateLogoutTask)
	BTree.RegisterTaskCreator("PM", CreatePMTask)
	BTree.RegisterTaskCreator("End", CreateEndTask)
	BTree.RegisterTaskCreator("RandomSleep", CreateRandomSleepTask)
	BTree.RegisterConditionCallback("IsOnline", IsOnline)
}

func IsOnline(testobj BTree.ObjectInterface, params string) bool {
	player := testobj.(*Player)
	//common.GStdout.Console("player online %v", player.Online)
	return player.Online
}

func CreateLoginTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("login task params error")
	}
	task := &LoginTask{}
	kingdomID, _ := strconv.ParseUint("params", 10, 32)
	task.kingdomID = (uint32)(kingdomID)
	return task, nil
}

func CreateLogoutTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &LogoutTask{}
	return task, nil
}

func CreatePMTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("pm task params error")
	}
	task := &PMTask{}
	task.pm = params[0]
	return task, nil
}

func CreateEndTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &EndTask{}
	return task, nil
}

func CreateRandomSleepTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("random sleep task params error")
	}

	minTime, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	maxTime, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	task := &RandomSleepTask{minTime: minTime, maxTime: maxTime}
	return task, nil
}

// LoginTask comment
type LoginTask struct {
	kingdomID uint32
}

// LogoutTask comment
type LogoutTask struct {
}

// EndTask comment
type EndTask struct {
}

// PMTask comment
type PMTask struct {
	pm string
}

// RandomSleepTask comment
type RandomSleepTask struct {
	minTime int
	maxTime int
}

func (t *LoginTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.Login(player)
}

func (t *LoginTask) Login(player *Player) bool {
	defer func() {
		if err := recover(); err != nil {
			//debug.PrintStack()
			common.GStdout.Error("%v", string(debug.Stack()))

			if player.GSession != nil {
				player.GSession.StopOnce()
			}
			player.Online = false
		}
	}()

	ip, port, loginsession, ok := t.LoginToLoginServer(player, t.kingdomID)
	if !ok {
		common.GStdout.Debug("login to login server error")
		return false
	}
	if !t.LoginToGameServer(player, ip, port, loginsession) {
		common.GStdout.Debug("login to game server error")
		return false
	}
	return true
}

func (t *LoginTask) LoginToLoginServer(player *Player, kingdomID uint32) (string, uint32, string, bool) {
	session := player.NewLoginSession()

	if session == nil {
		return "", 0, "", false
	}
	defer session.StopOnce()

	session.Start()

	msgCL2LSLoginRequest := &protomsg.MsgCL2LSLoginRequest{}
	msgCL2LSLoginRequest.Account = player.Account
	msgCL2LSLoginRequest.StrMacAddress = "00:00:00:00:00:00"
	msgCL2LSLoginRequest.WebSession = "sesskonkey"
	msgCL2LSLoginRequest.AreaId = 1
	msgCL2LSLoginRequest.GameId = "1234"
	msgCL2LSLoginRequest.CampId = 1

	msgLS2CLLoginReply, ok := player.SendAndWaitBySesision(session,
		msgCL2LSLoginRequest,
		msgtype.MsgType_kMsgLS2CLLoginReply).(*protomsg.MsgLS2CLLoginReply)

	if !ok {
		return "", 0, "", false
	}

	if msgLS2CLLoginReply.ErrorCode != 0 {
		common.GStdout.Error("login fail errorcode:%v", msgLS2CLLoginReply)
		return "", 0, "", false
	}

	return msgLS2CLLoginReply.NetAddress.Ip, msgLS2CLLoginReply.NetAddress.Port, msgLS2CLLoginReply.LoginSession, true
}

//Login Login
func (t *LoginTask) LoginToGameServer(player *Player, ip string, port uint32, loginsession string) bool {
	err := player.NewGameSession(ip, port)
	if err != nil {
		common.GStdout.Error("login fail errorcode:%v", err)
		return false
	}

	msgCL2GSLoginRequest := &protomsg.MsgCL2GSLoginRequest{}
	msgCL2GSLoginRequest.PlayerId = player.PlayerID
	msgCL2GSLoginRequest.Account = player.Account
	msgCL2GSLoginRequest.LoginSession = loginsession
	msgCL2GSLoginRequest.Udid = "fake udid"
	msgCL2GSLoginRequest.ClientVersion = "0.0.0"

	msgGS2CLLoginReply, ok := player.SendAndWait(
		msgCL2GSLoginRequest,
		msgtype.MsgType_kMsgGS2CLLoginReply).(*protomsg.MsgGS2CLLoginReply)

	if !ok {
		return false
	}

	if msgGS2CLLoginReply.ErrorCode != 0 {
		return false
	}

	msgCL2GSEnterGameRequest := &protomsg.MsgCL2GSEnterGameRequest{}

	msgGS2CLEnterGameReply, ok := player.SendAndWait(
		msgCL2GSEnterGameRequest,
		msgtype.MsgType_kMsgGS2CLEnterGameReply).(*protomsg.MsgGS2CLEnterGameReply)

	if !ok {
		return false
	}

	player.Online = true
	//common.GStdout.Console("set player online %v", player.Online)

	go player.KeepLive()
	player.GSession.Lsession.AddCloseCallback(player.GSession, func() {
		player.OnOffline()
	})

	if msgGS2CLEnterGameReply.ErrorCode != 0 {
		return false
	}

	//player.PM("superman")
	time.Sleep(2 * time.Second) // 停2秒,等数据初始化
	return true
}

func (t *LogoutTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.Logout(player)
}

//Logout LogoutTask
func (t *LogoutTask) Logout(player *Player) bool {
	if player.GSession != nil {
		player.GSession.StopOnce()
	}
	player.Online = false
	return true
}

func (t *PMTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	player.PM(t.pm)
	return true
}

func (t *EndTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	endPlayer := runtask.TaskObj.(*Player)
	runtask.StopRun()
	GPlayerMgr.RemovePlayer(endPlayer)
	return true
}

func (t *RandomSleepTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	randTime := t.minTime + rand.Intn((int)(t.maxTime-t.minTime))
	if randTime > 0 {
		time.Sleep(time.Duration(randTime) * time.Second)
	}
	return true
}
