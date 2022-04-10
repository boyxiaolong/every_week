package module

import (
	"consoleclient/application"
	"public/common"
	"public/config"
	"public/message/error_code"
	"time"
)

//Call Call
type Call func() (bool, int32)

//CallP CallP
type CallP func(str *common.StringParse) (bool, int32)

//Test Test
type Test interface {
	TestModule(application *application.Application)
	TestCommand(cmd string, application *application.Application)
	TestCommandParam(cmd string, application *application.Application, str *common.StringParse)
	ListCommand()
}

//TestBase TestBase
type TestBase struct {
	Cmds      map[string]Call
	CmdList   []string
	Name      string
	CmdParams map[string]CallP
}

//InitCmds InitCmds
func (m *TestBase) InitCmds(name string) {
	m.Cmds = make(map[string]Call)
	m.CmdList = make([]string, 0)
	m.Name = name
	m.CmdParams = make(map[string]CallP)
	AddModule(name, m)
}

//TestModule TestModule
func (m *TestBase) TestModule(application *application.Application) {
	for _, v := range m.CmdList {
		m.TestCommand(v, application)
		time.Sleep(700 * time.Millisecond)
	}
}

//RegCommand RegCommand
func (m *TestBase) RegCommand(cmd string, call Call) {
	m.Cmds[cmd] = call
	m.CmdList = append(m.CmdList, cmd)
	//GStdout.Info("%v", cmd)
}

//RegCommandParam RegCommand
func (m *TestBase) RegCommandParam(cmd string, call CallP) {
	m.CmdParams[cmd] = call
	//GStdout.Info("%v", cmd)
}

//TestCommand TestCommand
func (m *TestBase) TestCommand(cmd string, application *application.Application) {
	if call, ok := m.Cmds[cmd]; ok {
		GGameInfo.WaitMessage.SetCaseName(cmd)
		common.GStdout.Success("start test %v %v", m.Name, cmd)
		result, code := call()

		if result {
			common.GStdout.Success("test %v %v success", m.Name, cmd)
		} else {
			common.GStdout.Error("test %v %v false,code:%v", m.Name, cmd, error_code.ErrorCode_name[code])
			if config.Mode == common.MODE_EXIT {
				common.QuitClient("test failed", -1)
			} else if config.Mode == common.MODE_WAIT {
				wait := make(chan int)
				<-wait
			} else if config.Mode == common.MODE_CONTINUE {
				return
			}
		}
	} else {
		common.GStdout.Error("cmd not exist:%v", cmd)
	}
}

//TestCommandParam TestCommand
func (m *TestBase) TestCommandParam(cmd string, application *application.Application, str *common.StringParse) {
	if call, ok := m.CmdParams[cmd]; ok {
		GGameInfo.WaitMessage.SetCaseName(cmd)
		result, code := call(str)

		if result {
			common.GStdout.Success("test %v %v success", m.Name, cmd)
		} else {
			common.GStdout.Error("test %v %v false,code:%v", m.Name, cmd, error_code.ErrorCode_name[code])
			if config.Mode == common.MODE_EXIT {
				common.QuitClient("test failed", -1)
			} else if config.Mode == common.MODE_WAIT {
				wait := make(chan int)
				<-wait
			} else if config.Mode == common.MODE_CONTINUE {
				return
			}
		}
	} else {
		common.GStdout.Error("cmd not exist:%v", cmd)
	}
}

//ListCommand ListCommand
func (m *TestBase) ListCommand() {
	for k := range m.Cmds {
		common.GStdout.Info("%v", k)
	}
}
