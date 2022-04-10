package objectlog

import (
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
)

type ObjectLog struct {
	XMLName     xml.Name `xml:"metalib"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
	ObjectLogSubs   []ObjectLogSub `xml:"struct"`
}

type ObjectLogSub struct {
	XMLName     xml.Name `xml:"struct"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
	ObjectLogEntrys []ObjectLogEntry `xml:"entry"`
}

type ObjectLogEntry struct {
	XMLName     xml.Name `xml:"entry"`
	Name string   `xml:"name,attr"`
	Type string   `xml:"type,attr"`
	Desc string   `xml:"desc,attr"`
	ObjectLogItemInfo []ObjectLogItem `xml:"items"`
}

type ObjectLogItem struct {
	XMLName     xml.Name `xml:"items"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
}

func LoadObjectLogFile() {
	file, err := os.Open("ObjectLog.xml")
	if err != nil {
		fmt.Println("open xml file error")
		return
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file stream error")
		return
	}

	log := ObjectLog{}
	err = xml.Unmarshal(data, &log)
	if err != nil {
		fmt.Println("format xml data failed")
		return
	}

	fmt.Println("load objectlog success")
}