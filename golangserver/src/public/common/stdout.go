package common

import (
	"fmt"
	"os"
	"public/config"
	"public/timber"
	"runtime/debug"
	"time"

	"github.com/issue9/term/ansi"
	"github.com/issue9/term/colors"
	//"runtime"
)

var OUT_PUT_MODULE = "windows"

var ConsoleLogSetting int

func SetOutPutModule(output_module string) {
	OUT_PUT_MODULE = output_module
}

func GetOutPutModule() string {
	return OUT_PUT_MODULE
}

var ConsoleLogStrSetting = []string{
	"DEBUG",
	"INFO",
	"SUCCESS",
	"ERROR",
}

func SetConsoleLevel(levels []string) {
	ConsoleLogSetting = 0
	for _, v := range levels {
		if "DEBUG" == v {
			ConsoleLogSetting |= CONSOLE_LOG_DEBUG
		} else if "INFO" == v {
			ConsoleLogSetting |= CONSOLE_LOG_INFO
		} else if "SUCCESS" == v {
			ConsoleLogSetting |= CONSOLE_LOG_SUCCESS
		} else if "ERROR" == v {
			ConsoleLogSetting |= CONSOLE_LOG_ERROR
		}
	}
	fmt.Printf("console log level %v\n", ConsoleLogSetting)
}

func CheckConsoleSetting(console_type int) bool {
	return ConsoleLogSetting&console_type != 0
}

func init() {
	SetConsoleLevel(config.ConsoleLogLevel)
}

type Stdout struct {
	StringChan chan func()
	GLog       *timber.Timber
}

func (m *Stdout) Info(format string, v ...interface{}) {
	format = fmt.Sprintf("%v\n", format)

	m.Add(func() {
		if CheckConsoleSetting(CONSOLE_LOG_INFO) {
			fmt.Printf(format, v...)
		}

		GLog.Info(format, v...)
	})
}

func (m *Stdout) Console(format string, v ...interface{}) {
	format = fmt.Sprintf("%v\n", format)

	m.Add(func() {
		fmt.Printf(format, v...)
	})
}

func (m *Stdout) Success(format string, v ...interface{}) {
	format = fmt.Sprintf("%v\n", format)

	m.Add(func() {
		if CheckConsoleSetting(CONSOLE_LOG_SUCCESS) {
			if OUT_PUT_MODULE == "windows" {
				colors.Printf(colors.Green, colors.Black, format, v...)
			} else if OUT_PUT_MODULE == "linux" {
				format = fmt.Sprintf("%v%v%v%v", ansi.FGreen, ansi.BBlack, format, ansi.Reset)
				fmt.Printf(format, v...)
			} else {
				fmt.Printf(format, v...)
			}
		}

		GLog.Info(format, v...)
	})
}

func (m *Stdout) Debug(format string, v ...interface{}) {
	format = fmt.Sprintf("%v\n", format)

	m.Add(func() {
		if CheckConsoleSetting(CONSOLE_LOG_DEBUG) {
			fmt.Printf(format, v...)
		}

		GLog.Debug(format, v...)
	})
}

func (m *Stdout) Error(format string, v ...interface{}) {
	format = fmt.Sprintf("%v\n", format)

	m.Add(func() {
		//colors.Printf(colors.Red, colors.Black, format, v...)
		if CheckConsoleSetting(CONSOLE_LOG_ERROR) {
			fmt.Printf(format, v...)
			if OUT_PUT_MODULE == "windows" {
				colors.Printf(colors.Red, colors.Black, format, v...)
			} else if OUT_PUT_MODULE == "linux" {
				format = fmt.Sprintf("%v%v%v%v", ansi.FRed, ansi.BBlack, format, ansi.Reset)
				fmt.Printf(format, v...)
			} else {
				fmt.Printf(format, v...)
			}
		}

		//GLog.Error(format, v...)
	})
}

func (m *Stdout) Add(call func()) {
	m.StringChan <- call
}

func (m *Stdout) Do() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	for {
		select {
		case out := <-m.StringChan:
			out()
		}
	}
}

var (
	GStdout *Stdout
)

func init() {
	GStdout = &Stdout{
		StringChan: make(chan func(), 100),
	}

	go GStdout.Do()
}

//QuitClient 退出处理
func QuitClient(reason string, exitCode int) {
	//GLog.Info("quit reason %v", reason)
	time.Sleep(5 * time.Second)
	os.Exit(exitCode)
}
