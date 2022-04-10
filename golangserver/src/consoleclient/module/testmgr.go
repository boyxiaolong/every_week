package module

import (
	"consoleclient/application"
	"consoleclient/loadconfig"
	"public/command"
	"public/common"
	"public/config"
	"runtime/debug"
	"sort"
	"time"
)

//GTestMgr GTestMsg
var GTestMgr *TestMgr

func init() {
	command.GCommand.RegCommand("test", TestCase, "Execute special test case")
	command.GCommand.RegCommand("testsuit", TestSuit, "Execute special test suit")
	command.GCommand.RegCommand("testall", TestAll, "Execute all test cases")
	command.GCommand.RegCommand("listcase", ListCase, "Show all test suits")
	command.GCommand.RegCommand("listsuit", ListSuit, "Show all test case of special suit")

	InitTestMgr()
}

//InitTestMgr comment
func InitTestMgr() {
	if GTestMgr == nil {
		GTestMgr = &TestMgr{
			Tests:     make(map[string]Test),
			FuncChan:  make(chan func(), 50),
			TestsName: make([]string, 0),
		}
	}
}

//AddModule AddModule
func AddModule(name string, test Test) {
	InitTestMgr()
	GTestMgr.AddModule(name, test)
}

//TestSuit TestSuit
func TestSuit(str *common.StringParse) (err error) {
	if str.Len() <= 1 {
		common.GStdout.Error("test param error")
		return
	}

	name := str.GetString(1)
	GTestMgr.AddTestSuit(name)

	return
}

//TestCase TestCase
func TestCase(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("testcase param error")
		return
	}

	if str.Len() == 2 {
		name := str.GetString(1)
		GTestMgr.AddTestSuit(name)

	} else if str.Len() == 3 {
		module := str.GetString(1)
		cmd := str.GetString(2)
		GTestMgr.AddTestModule(module, cmd)
	} else {
		module := str.GetString(1)
		cmd := str.GetString(2)
		GTestMgr.AddTestModuleP(module, cmd, str)
	}

	return
}

//TestAll TestAll
func TestAll(str *common.StringParse) (err error) {
	AddTestAll()
	return
}

//AddTestAll TestAll
func AddTestAll() (err error) {
	if GGameInfo.Online == false {
		common.GStdout.Error("user not online")
		return
	}

	GTestMgr.TestAllModule()
	return
}

//ListSuit ListSuit
func ListSuit(str *common.StringParse) (err error) {
	for k := range GTestMgr.Tests {
		common.GStdout.Info(k)
	}
	return
}

//ListCase ListCase
func ListCase(str *common.StringParse) (err error) {
	if str.Len() < 2 {
		common.GStdout.Error("ListCase param error")
		return
	}

	if suit, ok := GTestMgr.Tests[str.GetString(1)]; ok {
		suit.ListCommand()
	} else {
		common.GStdout.Error("not suit name:%v", str.GetString(1))
	}

	return
}

//TestMgr TestMgr
type TestMgr struct {
	Tests       map[string]Test
	FuncChan    chan func()
	TestsName   []string
	Application *application.Application
}

//InitApplication comment
func (m *TestMgr) InitApplication(application *application.Application) {
	m.Application = application
}

//AddTestSuit comment
func (m *TestMgr) AddTestSuit(name string) {
	test := GTestMgr.GetModule(name)

	if test == nil {
		common.GStdout.Error("%v module not exist", name)
		return
	}

	if GGameInfo.Online == false {
		common.GStdout.Error("user not online")
		return
	}

	test.TestModule(GTestMgr.Application)

	return
}

//AddTestModule comment
func (m *TestMgr) AddTestModule(module string, cmd string) {
	if GGameInfo.Online == false {
		common.GStdout.Error("user not online")
		return
	}

	test := m.GetModule(module)

	if test == nil {
		common.GStdout.Error("Test module not exist")
		return
	}

	test.TestCommand(cmd, m.Application)
}

//AddTestModuleP comment
func (m *TestMgr) AddTestModuleP(module string, cmd string, str *common.StringParse) {
	if GGameInfo.Online == false {
		common.GStdout.Error("user not online")
		return
	}

	test := m.GetModule(module)

	if test == nil {
		common.GStdout.Error("Test module not exist")
		return
	}

	test.TestCommandParam(cmd, m.Application, str)
}

//GetModule GetModule
func (m *TestMgr) GetModule(name string) Test {
	if test, ok := m.Tests[name]; ok {
		return test
	}

	return nil
}

//AddModule AddModule
func (m *TestMgr) AddModule(name string, test Test) {
	m.Tests[name] = test
	common.GStdout.Info("add module : %v", name)
	m.TestsName = append(m.TestsName, name)
}

//TestModule TestModule
func (m *TestMgr) TestModule(name string) {
	if test, ok := m.Tests[name]; ok {
		test.TestModule(m.Application)
		time.Sleep(700 * time.Millisecond)
	} else {
		common.GStdout.Error("%v module not exist", name)
	}
}

//TestAllModule TestAllModule
func (m *TestMgr) TestAllModule() {
	sort.Strings(m.TestsName)

	m.Application.MainTask.AddTask(func() {
		common.GStdout.Info("begin test all module")
	})

	if loadconfig.TestAllType == 1 {
		for _, v := range m.TestsName {
			module := m.GetModule(v)
			module.TestModule(m.Application)
		}
	} else {
		for _, v := range m.Tests {
			v.TestModule(m.Application)
		}
	}

	m.Application.MainTask.AddTask(func() {
		common.GStdout.Info("end test all module")
		if config.Mode == common.MODE_EXIT {
			common.QuitClient("wait time out", 0)
		}
	})

}

//AddCall AddCall
func (m *TestMgr) AddCall(call func()) {
	select {
	case m.FuncChan <- call:
	}
}

//StartTest StartTest
func (m *TestMgr) StartTest() {
	LoginToServer(loadconfig.Account)
}

//CallFunc CallFunc
func (m *TestMgr) CallFunc(callFunc func()) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	callFunc()
}
