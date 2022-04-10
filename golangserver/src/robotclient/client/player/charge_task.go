package player

import (
	"fmt"
	BTree "public/behaviortree"
	"strconv"
	"encoding/binary"
	"net"
	"encoding/json"
	"time"
	"bytes"
	"robotclient/loadconfig"
	"public/common"
	"math/rand"
	//"public/common"
)

var ValidPCID = [...]uint64{30866,31423,1003,1004,1005,1011,1012,1013,1014,1015,1021,1022,1023,1024,1025,1031,1032,1033,1034,1035,1041,1042,1043,1044,1045,1051,1052,1053,1054,1055,1101,1102,1103,1104,1105,1106,1111,1112,1113,1114,1115,1121,1122,1123,1124,1125,2001,2002,2003,2004,2005,4001,4002,4003,4004,4005,4006,4011,4012,4013,4014,4015,4016,4051,4052,4053,4054,4055,4056,4057,4058,4059,4060,4061,4062,4063,4064,4065,4066,3011,3012,5001,5002,5003,5004,5005,5006,5007,5008,5009,5011,5012,5013,5014,5015,5016,5021,5022,5023,5031,5032,5033,5041,5042,5043,5044,5045,5051,5052,5053,5056,5057,5058,5061,5062,5063,5066,5067,5068,4102,4103,4104,4112,4113,4114,4122,4123,4124,4132,4133,4134,4142,4143,4144,4152,4153,4154,6001,6002,6003,6004,6005,6006,6007,6011,6012,6013,6014,6015,6016,6017,7001,7002,7003,7004,8001,8002,8003,8004,9001,9002,9003}

// Init comment
func init() {
	BTree.RegisterTaskCreator("Charge", CreateChargeTask)
	BTree.RegisterTaskCreator("RandomCharge", CreateRandomChargeTask)
}

type ChargeResultJson struct {
	Desc string `json:"desc"` 
	Err  uint32 `json:"error"` 
}

func ChargeRequest(iggId uint64,pcid uint64,giftInfo string) error {
	address := fmt.Sprintf("%v:%v",loadconfig.GetChargeIp(), loadconfig.GetChargePort())

	buf := new(bytes.Buffer)
	var msg_type uint16 = 6880
	var sn = fmt.Sprintf("test_%v_%v_%v",iggId,pcid,time.Now().Unix())
	var msg string = fmt.Sprintf("{\"sn\":\"%v\",\"iggid\":%v,\"pc_id\":%v,\"gift_info\":\"%v\"}", sn,iggId,pcid,giftInfo)
	var msg_len uint16 = 4 + (uint16)(len(msg))

	err := binary.Write(buf, binary.LittleEndian, msg_len)
	if err != nil {
		return err
	}

	err = binary.Write(buf, binary.LittleEndian, msg_type)
	if err != nil {
		return err
	}

	msgb := []byte(msg)
	err = binary.Write(buf, binary.LittleEndian, msgb)
	if err != nil {
		return err
	}

	conn,err:= net.Dial("tcp", address)
	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.Write(buf.Bytes()); err != nil {
		return err
	}

	readbuf := make([]byte, 1024)
	readlen,err := conn.Read(readbuf)
	if err != nil {
		return err
	}

	var res ChargeResultJson
	err = json.Unmarshal(readbuf[4:readlen], &res)
	if err != nil {
		return err
	}

	if res.Err == 0 {
		return nil
	}
	return fmt.Errorf("charge request error% v!",res.Err)
}

// CreateChargeTask comment
func CreateChargeTask(params []string) (res BTree.BTTaskInterface, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("charge Task params error!")
	}

	pcid, error := strconv.ParseUint(params[0], 10, 64)
	if error != nil {
		return nil, error
	}

	giftInfo := params[1]
	if error != nil {
		return nil, error
	}

	task := &ChargeTask{pcid: pcid,giftInfo:giftInfo}
	return task, nil
}

// ChargeTask comment
type ChargeTask struct {
	pcid 		 uint64
	giftInfo string
}

// DoTask comment
func (t *ChargeTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	if err := ChargeRequest(player.PlayerID,t.pcid,t.giftInfo); err != nil {
		common.GStdout.Info("ChargeRequest error %v", err)
		return false
	}
	return true
}

// RandomChargeTask comment
func CreateRandomChargeTask(params []string) (res BTree.BTTaskInterface, err error) {
	task := &RandomChargeTask{}
	return task, nil
}

// RandomChargeTask comment
type RandomChargeTask struct {
}

// DoTask comment
func (t *RandomChargeTask) DoTask(runtask *BTree.BTreeRunTask, taskindex uint32) bool {
	player := runtask.TaskObj.(*Player)
	pcid := ValidPCID[rand.Intn(len(ValidPCID))]
	if err := ChargeRequest(player.PlayerID,pcid,""); err != nil {
		common.GStdout.Info("ChargeRequest error %v", err)
		return false
	}
	return true
}
