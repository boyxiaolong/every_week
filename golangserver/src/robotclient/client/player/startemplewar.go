package player

import (
	BTree "public/behaviortree"
	"public/common"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("StartTempleWar", CreateStartTempleWarTask)
}

type StartTempleWarTask struct {
	TempeldId uint64
}

func CreateStartTempleWarTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &StartTempleWarTask{}
	temple_id, error := strconv.ParseUint(params[0], 10, 64)
	if error != nil {
		return nil, error
	}

	task.TempeldId = temple_id
	return task, nil
}

func (t *StartTempleWarTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)

	request := &protomsg.MsgCL2GSStartTempleWarRequest{}
	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLStartTempleWarReply).(*protomsg.MsgGS2CLStartTempleWarReply)
	if !ok {
		return false
	}

	if response.ErrorCode != 0 {
		common.GStdout.Debug("start templewar error:%v", response.ErrorCode)
		return false
	}
	return true
}
