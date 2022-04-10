package eventlog

import (
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
)

type EventLog struct {
	XMLName     xml.Name `xml:"metalib"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
	EventLogSubs   []EventLogSub `xml:"struct"`
}

type EventLogSub struct {
	XMLName     xml.Name `xml:"struct"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
	EventLogSubEntrys []EventLogSubEntry `xml:"entry"`
}

type EventLogSubEntry struct {
	XMLName     xml.Name `xml:"entry"`
	Name string   `xml:"name,attr"`
	Type string   `xml:"type,attr"`
	Desc string   `xml:"desc,attr"`
	EventLogItemInfo []EventLogSubItem `xml:"items"`
}

type EventLogSubItem struct {
	XMLName     xml.Name `xml:"items"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
}

func LoadEventLogFile() {
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

	log := EventLog{}
	err = xml.Unmarshal(data, &log)
	if err != nil {
		fmt.Println("format xml data failed")
		return
	}

	fmt.Println("load objectlog success")
}