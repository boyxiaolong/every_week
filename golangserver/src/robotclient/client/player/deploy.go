package player

import (
	"fmt"
	BTree "public/behaviortree"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

var TaskParamsCount int

// Init comment
func init() {
	BTree.RegisterTaskCreator("GetDeploy", CreateGetDeployTask)
	BTree.RegisterTaskCreator("SetDeploy1", CreateSetDeployTask)
	BTree.RegisterTaskCreator("SetDeploy2", CreateSetDeployTask)
	BTree.RegisterTaskCreator("SetDeploy3", CreateSetDeployTask)
	BTree.RegisterTaskCreator("SetDeploy4", CreateSetDeployTask)
	BTree.RegisterTaskCreator("SetDeploy5", CreateSetDeployTask)
	TaskParamsCount = 12
}

type GetDeployTask struct {
}

func CreateGetDeployTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &GetDeployTask{}
	return task, nil
}

func (t *GetDeployTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	request := &protomsg.MsgGS2CLGetDeployInfoRequest{}
	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLGetDeployInfoReply).(*protomsg.MsgGS2CLGetDeployInfoReply)
	if !ok {
		return false
	}

	for _, deploy_info := range response.Infos {
		hero_id_list := make([]uint32, 5)
		for _, general := range deploy_info.GetGenerals() {
			hero_id_list = append(hero_id_list, general.GeneralId)
		}

		player.Deploys[deploy_info.Id] = hero_id_list
	}
	return true
}

// LoginTask comment
type SetDeployTask struct {
	Infos protomsg.DeployInfo
}

func CreateSetDeployTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) != TaskParamsCount {
		return nil, fmt.Errorf("set deploy task params error")
	}

	task := &SetDeployTask{}

	id, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	main_hero_id, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	position0, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	second_hero_id, error := strconv.Atoi(params[3])
	if error != nil {
		return nil, error
	}

	position1, error := strconv.Atoi(params[4])
	if error != nil {
		return nil, error
	}

	third_hero_id, error := strconv.Atoi(params[5])
	if error != nil {
		return nil, error
	}

	position2, error := strconv.Atoi(params[6])
	if error != nil {
		return nil, error
	}

	fourth_hero_id, error := strconv.Atoi(params[7])
	if error != nil {
		return nil, error
	}

	position3, error := strconv.Atoi(params[8])
	if error != nil {
		return nil, error
	}

	fifth_hero_id, error := strconv.Atoi(params[9])
	if error != nil {
		return nil, error
	}

	position4, error := strconv.Atoi(params[10])
	if error != nil {
		return nil, error
	}

	siege_id, error := strconv.Atoi(params[11])
	if error != nil {
		return nil, error
	}

	Generals1 := &protomsg.GeneralInfo{}
	Generals1.GeneralId = (uint32)(main_hero_id)
	Generals1.Position = (protomsg.DeployPosition)(position0)

	Generals2 := &protomsg.GeneralInfo{}
	Generals2.GeneralId = (uint32)(second_hero_id)
	Generals2.Position = (protomsg.DeployPosition)(position1)

	Generals3 := &protomsg.GeneralInfo{}
	Generals3.GeneralId = (uint32)(third_hero_id)
	Generals3.Position = (protomsg.DeployPosition)(position2)

	Generals4 := &protomsg.GeneralInfo{}
	Generals4.GeneralId = (uint32)(fourth_hero_id)
	Generals4.Position = (protomsg.DeployPosition)(position3)

	Generals5 := &protomsg.GeneralInfo{}
	Generals5.GeneralId = (uint32)(fifth_hero_id)
	Generals5.Position = (protomsg.DeployPosition)(position4)

	task.Infos.Id = (uint32)(id)
	task.Infos.Generals = append(task.Infos.Generals, Generals1)
	task.Infos.Generals = append(task.Infos.Generals, Generals2)
	task.Infos.Generals = append(task.Infos.Generals, Generals3)
	task.Infos.Generals = append(task.Infos.Generals, Generals4)
	task.Infos.Generals = append(task.Infos.Generals, Generals5)
	task.Infos.SiegeId = (uint32)(siege_id)
	task.Infos.IsAutoFill = true
	task.Infos.IsMarch = false
	return task, nil
}

func (t *SetDeployTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	_, ok := player.Deploys[t.Infos.Id]
	if ok {
		common.GStdout.Debug("deploy already exist id:%v", t.Infos.Id)
		return true
	} else {
		hero_id_list := make([]uint32, 5)
		for _, general := range t.Infos.GetGenerals() {
			hero_id_list = append(hero_id_list, general.GeneralId)
		}

		player.Deploys[t.Infos.Id] = hero_id_list
		request := &protomsg.MsgGS2CLSetDeployInfoRequest{}
		request.Infos = append(request.Infos, &(t.Infos))

		response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLSetDeployInfoReply).(*protomsg.MsgGS2CLSetDeployInfoReply)
		if !ok {
			return false
		}

		if response.ErrorCode != 0 {
			common.GStdout.Debug("set deploy request error:%v", response.ErrorCode)
			return false
		}
	}
	return true
}
