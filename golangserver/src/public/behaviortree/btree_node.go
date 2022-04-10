package behaviortree

import (
	"fmt"
	"math/rand"
	"public/common"
	"strconv"
	"strings"
	"time"
)

var nodeIndex uint32

func init() {
	nodeIndex = 1
}

// NodeImpInterface comment
type NodeImpInterface interface {
	DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool
}

// Node comment
type Node struct {
	nodeType      string
	nodeName      string
	nodeIndex     uint32
	doTimes       uint32
	resetIntervel int64
	sleepTime     int
	nodeImp       NodeImpInterface
}

// GetType comment
func (node *Node) GetType() string {
	return node.nodeType
}

// GetName comment
func (node *Node) GetName() string {
	return node.nodeName
}

// DoTask comment
func (node *Node) DoTask(runtask *BTreeRunTask) bool {
	common.GStdout.Debug("Start DoTask %v %v", node.GetName(), node.GetType())
	timeInfo := runtask.GetTaskTimeData(node.nodeIndex)
	curTime := time.Now().Unix()
	if node.doTimes > 0 {
		if node.resetIntervel > 0 && timeInfo.preResetTime+node.resetIntervel < curTime {
			timeInfo.times = 0
			timeInfo.preResetTime = 0
		}

		if timeInfo.times >= node.doTimes {
			return true
		}
	}

	res := node.nodeImp.DoTask(runtask, node.nodeIndex)
	if res {
		timeInfo.times++
		timeInfo.totalTimes++
		timeInfo.preResetTime = curTime
	} else {
		time.Sleep(time.Duration(60000) * time.Millisecond)
	}

	if node.sleepTime > 0 {
		time.Sleep(time.Duration(node.sleepTime) * time.Millisecond)
	}
	common.GStdout.Debug("DoTask %v %v Result %v", node.GetName(), node.GetType(), res)
	return res
}

// NodeCreatorCallback comment
type NodeCreatorCallback func(xmlnode XMLNode) (res NodeImpInterface, err error)

// RegisterNodeCreator comment
func RegisterNodeCreator(nodeType string, createFunc NodeCreatorCallback) error {
	nodeCreator := getNodesCreator()
	_, ok := nodeCreator.creator[nodeType]
	if ok {
		return fmt.Errorf("already register node creator %v", nodeType)
	}
	nodeCreator.creator[nodeType] = createFunc
	return nil
}

// NodeConditionCallback comment
type NodeConditionCallback func(testobj ObjectInterface, params string) bool

// RegisterConditionCallback comment
func RegisterConditionCallback(conditionType string, cb NodeConditionCallback) error {
	nodeCreator := getNodesCreator()
	_, ok := nodeCreator.conditionCb[conditionType]
	if ok {
		return fmt.Errorf("already register node creator %v", conditionType)
	}
	nodeCreator.conditionCb[conditionType] = cb
	return nil
}

func createNodeImp(xmlnode XMLNode) (NodeImpInterface, error) {
	nodeCreator := getNodesCreator()
	creator, ok := nodeCreator.creator[xmlnode.Type]
	if !ok {
		return nil, fmt.Errorf("not register node creator %v %v", xmlnode.Type, xmlnode.Name)
	}
	return creator(xmlnode)
}

// GetCondition comment
func GetCondition(conditionName string) (NodeConditionCallback, error) {
	nodeCreator := getNodesCreator()
	cb, ok := nodeCreator.conditionCb[conditionName]
	if !ok {
		return nil, fmt.Errorf("not register condition %v", conditionName)
	}
	return cb, nil
}

func createNode(xmlnode XMLNode) (*Node, error) {
	node := &Node{}
	node.nodeName = xmlnode.Name
	node.nodeType = xmlnode.Type
	node.nodeIndex = nodeIndex
	nodeIndex++
	node.sleepTime = xmlnode.SleepTime
	node.resetIntervel = xmlnode.ResetIntervelTime
	node.doTimes = xmlnode.DoTimes

	if node.resetIntervel > 0 && node.doTimes == 0 { // 默认为1
		node.doTimes = 1
	}

	nodeImp, error := createNodeImp(xmlnode)
	if nodeImp == nil {
		return nil, error
	}
	node.nodeImp = nodeImp
	return node, nil
}

