package main

import (
	//. "game/common"
	//"github.com/issue9/term/colors"
	//"reflect"

	"runtime/debug"
	//"time"
	"public/command"
	//"game/config"
	//"game/connect"
	"consoleclient/application"
	"consoleclient/module"
	"os"
	"public/common"
	"public/config"
	"public/winapi"
	_"consoleclient/data"
)

func RunApp() {

}

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	winapi.SetConsoleMode()

	if len(os.Args) > 1 {
		common.SetOutPutModule(os.Args[1])
	}

	module.GTestMgr.InitApplication(application.GetApplication())

	if config.Mode == common.MODE_EXIT {
		module.GTestMgr.StartTest()
	}

	common.GStdout.Success("INIT SUCCESS")
	command.SetCommandTask(command.GCommand,application.GetApplication().MainTask)
	command.GCommand.Read()
}
