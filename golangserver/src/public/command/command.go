package command

import (
	"bufio"
	"fmt"
	"os"
	"public/common"
	"public/task"
	"runtime/debug"
	"strings"
)

type Command struct {
	commands map[string]*KeyCommand
	Task     *task.Task
}

func (m *Command) RegCommand(command string, callback Callback, note string) {
	if _, ok := m.commands[command]; ok {
		common.GStdout.Error("command:%v has reg", command)
		return
	}

	keyCommand := &KeyCommand{}
	keyCommand.Command = command
	keyCommand.Callback = callback
	keyCommand.Note = note
	m.commands[command] = keyCommand

	common.GStdout.Console("AddCommond Success:%v", command)
}

func (m *Command) UnRegCommand(command string) {
	delete(m.commands, command)
}

func (m *Command) Call(pStr *common.StringParse) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	if pStr == nil {
		common.GStdout.Error("command call is nil")
		return
	}
	if pStr.Len() <= 0 {
		common.GStdout.Error("command call is empty")
		return
	}
	keyCommand := m.commands[pStr.GetString(0)]

	if keyCommand == nil {
		common.GStdout.Error("command is not exist :%v", pStr.GetString(0))
		return
	}

	keyCommand.Callback(pStr)
}

func (m *Command) PrintCmd() {
	for k, v := range m.commands {
		common.GStdout.Console("%-20s%-10s%-30s", k, "-", v.Note)
	}
}

func (m *Command) ExecuteCommand(cmd string) {
	cmd = strings.Replace(cmd, "\n", "", -1)
	cmd = strings.Replace(cmd, "\r", "", -1)

	pStr := &common.StringParse{}
	pStr.ParseString(cmd, " ")
	m.Call(pStr)
}

func (m *Command) Read() {
	for {
		inputReader := bufio.NewReader(os.Stdin)
		cmd, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("There ware errors reading,exiting program.")
			return
		}

		m.Task.AddTask(func() { m.ExecuteCommand(cmd) })
	}
}

func NewCommand() *Command {
	return &Command{
		commands: make(map[string]*KeyCommand, 0),
	}
}
