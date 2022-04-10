package player

import (
	"fmt"
	"public/common"
	"public/config"
	"public/connect"
	"public/message/msgtype"
	"public/message/protomsg"
	"public/wait"
	"robotclient/client/player/data"
	"robotclient/loadconfig"
	"strconv"
	"strings"
	"sync"
	"time"

	BTree "public/behaviortree"

	"public/task"

	"github.com/golang/protobuf/proto"
)

func init() {
}

type HandleMsg struct {
	connect.HandlerMsgBase
}

// Player comment
type Player struct {
	connect.BaseDispatch
	GSession    *connect.Session
	PlayerID    uint64
	Account     uint64
	Online      bool
	Deploys     map[uint32][]uint32 // deploy_id, hero_id_list
	CurTask     *BTree.BTreeRunTask
	WaitMessage wait.Wait
	HandleMsg   *HandleMsg
	data_center *data.DataCenter
	main_task   *task.Task
	wait_task   *task.Task
}

// Init comment
func (p *Player) Init(playerid uint64, account uint64) {
	p.PlayerID = playerid
	p.Account = account
	p.WaitMessage.Init(config.Mode)
	p.data_center = data.CreateDataCenter()
	p.data_center.Start()
	p.HandleMsg = &HandleMsg{}
	p.HandleMsg.Init()
	p.HandleMsg.RegDispatch(1, p.data_center)
	p.HandleMsg.RegDispatch(2, p)
	p.Online = false
	p.Deploys = make(map[uint32][]uint32)
	p.CurTask = nil
	p.main_task = task.MakeTask(1000)
	p.wait_task = task.MakeTask(1000)
	p.Start()
}

// Init comment
func (p *Player) OnOffline() {
	p.data_center.Stop()
	p.data_center.Reset()
	p.Online = false
}

func (p *Player) Release() {
	p.Online = false
	p.HandleMsg.ReleaseHandler()
	if p.GSession != nil {
		p.GSession.StopOnce()
	}
	p.StopBTreeRun()
	p.StopRunning()
}

//Dispatch comment
func (p *Player) Dispatch(msgtype msgtype.MsgType, msg proto.Message) {
	p.wait_task.AddTask(func() {
		p.WaitMessage.Done(uint16(msgtype), msg)
	})
}

//GetObjectId GetObjectId
func (p *Player) GetObjectId() uint64 {
	return p.PlayerID
}

//SetPlayerID SetPlayerID
func (p *Player) SetPlayerID(playerID uint64) {
	p.PlayerID = playerID
}

//SetAccount SetAccount
func (p *Player) SetAccount(account uint64) {
	p.Account = account
}

//NewLoginSession NewLoginSession
func (p *Player) NewLoginSession() *connect.Session {
	session, err := connect.GetConnectSessionByMsgHander(loadconfig.GetIp(), loadconfig.GetPort(), common.CLIENT_TYPE_LOGIN_SERVER, p.HandleMsg)

	if err != nil {
		common.GStdout.Error("login new session error:%v", err)
		return nil
	}
	return session
}

//NewGSession NewGSession
func (p *Player) NewGameSession(ip string, port_id uint32) error {
	port := fmt.Sprintf("%v", port_id)
	session, err := connect.GetConnectSessionByMsgHander(ip, port, common.CLIENT_TYPE_GAME_SERVER, p.HandleMsg)

	if err != nil {
		common.GStdout.Error("game new session error:%v", err)
		return err
	}

	p.GSession = session
	p.GSession.Start()
	return nil
}

func (p *Player) DataHander() {

}

//Dispatch comment
func (p *Player) Start() {
	if p.main_task.IsRunning {
		return
	}
	p.main_task.Start()

	if p.wait_task.IsRunning {
		return
	}
	p.wait_task.Start()
}

func (p *Player) StopRunning() {
	if !p.main_task.IsRunning {
		return
	}
	p.main_task.Stop()

	if !p.wait_task.IsRunning {
		return
	}
	p.wait_task.Stop()
}

func (p *Player) StartTest(btree *BTree.BTree) {
	if p.CurTask != nil {
		p.CurTask.IsRun = false
		p.CurTask = nil
	}

	new_task := BTree.CreateRunTask(p)
	p.CurTask = new_task

	p.main_task.AddTask(func() {
		btree.Do(new_task)
	})
}

func (p *Player) StopBTreeRun() {
	if p.CurTask != nil {
		common.GStdout.Debug("StopBTreeRun ======================== ")
		p.CurTask.IsRun = false
		p.CurTask = nil
	}
}

