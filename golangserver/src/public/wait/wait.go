package wait

import (
	"public/common"
	"sync"
	"sync/atomic"
	"time"

	"public/message/msgtype"

	"github.com/golang/protobuf/proto"
)

var BREAK_STATUS_NORMAL uint32 = 0
var BREAK_STATUS_STOP uint32 = 1
var DefaultWaitSeconds uint16 = 20

type Wait struct {
	WaitMap      map[uint16]chan proto.Message
	mutex        sync.RWMutex
	Mode         int
	casename     string
	break_status uint32
}

func (m *Wait) SetCaseName(casename string) {
	m.casename = casename
}

func (m *Wait) add(key uint16) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.WaitMap[key]; ok {
		return
	}

	m.WaitMap[key] = make(chan proto.Message, 1)
}

func (m *Wait) get(key uint16) chan proto.Message {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.WaitMap[key]
}

func (m *Wait) del(key uint16) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.WaitMap, key)
}

// DoWait comment
func (m *Wait) DoWait(key uint16, seconds uint16) proto.Message {
	m.add(key)

	for {
		select {
		case message := <-m.get(key):
			m.del(key)
			return message
		case <-time.After(time.Duration(seconds) * time.Second): //超时1s
			if m.Mode == common.MODE_EXIT {
				common.GStdout.Error("wait msg type %v in %v", key, m.casename)
				common.QuitClient("wait msg time out", -1)
			} else if m.Mode == common.MODE_WAIT {
				if m.GetBreakStatus() == BREAK_STATUS_STOP {
					m.del(key)
					m.SetBreakStatus(0)
					common.GStdout.Error("wait (wait) msg type %v in %v error", key, m.casename)
					return nil
				}
				common.GStdout.Error("wait msg type %v in %v", key, m.casename)
			} else if m.Mode == common.MODE_CONTINUE {
				m.del(key)
				common.GStdout.Error("wait continue msg type %v in %v error", key, m.casename)
				return nil
			}
		}

	}

	return nil
}

func (m *Wait) Done(key uint16, message proto.Message) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	mc := m.WaitMap[key]

	if mc == nil {
		return
	}

	mc <- message
}

func (m *Wait) SetMode(mode int) {
	m.Mode = mode
}

func (m *Wait) SetBreakStatus(status uint32) {
	if m.Mode != common.MODE_WAIT {
		return
	}
	atomic.StoreUint32(&m.break_status, status)
}

func (m *Wait) GetBreakStatus() uint32 {
	if m.Mode != common.MODE_WAIT {
		return 0
	}
	return atomic.LoadUint32(&m.break_status)
}

func (m *Wait) Init(mode int) {
	m.WaitMap = make(map[uint16]chan proto.Message, 0)
	m.Mode = mode
}

func (m *Wait) WaitMultiple(args ...msgtype.MsgType) []proto.Message {
	var syncWait sync.WaitGroup
	messages := make([]proto.Message, len(args))
	syncWait.Add(len(args))
	for k, v := range args {
		go func(i int, msg_type msgtype.MsgType) {
			messages[i] = m.DoWait(uint16(msg_type), 7)
			syncWait.Done()
		}(k, v)
	}

	syncWait.Wait()
	return messages
}

func (m *Wait) Wait(msg_type msgtype.MsgType) proto.Message {
	return m.DoWait(uint16(msg_type), DefaultWaitSeconds)
}

func (m *Wait) WaitSeconds(msg_type msgtype.MsgType) proto.Message {
	return m.DoWait(uint16(msg_type), DefaultWaitSeconds)
}

func (m *Wait) Break() {
	m.SetBreakStatus(1)
}
