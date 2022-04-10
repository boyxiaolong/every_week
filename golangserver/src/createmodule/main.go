package main

import (
	//. "game/common"
	//"github.com/issue9/term/colors"
	//"reflect"
	"runtime/debug"
	//"time"
	//"public/command"
	//"game/config"
	//"game/connect"
	"createmodule/module"
	"fmt"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	fmt.Println(os.Args)

	module_case := &module.ModuleDescribe{}
	module_case.SetModuleName(os.Args[1])
	generator := module.Generator{}
	generator.GeneratorAllFile(module_case)
}
