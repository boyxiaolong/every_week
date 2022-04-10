package player

import (
	BTree "public/behaviortree"
	"public/message/protomsg"
	"strconv"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("RadarComplete", CreateRadarCompleteTaskTask)
}

type RadarCompleteTask struct {
	TempeldId uint64
}

func CreateRadarCompleteTaskTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &RadarCompleteTask{}
	return task, nil
}

func (t *RadarCompleteTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	radars := player.GetData("RadarData.Radars").([]*protomsg.RadarData)
	for _, radar := range radars {
		player.PM("completeradar " + strconv.Itoa(int(radar.RadarId)))
	}
	return true
}
