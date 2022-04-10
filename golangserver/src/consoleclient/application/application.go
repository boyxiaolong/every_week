package application

import (
	"public/task"
)

var Gapplicaiton *Application

func init() {
	Gapplicaiton = InitApplication()
	Gapplicaiton.MainTask.Start()
	Gapplicaiton.DataTask.Start()
}

func GetApplication() *Application {
	return Gapplicaiton
}

func InitApplication() *Application {
	return &Application{
		MainTask: task.MakeTask(1000),
		DataTask: task.MakeTask(1000),
	}
}

type Application struct {
	MainTask *task.Task
	DataTask *task.Task
}
