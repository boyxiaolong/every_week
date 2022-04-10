package behaviortree

import (
	"fmt"
	"public/common"
	"strings"
	"sync/atomic"
)

// TaskCreatorCallback comment
type TaskCreatorCallback func(params []string) (res BTTaskInterface, err error)

// RegisterTaskCreator comment
func RegisterTaskCreator(taskType string, createFunc TaskCreatorCallback) error {
	taskCreator := getTasksCreator()
	_, ok := taskCreator.creator[taskType]
	if ok {
		return fmt.Errorf("already register task creator %v", taskType)
	}
	taskCreator.creator[taskType] = createFunc

	GTasksRecord.Record[taskType] = &TaskRecord{}
	return nil
}

// CreateTask comment
func CreateTask(taskType string, resetIntervel int64, doTimes uint32, params string) (*BTTask, error) {
	task := &BTTask{}
	task.taskType = taskType
	taskImp, error := createTaskImp(taskType, params)
	if taskImp == nil {
		return nil, error
	}
	task.taskImp = taskImp
	return task, nil
}

func createTaskImp(taskType string, params string) (BTTaskInterface, error) {
	taskCreator := getTasksCreator()
	creator, ok := taskCreator.creator[taskType]
	if !ok {
		return nil, fmt.Errorf("not register task creator %v", taskType)
	}

	paramsList := make([]string, 0, 0)
	if len(params) > 0 {
		paramsList = strings.Split(params, ",")
	}
	return creator(paramsList)
}

func getTasksCreator() *TasksCreator {
	if GTasksCreator == nil {
		GTasksCreator = &TasksCreator{
			creator: make(map[string]TaskCreatorCallback),
		}
	}
	return GTasksCreator
}

// BTTask comment
type BTTask struct {
	taskType string
	taskImp  BTTaskInterface
}

// DoTask comment
func (task *BTTask) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {

	common.GStdout.Debug("=========obj %v %v start ", runtask.TaskObj.GetObjectId(), task.GetType())
	res := task.taskImp.DoTask(runtask, nodeIndex)
	common.GStdout.Debug("=========obj %v %v end %v ", runtask.TaskObj.GetObjectId(), task.GetType(), res)
	GTasksRecord.RecordTask(task.GetType(), res)
	return res
}

// GetType comment
func (task *BTTask) GetType() string {
	return task.taskType
}

// GTasksCreator comment
var GTasksCreator *TasksCreator

// TasksCreator comment
type TasksCreator struct {
	creator map[string]TaskCreatorCallback
}

// TaskRecord comment
type TaskRecord struct {
	SuccessCount int32
	FailureCount int32
}

// TasksRecord comment
type TasksRecord struct {
	Record map[string]*TaskRecord
}

// GTasksRecord comment
var GTasksRecord *TasksRecord

// RecordTask comment
func (t *TasksRecord) RecordTask(taskName string, isSuc bool) {
	if v, k := t.Record[taskName]; k {
		if isSuc {
			atomic.AddInt32(&v.SuccessCount, 1)
		} else {
			atomic.AddInt32(&v.FailureCount, 1)
		}
	}
}

// PrintTaskRecord comment
func PrintTaskRecord() {
	for k, v := range GTasksRecord.Record {
		successCount := atomic.LoadInt32(&v.SuccessCount)
		failureCount := atomic.LoadInt32(&v.FailureCount)
		if successCount > 0 || successCount > 0 {
			common.GStdout.Console("task %v SuccessCount[%v]  , FailureCount[%v]", k, successCount, failureCount)
		}
	}
}

// BTTaskInterface comment
type BTTaskInterface interface {
	DoTask(runtask *BTreeRunTask, taskindex uint32) bool
}

// Init comment
func init() {
	GTasksRecord = &TasksRecord{
		Record: make(map[string]*TaskRecord)}

	RegisterTaskCreator("Default", CreateDefaultTask)
	RegisterTaskCreator("Stop", CreateStopTask)

}

// DefaultTask comment
type DefaultTask struct {
	params string
}

// DoTask comment
func (task *DefaultTask) DoTask(runtask *BTreeRunTask, taskindex uint32) bool {
	return true
}

// CreateDefaultTask comment
func CreateDefaultTask(params []string) (res BTTaskInterface, err error) {
	task := &DefaultTask{}
	return task, nil
}

// StopTask comment
type StopTask struct {
}

// DoTask comment
func (t *StopTask) DoTask(runtask *BTreeRunTask, taskindex uint32) bool {
	runtask.StopRun()
	return true
}

// CreateStopTask comment
func CreateStopTask(params []string) (res BTTaskInterface, err error) {
	task := &StopTask{}
	return task, nil
}
