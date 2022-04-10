package data

import (
	"public/common"
	"public/connect"
	"public/message/msgtype"
	"public/task"
	"reflect"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

type DataInterface interface {
}

type DataCreatorCallBack func(data_center *DataCenter)

var GDataCenterCreator *DataCenterCreator

func RegisterDataCreator(cb DataCreatorCallBack) {
	if GDataCenterCreator == nil {
		GDataCenterCreator = &DataCenterCreator{}
		GDataCenterCreator.data_creators = make([]DataCreatorCallBack, 0)
	}
	GDataCenterCreator.data_creators = append(GDataCenterCreator.data_creators, cb)
}

func init() {
}

type DataCenterCreator struct {
	data_creators []DataCreatorCallBack
}

type DataCenter struct {
	connect.BaseDispatch
	datas          *sync.Map
	reflect_datas  *sync.Map
	lock           *sync.RWMutex
	msgHandlerTask *task.Task
}

func CreateDataCenter() *DataCenter {
	if GDataCenterCreator == nil {
		common.GStdout.Error("init data center creator error.")
		return nil
	}

	data_center := &DataCenter{}
	data_center.datas = &sync.Map{}
	data_center.reflect_datas = &sync.Map{} //make(map[string]*reflect.Value)
	data_center.lock = new(sync.RWMutex)
	data_center.msgHandlerTask = task.MakeTask(1000)
	data_center.Init()

	for _, v := range GDataCenterCreator.data_creators {
		v(data_center)
	}
	return data_center
}

func (t *DataCenter) DataRegister(data DataInterface) {
	s := reflect.ValueOf(data).Elem()
	typeOfT := s.Type()

	type_names := strings.Split(s.Type().String(), ".")
	type_name := type_names[len(type_names)-1]

	t.datas.Store(type_name, data)

	v := reflect.ValueOf(data)
	common.GStdout.Info("%v %v = %v\n", type_name, v.Type(), v.Interface())
	common.GStdout.Info("(%v) %v\n", type_name, reflect.ValueOf(data).Interface())
	t.reflect_datas.Store(type_name, &v)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		register_name := type_name + "." + typeOfT.Field(i).Name
		common.GStdout.Info("%d: (%v) %s.%s %s = %v\n", i, register_name, type_name, typeOfT.Field(i).Name, f.Type(), f.Interface())
		t.reflect_datas.Store(register_name, &f)
	}
}

func (m *DataCenter) Dispatch(msgtype msgtype.MsgType, msg proto.Message) {
	m.msgHandlerTask.AddTask(func() {
		if call, ok := m.CallBacks[msgtype]; ok {
			m.lock.Lock()
			defer m.lock.Unlock()
			call(msg)
		}
	})
}

func (t *DataCenter) GetData(data_name string) interface{} {
	data, ok := t.getData(data_name)
	if !ok {
		common.GStdout.Info("(%v) error\n", data_name)
		return nil
	}
	return data.Interface()
}

func (t *DataCenter) getData(data_name string) (*reflect.Value, bool) {
	if data, ok := t.reflect_datas.Load(data_name); ok {
		value, valueok := data.(*reflect.Value)
		return value, valueok
	}
	return nil, false
}

func (t *DataCenter) SetData(data_name string, value interface{}) {
	data, ok := t.getData(data_name)
	if !ok {
		return
	}

	if !data.CanSet() {
		return
	}

	if data.Type() != reflect.ValueOf(value).Type() {
		return
	}

	data.Set(reflect.ValueOf(value))
	return
}

func (t *DataCenter) Start() {
	t.msgHandlerTask.Start()
}

func (t *DataCenter) Stop() {
	t.msgHandlerTask.StopOnce()
}

func (t *DataCenter) Reset() {
	tlock := t.lock
	tlock.Lock()
	defer tlock.Unlock()

	t.datas = &sync.Map{}
	t.reflect_datas = &sync.Map{} //make(map[string]*reflect.Value)
	t.lock = new(sync.RWMutex)
	t.msgHandlerTask = task.MakeTask(1000)
	t.Init()
	for _, v := range GDataCenterCreator.data_creators {
		v(t)
	}
	t.Start()
}