//KeepLive 心跳
func (p *Player) KeepLive() {
	if p.Online {
		//common.GStdout.Debug("Keep live")
		request := &protomsg.MsgCL2GSKeepLiveRequest{}
		p.Send(request)
		time.AfterFunc(time.Second*15, p.KeepLive)
	}
}

//Send comment
func (p *Player) Send(pb proto.Message) {
	p.GSession.AddSend(pb)
}

//SendAndWait comment
func (p *Player) SendAndWait(pb proto.Message, msgType msgtype.MsgType) proto.Message {
	p.GSession.AddSend(pb)
	waitMessage := p.WaitMessage.DoWait(uint16(msgType), wait.DefaultWaitSeconds)
	return waitMessage
}

//SendAndWaitBySesision comment
func (p *Player) SendAndWaitBySesision(session *connect.Session, pb proto.Message, msgType msgtype.MsgType) proto.Message {
	session.AddSend(pb)
	waitMessage := p.WaitMessage.DoWait(uint16(msgType), wait.DefaultWaitSeconds)
	return waitMessage
}

//SendAndWaitMultiple comment
func (p *Player) SendAndWaitMultiple(pb proto.Message, args ...msgtype.MsgType) []proto.Message {
	p.GSession.AddSend(pb)
	return p.WaitMessage.WaitMultiple(args...)
}

//WaitMultiple comment
func (p *Player) WaitMultiple(args ...msgtype.MsgType) []proto.Message {
	var syncWait sync.WaitGroup
	messages := make([]proto.Message, len(args))
	syncWait.Add(len(args))
	for k, v := range args {
		go func(i int, msg_type msgtype.MsgType) {
			messages[i] = p.WaitMessage.DoWait(uint16(msg_type), wait.DefaultWaitSeconds)
			syncWait.Done()
		}(k, v)
	}

	syncWait.Wait()
	return messages
}

//Wait comment
func (p *Player) Wait(msgType msgtype.MsgType) proto.Message {
	return p.WaitMessage.DoWait(uint16(msgType), wait.DefaultWaitSeconds)
}

//Break comment
func (p *Player) Break() {
	p.WaitMessage.SetBreakStatus(1)
}

//PM 执行PM命令
func (p *Player) PM(cmd string) (err error) {

	cmd = strings.Replace(cmd, "\n", "", -1)
	cmd = strings.Replace(cmd, "\r", "", -1)

	str := &common.StringParse{}
	str.ParseString(cmd, " ")

	if str.Len() < 1 {
		common.GStdout.Error("Invalid PM Command")
		return
	}

	if !p.Online {
		common.GStdout.Error("user not online")
		return
	}

	params := make([]string, 0)
	params = append(params, str.Strs...)

	cmd = strings.Join(params, " ")

	common.GStdout.Debug("pm:%v", cmd)
	msgCL2GSPMCommandRequest := &protomsg.MsgCL2GSPMCommandRequest{}
	msgCL2GSPMCommandRequest.Command = cmd

	response, ok := p.SendAndWait(msgCL2GSPMCommandRequest, msgtype.MsgType_kMsgGS2CLPMCommandReply).(*protomsg.MsgGS2CLPMCommandReply)
	if !ok {
		common.GStdout.Debug("pm command send wait timeout error")
		return
	}

	if response.ErrorCode != 0 {
		common.GStdout.Debug("pm command request error:%v", response.ErrorCode)
		return
	}
	return
}

func (p *Player) DoCreateBuilding(constructionType uint32) (uint32, int32) {
	msgCL2GSPMCommandRequest := &protomsg.MsgCL2GSPMCommandRequest{}
	msgCL2GSPMCommandRequest.Command = "createbuilding " + strconv.FormatUint(uint64(constructionType), 10) + " 1 1 1"
	if reply, ok := p.SendAndWait(msgCL2GSPMCommandRequest, msgtype.MsgType_kMsgGS2CLCityCreateBuildingReply).(*protomsg.MsgGS2CLCityCreateBuildingReply); ok {
		if reply.ErrorCode != 0 {
			return 0, reply.ErrorCode
		}

		return reply.BuildingId, 0
	}
	return 0, 1
}

func (p *Player) GetData(data_name string) interface{} {
	return p.data_center.GetData(data_name)
}

func (p *Player) SetData(data_name string, value interface{}) {
	p.data_center.SetData(data_name, value)
}

func (p *Player) GetDeployInfoId() *uint32 {
	request := &protomsg.MsgGS2CLGetDeployInfoRequest{}
	response, ok := p.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGetDeployInfoReply).(*protomsg.MsgGS2CLGetDeployInfoReply)
	if !ok {
		return nil
	}

	for _, deploy := range response.Infos {
		if !deploy.IsMarch {
			return &deploy.Id
		}
	}
	return nil
}
