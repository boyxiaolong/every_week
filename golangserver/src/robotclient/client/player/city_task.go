package player

import (
	"fmt"
	BTree "public/behaviortree"
	"public/message/msgtype"
	"public/message/protomsg"
	"strconv"
	"math/rand"
	//"public/common"
)

// Init comment
func init() {
	BTree.RegisterTaskCreator("QueryCity", CreateQueryCityTask)
	BTree.RegisterTaskCreator("CopyCityTheme", CreateCopyCityThemeTask)
	BTree.RegisterTaskCreator("ApplyTheme", CreateApplyThemeTask)
}

// RandomPlayerId comment
func RandomPlayerId( startPlayerId uint64 , endPlayerId uint64 ) uint64 {
	if startPlayerId > endPlayerId {
		startPlayerId,endPlayerId = endPlayerId,startPlayerId
	}
	var rangeValue = endPlayerId - startPlayerId
	var randomValue uint64
	if rangeValue != 0 {
		randomValue = rand.Uint64() % rangeValue
	}
	return startPlayerId + randomValue
}

// QueryCityTask comment
type QueryCityTask struct {
	startPlayerId uint64
	endPlayerId   uint64
}

// CreateQueryCityTask comment
func CreateQueryCityTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("query city params error!")
	}

	startPlayerId, error := strconv.ParseUint(params[0], 10, 64)
	if error != nil {
		return nil, error
	}

	endPlayerId, error := strconv.ParseUint(params[1], 10, 64)
	if error != nil {
		return nil, error
	}

	if startPlayerId > endPlayerId {
		return nil, fmt.Errorf("copy City Task Theme params error!")
	}

	task := &QueryCityTask{startPlayerId: startPlayerId, endPlayerId: endPlayerId}
	return task, nil
}

// DoTask comment
func (t *QueryCityTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.QueryCity(player)
}

// QueryCity comment
func (t *QueryCityTask) QueryCity(player *Player) bool {

	request := &protomsg.MsgCL2GSQueryKingdomPlayerCityInfoRequest{}
	request.PlayerId = RandomPlayerId(t.startPlayerId,t.endPlayerId)
	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLQueryKingdomPlayerCityInfoReply).(*protomsg.MsgGS2CLQueryKingdomPlayerCityInfoReply)
	if !ok {
		return false
	}

	return int32(response.ErrorCode) == 0
}

// CopyCityThemeTask comment
type CopyCityThemeTask struct {
	startPlayerId uint64
	endPlayerId   uint64
	themeIndex    uint32
}

// CreateCopyCityThemeTask comment
func CreateCopyCityThemeTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("copy city theme params error!")
	}

	startPlayerId, error := strconv.ParseUint(params[0], 10, 64)
	if error != nil {
		return nil, error
	}

	endPlayerId, error := strconv.ParseUint(params[1], 10, 64)
	if error != nil {
		return nil, error
	}

	if startPlayerId > endPlayerId {
		return nil, fmt.Errorf("copy City Task Theme params error!")
	}

	themeIndex, error := strconv.Atoi(params[2])
	if error != nil {
		return nil, error
	}

	task := &CopyCityThemeTask{startPlayerId: startPlayerId, endPlayerId: endPlayerId,themeIndex:(uint32)(themeIndex)}
	return task, nil
}

// DoTask comment
func (t *CopyCityThemeTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.CopyCityTheme(player)
}

// CopyCityTheme comment
func (t *CopyCityThemeTask) CopyCityTheme(player *Player) bool {
	request := &protomsg.MsgCL2GSCityCopyThemeRequest{}
	request.TargetPlayerId = RandomPlayerId(t.startPlayerId,t.endPlayerId)
	request.ThemeIndex = t.themeIndex

	response, ok := player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityCopyThemeReply).(*protomsg.MsgGS2CLCityCopyThemeReply)
	if !ok {
		return false
	}

	if int32(response.ErrorCode) != 0 {
		return false
	}

	updateRequest := &protomsg.MsgCL2GSCityUpdateThemeRequest{}
	updateRequest.ThemeIndex = t.themeIndex
	updateRequest.Theme = response.Theme

	updateResponse, ok := player.SendAndWait(updateRequest, msgtype.MsgType_kMsgGS2CLCityUpdateThemeReply).(*protomsg.MsgGS2CLCityUpdateThemeReply)
	if !ok {
		return false
	}
	return int32(updateResponse.ErrorCode) == 0
}

// ApplyThemeTask comment
type ApplyThemeTask struct {
	srcThemeIndex  uint32
	destThemeIndex uint32
}

// CreateApplyThemeTask comment
func CreateApplyThemeTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("apply theme params error!")
	}

	srcThemeIndex, error := strconv.Atoi(params[0])
	if error != nil {
		return nil, error
	}

	destThemeIndex, error := strconv.Atoi(params[1])
	if error != nil {
		return nil, error
	}

	task := &ApplyThemeTask{srcThemeIndex: (uint32)(srcThemeIndex), destThemeIndex: (uint32)(destThemeIndex)}
	return task, nil
}

// DoTask comment
func (t *ApplyThemeTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	return t.ApplyTheme(player)
}

// ApplyTheme comment
func (t *ApplyThemeTask) ApplyTheme(player *Player) bool {
	request := &protomsg.MsgCL2GSCityApplyThemeRequest{}
	request.ThemeIndex = t.destThemeIndex

	player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityApplyThemeReply)
	
	request.ThemeIndex = t.srcThemeIndex
	player.SendAndWait(request, msgtype.MsgType_kMsgGS2CLCityApplyThemeReply)

	return true
}
