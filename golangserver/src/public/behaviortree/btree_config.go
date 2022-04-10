package behaviortree

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// XMLTask comment
type XMLTask struct {
	TaskType   string `xml:"task_type,attr"`
	TaskParams string `xml:"task_params,attr"`
}

// XMLLoopNodeData comment
type XMLLoopNodeData struct {
	LoopMaxTimes int `xml:"loop_max_times,attr"`
}

// XMLRandomNodeData comment
type XMLRandomNodeData struct {
}

// XMLRandomWeightNodeData comment
type XMLRandomWeightNodeData struct {
	RandomWeigths string `xml:"random_weigth,attr"`
}

// XMLConditionNodeData comment
type XMLConditionNodeData struct {
	ConditionCheck  string `xml:"condition_check,attr"`
	ConditionResult bool   `xml:"condition_result,attr"`
	ConditionParams string `xml:"condition_params,attr"`
}

// XMLNode comment
type XMLNode struct {
	Name              string `xml:"name,attr"`
	Type              string `xml:"type,attr"`
	SleepTime         int    `xml:"sleep_time,attr"`
	ResetIntervelTime int64  `xml:"reset_intervel,attr"`
	DoTimes           uint32 `xml:"do_times,attr"`
	XMLTask
	XMLLoopNodeData
	XMLRandomNodeData
	XMLConditionNodeData
	XMLRandomWeightNodeData
	Childs []XMLNode `xml:"node"`
}

// XML comment
type XML struct {
	XMLName xml.Name `xml:"btree"`
	Root    XMLNode  `xml:"root"`
}

// LoadConfig comment
func LoadConfig(fileName string) (*XML, error) {
	if len(fileName) <= 0 {
		return nil, fmt.Errorf("emply file name")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("TIMBER! Can't load xml config file: %s %v", fileName, err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file stream error")
		return nil, fmt.Errorf("TIMBER! Can't parse xml config file: %s %v", fileName, err)
	}

	config := XML{}
	err = xml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("TIMBER! Can't parse xml config file: %s %v", fileName, err)
	}
	return &config, nil
}
