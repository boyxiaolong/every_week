package command

import (
	//"fmt"
	"os"
	"os/exec"
	//"path/filepath"
	"fmt"
	"io/ioutil"
	"net/http"
	"public/common"
	"public/task"
)

var (
	GCommand *Command
)

func init() {
	GCommand = NewCommand()
	GCommand.RegCommand("help", PrintCommand, "echo all command")
	GCommand.RegCommand("quit", Quit, "quit")
	GCommand.RegCommand("cls", Cls, "cls")
	GCommand.RegCommand("temp", temp, "temp")
}

func SetCommandTask(command *Command, Task *task.Task) {
	command.Task = Task
}

func PrintCommand(str *common.StringParse) (err error) {
	GCommand.PrintCmd()
	return
}

func Quit(str *common.StringParse) (err error) {
	common.QuitClient("quit command", 0)
	return
}

func Cls(str *common.StringParse) (err error) {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return
}

func temp(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		return
	}

	url := str.GetString(1)

	resp, e := http.Get(url)
	if e != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", e)
		return e
	}

	b, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		fmt.Fprintf(os.Stderr, "fetch reading:%s %v\n", url, e)
		return e
	}

	fmt.Printf("%s", b)
	return
}
