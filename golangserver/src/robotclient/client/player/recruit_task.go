package player

import (
	"fmt"
	"math/rand"
	BTree "public/behaviortree"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("Recruit", CreateRecruitTask)
}

var RecruitItems = [2]uint32{50000, 50001}

const RandomRecruitItemCount = 1

// CreateRecruitTask comment
func CreateRecruitTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("Recruit Task params error!")
	}

	RecruitType, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	if RecruitType > len(RecruitItems) || RecruitType == 0 {
		return nil, fmt.Errorf("Recruit Task params error!")
	}

	RecruitItem := RecruitItems[RecruitType-1]
	task := &RecruitTask{RecruitType: (protomsg.RecruitType)(RecruitType), RecruitItem: RecruitItem}
	return task, nil
}

// RecruitTask comment
type RecruitTask struct {
	RecruitType protomsg.RecruitType
	RecruitItem uint32
}

// DoTask comment
func (t *RecruitTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.Recruit(player)
}

// Recruit comment
func (t *RecruitTask) Recruit(player *Player) bool {

	random_count := rand.Intn(RandomRecruitItemCount) + 1
	pm := fmt.Sprintf("additem %v %v", t.RecruitItem, random_count)
	player.PM(pm)

	request := &protomsg.MsgCL2GSRecruitRequest{}
	request.Type = t.RecruitType
	request.IsAll = true

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLRecruitReply).(*protomsg.MsgGS2CLRecruitReply)
	if !ok {
		return false
	}
	return int32(response.ErrorCode) == 0
}
