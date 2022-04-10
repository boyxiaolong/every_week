package main

import "behaviorlog/objectlog"
import (
	"behaviorlog/eventlog"
	"behaviorlog/items"
	"behaviorlog/loadlog"
	"public/command"
)

func main() {
	objectlog.LoadObjectLogFile()
	eventlog.LoadEventLogFile()
	items.GetItemsDirectory()
	loadlog.GetAllLogFileName()

	command.GCommand.Read()
}