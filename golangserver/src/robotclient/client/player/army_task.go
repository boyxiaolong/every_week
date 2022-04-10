package player

import (
	"fmt"
	BTree "public/behaviortree"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("InstantTrain", CreateInstantTrainTask)
	BTree.RegisterConditionCallback("IsRich", IsRich)
}

// IsRich comment
func IsRich(testobj BTree.ObjectInterface, params string) bool {
	return true
}

// CreateInstantTrainTask comment
func CreateInstantTrainTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("Instant Train Task params error!")
	}

	ArmyId, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	Count, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	task := &InstantTrainTask{ArmyId: (uint32)(ArmyId), Count: (uint64)(Count)}
	return task, nil
}

// InstantTrainTask comment
type InstantTrainTask struct {
	ArmyId uint32
	Count  uint64
}

// DoTask comment
func (t *InstantTrainTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.InstantTrain(player)
}

// InstantTrain comment
func (t *InstantTrainTask) InstantTrain(player *Player) bool {

	emoney := player.GetData("BaseData.Player_Emoney").(uint64)
	if emoney < 1000 {
		player.PM("addcurrency 1 20000")
	}

	request := &protomsg.MsgCL2GSInstantTrainRequest{}
	request.Army = &protomsg.ArmyInfo{}
	request.Army.ArmyId = t.ArmyId
	request.Army.Count = t.Count

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLInstantTrainReply).(*protomsg.MsgGS2CLInstantTrainReply)
	if !ok {
		return false
	}

	common.GStdout.Info("=================== army id: %v, count: %v", request.Army.ArmyId, request.Army.Count)

	return int32(response.ErrorCode) == 0
}
