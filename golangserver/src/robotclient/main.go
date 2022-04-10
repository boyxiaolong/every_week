package main

import (
	"flag"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"public/command"
	"public/common"
	"public/task"
	_ "robotclient/client"
	_ "robotclient/client/player"
	_ "robotclient/client/test"
	"runtime/debug"
	"time"
)

func main() {
	//runtime.GOMAXPROCS(1)
	go http.ListenAndServe("0.0.0.0:8080", nil)
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	rand.Seed(time.Now().Unix())

	cmd := flag.String("cmd", "", "default command")
	flag.Parse()

	main_task := task.MakeTask(1000)
	command.SetCommandTask(command.GCommand, main_task)
	main_task.Start()
	if *cmd != "" {
		common.GStdout.Console("==============================start cmd [%v] ================================", *cmd)
		command.GCommand.Task.AddTask(func() { command.GCommand.ExecuteCommand(*cmd) })
	}
	command.GCommand.Read()
}
