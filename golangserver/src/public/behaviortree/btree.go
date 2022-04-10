package behaviortree

import (
	"fmt"
	"path/filepath"
	"public/common"
	"runtime/debug"
	"strings"
	"sync"
)

// ObjectInterface comment
type ObjectInterface interface {
	GetObjectId() uint64
}

// TaskTimeInfo comment
type TaskTimeInfo struct {
	times        uint32
	totalTimes   uint32
	preResetTime int64
}

// BTreeRunTask comment
type BTreeRunTask struct {
	TaskObj      ObjectInterface
	TaskData     *sync.Map
	TaskTimeData *sync.Map
	IsRun        bool
}

// StopRun comment
func (b *BTreeRunTask) StopRun() {
	b.IsRun = false
}

// GetTaskData comment
func (b *BTreeRunTask) GetTaskData(taskIndex uint32) interface{} {
	data, ok := b.TaskData.Load(taskIndex)
	if ok {
		return data
	}
	return nil
}

// SetTaskData comment
func (b *BTreeRunTask) SetTaskData(taskIndex uint32, data interface{}) {
	b.TaskData.Store(taskIndex,data)
}

// GetTaskTimeData comment
func (b *BTreeRunTask) GetTaskTimeData(taskIndex uint32) *TaskTimeInfo {
	data, ok := b.TaskTimeData.Load(taskIndex)
	if ok {
		return data.(*TaskTimeInfo)
	}
	timeInfo := &TaskTimeInfo{0, 0, 0}
	b.TaskTimeData.Store(taskIndex,timeInfo)
	return timeInfo
}

// CreateRunTask comment
func CreateRunTask(p ObjectInterface) *BTreeRunTask {
	newTask := &BTreeRunTask{TaskObj: p, TaskData: &sync.Map{}, IsRun: true, TaskTimeData: &sync.Map{}}
	return newTask
}

// BTree comment
type BTree struct {
	name string
	root *Node
}

// Do comment
func (tree *BTree) Do(runtask *BTreeRunTask) {
	if tree == nil || tree.root == nil || runtask.TaskObj == nil {
		return
	}

	for runtask.IsRun {
		tree.DoTask(runtask)
	}
	//common.GStdout.Console("stop run task ")
}

// GetName comment
func (tree *BTree) GetName() string {
	return tree.name
}

// DoTask comment
func (tree *BTree) DoTask(runtask *BTreeRunTask) {
	defer func() {
		if err := recover(); err != nil {
			//debug.PrintStack()
			common.GStdout.Error("%v", string(debug.Stack()))
		}
	}()

	tree.root.DoTask(runtask)
	return
}

// CreateTree comment
func CreateTree(name string) *BTree {
	btree := &BTree{}
	btree.name = name
	return btree
}

// Load comment
func Load(fileName string) (*BTree, error) {
	btreeConfig, error := LoadConfig(fileName)
	if btreeConfig == nil {
		common.GStdout.Error("load %v error %v\n", fileName, error)
		return nil, fmt.Errorf("load config %v error", fileName)
	}

	baseName := filepath.Base(fileName)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	common.GStdout.Console("%v", fileName)
	btree := &BTree{}
	btree.name = name
	node, error := createNode(btreeConfig.Root)
	if node == nil {
		return nil, error
	}
	btree.root = node
	return btree, nil
}
