package config

import (
	//"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/unknwon/com"
	"github.com/unknwon/goconfig"
)

var (
	WorkDir         string
	Cfg             *goconfig.ConfigFile
	LogXmlConf      string
	Mode            int
	ConsoleLogLevel []string
	MsgLogLevel     int
)

func GetWorkDir() (string, error) {
	p := os.Args[0]
	p, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return path.Dir(strings.Replace(p, "\\", "/", -1)), err
}

func init() {
	Load()
}

func Load() {
	wDir, err := GetWorkDir()
	if err != nil {
		panic("Fail to get work directory")
	}
	WorkDir = wDir
	filename := filepath.Join(WorkDir, "conf/config.ini")
	if !com.IsFile(filename) {
		//为initdo目录做一个入口补丁
		WorkDir = path.Dir(WorkDir)
		filename = filepath.Join(WorkDir, "conf/config.ini")
		if !com.IsFile(filename) {
			panic("Fail to read:" + filename)
		}
	}

	Cfg, err = goconfig.LoadConfigFile(filename)
	if err != nil {
		panic("Fail to read:" + filename)
	}

	Mode = Cfg.MustInt("test", "mode", 1)
	LogXmlConf = Cfg.MustValue("log", "log_xmlconf", "log.xml")
	ConsoleLogLevel = Cfg.MustValueArray("log", "console_log_level", ",")
	MsgLogLevel = Cfg.MustInt("log", "msg_log_level", 1)
}
