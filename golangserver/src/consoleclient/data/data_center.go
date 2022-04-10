package data

import (
	"consoleclient/application"
	"reflect"
)

// GDataCenter xxx
var GDataCenter *DataCenter

//GetDataCallBack 获取回调
type GetDataCallBack func() interface{}

// GetDataCallBackParam 带参数回调
type GetDataCallBackParam func(args ...interface{}) interface{}

//type SetDataCallBack func(param interface{})

// DataCenter 数据
type DataCenter struct {
}

func init() {
	GDataCenter = &DataCenter{}
}

//GetIntValue 获取数据
func (m *DataCenter) GetIntValue(call GetDataCallBack) int32 {
	value := GDataCenter.GetData(call)
	return value.(int32)
}

//GetUintValue 获取数据
func (m *DataCenter) GetUintValue(call GetDataCallBack) uint32 {
	value := GDataCenter.GetData(call)
	return value.(uint32)
}

//GetInt64Value 获取数据
func (m *DataCenter) GetInt64Value(call GetDataCallBack) int64 {
	value := GDataCenter.GetData(call)
	return value.(int64)
}

//GetUint64Value 获取数据
func (m *DataCenter) GetUint64Value(call GetDataCallBack) uint64 {
	value := GDataCenter.GetData(call)
	return value.(uint64)
}

//GetStringValue 获取数据
func (m *DataCenter) GetStringValue(call GetDataCallBack) string {
	value := GDataCenter.GetData(call)
	return value.(string)
}

//GetData 获取数据
func (m *DataCenter) GetData(call GetDataCallBack) interface{} {
	wait := make(chan interface{}, 1)

	application.GetApplication().DataTask.AddTask(func() {
		msg := call()
		wait <- msg
	})

	return <-wait
}

//GetDataParam 获取数据
func (m *DataCenter) GetDataParam(call GetDataCallBackParam, args ...interface{}) interface{} {
	wait := make(chan interface{}, 1)

	application.GetApplication().DataTask.AddTask(func() {
		msg := call(args...)
		wait <- msg
	})

	return <-wait
}

// SetData 设置
func (m *DataCenter) SetData(call interface{}, param interface{}) bool {
	callValue := reflect.ValueOf(call)
	paramValue := reflect.ValueOf(param)

	params := make([]reflect.Value, 1)
	params[0] = paramValue
	wait := make(chan bool, 1)

	application.GetApplication().DataTask.AddTask(func() {
		callValue.Call(params)
		wait <- true
	})

	return <-wait
}
