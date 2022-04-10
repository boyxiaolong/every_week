package loadlog

import (
	"behaviorlog/db"
	"behaviorlog/operat_file"
	"bufio"
	"fmt"
	"io"
	"os"
	"public/common"
	"strings"
)

const (
	FILE_TYPE_NONE   = 0
	FILE_TYPE_EVENT  = 1
	FILE_TYPE_OBJECT = 2
)

func GetAllLogFileName() {
	logDirectory := operat_file.GetCurrentDirectory()
	itemsDirectory := logDirectory + "/log"

	dirs, err := operat_file.GetDirs(itemsDirectory)
	if err != nil {
		return
	}

	for _, dir := range dirs {
		filenames, err := operat_file.GetAllFileName(dir)

		if err != nil {
			continue
		}

		for _, v := range filenames {
			var r = []rune(v)
			HandleFile(v, string(r[len(dir)+1:len(v)-6]))
		}
	}

}

func HandleFile(fullname string, filename string) {
	file_type := FILE_TYPE_NONE

	filename = strings.Replace(filename, ".", "_", -1)
	result := strings.Contains(fullname, "ObjectLog")

	if result {
		db.CreateObjectTable(db.OperateSqlInst.Db, filename)
		file_type = FILE_TYPE_OBJECT
	}

	result = strings.Contains(fullname, "EventLog")

	if result {
		db.CreateEventTable(db.OperateSqlInst.Db, filename)
		file_type = FILE_TYPE_EVENT
	}

	if file_type == FILE_TYPE_NONE {
		fmt.Println("file is error", filename)
		return
	}

	file, err := os.OpenFile(fullname, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			var info common.StringParse
			info.ParseString(line, "\t")

			if file_type == FILE_TYPE_OBJECT {
				line = db.InsertObjectTable(filename, &info)
			} else {
				line = db.InsertEventTable(filename, &info)
			}

			db.OperateSqlInst.AddSend(line)
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!", filename)
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}
}