func getNodesCreator() *NodesCreator {
	if GNodesCreator == nil {
		GNodesCreator = &NodesCreator{
			creator:     make(map[string]NodeCreatorCallback),
			conditionCb: make(map[string]NodeConditionCallback),
		}
	}
	return GNodesCreator
}

// GNodesCreator comment
var GNodesCreator *NodesCreator

// NodesCreator comment
type NodesCreator struct {
	creator     map[string]NodeCreatorCallback
	conditionCb map[string]NodeConditionCallback
}

// Init comment
func init() {
	RegisterNodeCreator("loop", CreateLoopNode)
	RegisterNodeCreator("random", CreateRandomNode)
	RegisterNodeCreator("randomweigth", CreateRandomWeigthNode)
	RegisterNodeCreator("selector", CreateSelectorNode)
	RegisterNodeCreator("sequence", CreateSequenceNode)
	RegisterNodeCreator("parallel", CreateParallelNode)
	RegisterNodeCreator("condition", CreateConditionNode)
	RegisterNodeCreator("task", CreateTaskNode)
}

// NodeLoop comment
type NodeLoop struct {
	childs   []*Node
	maxTimes int
}

// CreateLoopNode comment
func CreateLoopNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	selector := &NodeLoop{childs: make([]*Node, 0), maxTimes: xmlnode.LoopMaxTimes}
	for _, child := range xmlnode.Childs {
		childNode, error := createNode(child)
		if childNode == nil {
			return nil, error
		}
		selector.childs = append(selector.childs, childNode)
	}
	return selector, nil
}

// DoTask comment
func (node *NodeLoop) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	if node.maxTimes > 0 {
		for i := 0; i < node.maxTimes; i++ {
			for _, child := range node.childs {
				if !child.DoTask(runtask) {
					break
				}
			}
		}
	} else {
		for {
			for _, child := range node.childs {
				if !child.DoTask(runtask) {
					break
				}
			}
		}
	}
	return true
}

// NodeSelector comment
type NodeSelector struct {
	childs []*Node
}

// CreateSelectorNode comment
func CreateSelectorNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	selector := &NodeSelector{childs: make([]*Node, 0)}
	for _, child := range xmlnode.Childs {
		childNode, error := createNode(child)
		if childNode == nil {
			return nil, error
		}
		selector.childs = append(selector.childs, childNode)
	}
	return selector, nil
}

// DoTask comment
func (node *NodeSelector) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	for _, child := range node.childs {
		if child.DoTask(runtask) {
			return true
		}
	}
	return false
}

// NodeSequence comment
type NodeSequence struct {
	childs []*Node
}

// CreateSequenceNode comment
func CreateSequenceNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	node := &NodeSequence{childs: make([]*Node, 0)}
	for _, child := range xmlnode.Childs {
		childNode, error := createNode(child)
		if childNode == nil {
			return nil, error
		}
		node.childs = append(node.childs, childNode)
	}
	return node, nil
}

// DoTask comment
func (node *NodeSequence) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	for _, child := range node.childs {
		if !child.DoTask(runtask) {
			return false
		}
	}
	return true
}

// NodeParallel comment
type NodeParallel struct {
	childs []*Node
}

// CreateParallelNode comment
func CreateParallelNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	node := &NodeParallel{childs: make([]*Node, 0)}
	for _, child := range xmlnode.Childs {
		childNode, error := createNode(child)
		if childNode == nil {
			return nil, error
		}
		node.childs = append(node.childs, childNode)
	}
	return node, nil
}

// DoTask comment
func (node *NodeParallel) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	for _, child := range node.childs {
		child.DoTask(runtask)
	}
	return true
}

// NodeCondition comment
type NodeCondition struct {
	trueNode        *Node
	falseNode       *Node
	conditionCb     NodeConditionCallback
	conditionParams string
	checkRes        bool
}

