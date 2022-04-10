package items

import (
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
)

type XmlLogData struct {
	XMLName     xml.Name `xml:"metalib"`
	ItemsInfo   []XmlItems  `xml:"items"`
}

type XmlItems struct {
	XMLName     xml.Name `xml:"items"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
	ItemInfo   []XmlItem `xml:"item"`
}

type XmlItem struct {
	XMLName     xml.Name `xml:"item"`
	Name string   `xml:"name,attr"`
	Desc string   `xml:"desc,attr"`
}


func LoadFile(filename string) *XmlLogData {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open xml file error")
		return nil
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file stream error")
		return nil
	}

	log := &XmlLogData{}
	err = xml.Unmarshal(data, &log)
	if err != nil {
		fmt.Println("format xml data failed", err)
		return nil
	}

	return log
}