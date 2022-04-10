package common

import (
	"fmt"
	"os"
	"path"
	"public/config"

	"public/timber"
)

type Logger struct {
	*timber.Timber
}

var (
	GLog *Logger
)

//提供一个全局的log对象
func GetLogger() *Logger {
	return GLog
}

func init() {
	mkdir := "gamelog"
	_, err := os.Stat(mkdir)
	if err == nil || os.IsExist(err) {

	} else {
		os.Mkdir(mkdir, 0777)
	}

	man := Logger{timber.Global}
	GLog = &man
	err = GLog.LoadXMLConfig(path.Join(config.WorkDir, "conf", config.LogXmlConf))

	if err != nil {
		fmt.Println(err.Error())
	}
}