// CreateConditionNode comment
func CreateConditionNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	node := &NodeCondition{}
	if len(xmlnode.Childs) == 0 {
		return nil, fmt.Errorf("create condition node error,child node size error")
	}

	trueNode, error := createNode(xmlnode.Childs[0])
	if trueNode == nil {
		return nil, error
	}
	node.trueNode = trueNode

	if len(xmlnode.Childs) > 1 {
		falseNode, error := createNode(xmlnode.Childs[1])
		if falseNode == nil {
			return nil, error
		}
		node.falseNode = falseNode
	}

	cb, error := GetCondition(xmlnode.ConditionCheck)
	if cb == nil {
		return nil, fmt.Errorf("create condition node error, not found condition [%v]", xmlnode.ConditionCheck)
	}

	node.conditionCb = cb
	node.conditionParams = xmlnode.ConditionParams
	node.checkRes = xmlnode.ConditionResult
	return node, nil
}

// DoTask comment
func (node *NodeCondition) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	if node.conditionCb(runtask.TaskObj, node.conditionParams) == node.checkRes {
		if node.trueNode == nil {
			return false
		}
		return node.trueNode.DoTask(runtask)
	}

	if node.falseNode == nil {
		return true
	}
	return node.falseNode.DoTask(runtask)
}

// NodeRandom comment
type NodeRandom struct {
	childs []*Node
}

// CreateRandomNode comment
func CreateRandomNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	if len(xmlnode.Childs) == 0 {
		return nil, fmt.Errorf("zero child")
	}

	random := &NodeRandom{
		childs: make([]*Node, 0),
	}

	for _, child := range xmlnode.Childs {
		childNode, error := createNode(child)
		if childNode == nil {
			return nil, error
		}
		random.childs = append(random.childs, childNode)
	}

	if len(random.childs) == 0 {
		return nil, fmt.Errorf("zero child")
	}
	return random, nil
}

// DoTask comment
func (node *NodeRandom) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	randomIndex := rand.Intn(len(node.childs))
	return node.childs[randomIndex].DoTask(runtask)
}

// NodeRandomWeigth comment
type NodeRandomWeigth struct {
	childs      []*Node
	nodeWeigths []int
	maxWeigths  int
}

// CreateRandomWeigthNode comment
func CreateRandomWeigthNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	randomWeigthsParams := strings.Split(xmlnode.RandomWeigths, ",")

	random := &NodeRandomWeigth{
		childs:      make([]*Node, 0),
		nodeWeigths: make([]int, 0),
	}

	if len(randomWeigthsParams) != len(xmlnode.Childs) {
		return nil, fmt.Errorf("create random node error,random data node error")
	}

	for k, child := range xmlnode.Childs {

		childNode, error := createNode(child)
		if childNode == nil {
			return nil, error
		}

		weigth, error := strconv.Atoi(randomWeigthsParams[k])
		if error != nil {
			return nil, fmt.Errorf("create random node error,random data node error")
		}
		random.childs = append(random.childs, childNode)
		random.nodeWeigths = append(random.nodeWeigths, weigth)
		random.maxWeigths += weigth
	}
	return random, nil
}

// DoTask comment
func (node *NodeRandomWeigth) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	randomWeigth := rand.Intn(node.maxWeigths)
	for k, v := range node.nodeWeigths {
		if v > randomWeigth {
			return node.childs[k].DoTask(runtask)
		}
		randomWeigth -= v
	}
	return false
}

// NodeTask comment
type NodeTask struct {
	task *BTTask
}

// CreateTaskNode comment
func CreateTaskNode(xmlnode XMLNode) (res NodeImpInterface, err error) {
	task, error := CreateTask(xmlnode.TaskType, xmlnode.ResetIntervelTime, xmlnode.DoTimes, xmlnode.TaskParams)
	if task == nil {
		return nil, error
	}
	node := &NodeTask{}
	node.task = task
	return node, nil
}

// DoTask comment
func (node *NodeTask) DoTask(runtask *BTreeRunTask, nodeIndex uint32) bool {
	return node.task.DoTask(runtask, nodeIndex)
}
