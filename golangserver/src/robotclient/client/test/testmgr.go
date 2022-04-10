package test

import (
	"os"
	"path"
	"path/filepath"
	BTree "public/behaviortree"
	"public/command"
	"public/common"
	"robotclient/client/player"
	"runtime"
	"time"
)

var GTestBTree *TestBTree

type TestBTree struct {
	testbtree map[string]*BTree.BTree
}

func (t *TestBTree) LoadAllTree(find_path string) {
	var files []string

	err := filepath.Walk(find_path, func(file_path string, info os.FileInfo, err error) error {
		if path.Ext(file_path) == ".xml" {
			files = append(files, file_path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		btree, error := BTree.Load(file)
		if btree == nil {
			common.GStdout.Console("load %v error %v", file, error)
			panic(error)
		}
		t.testbtree[btree.GetName()] = btree
	}
}

func (t *TestBTree) GetBTree(name string) *BTree.BTree {
	btree, ok := t.testbtree[name]
	if !ok {
		return nil
	}
	return btree
}

func init() {
	command.GCommand.RegCommand("listtest", listTest, "List All Test")
	command.GCommand.RegCommand("test", startTest, "Test")
	command.GCommand.RegCommand("endtest", endTest, "End Test")
	command.GCommand.RegCommand("playercount", showPlayerCount, "player count")
	command.GCommand.RegCommand("record", showTestRecord, "test record")
	command.GCommand.RegCommand("gc", startGC, "gc")
	GTestBTree = &TestBTree{}
	GTestBTree.testbtree = make(map[string]*BTree.BTree)
	GTestBTree.LoadAllTree("./btree")
	listTest(nil)
}

func listTest(str *common.StringParse) (err error) {
	common.GStdout.Console("=========================Test List========================================")
	for k := range GTestBTree.testbtree {
		common.GStdout.Console("%v", k)
	}
	common.GStdout.Console("===========================================================================")
	return
}

func startTest(str *common.StringParse) (err error) {
	if str.Len() < 4 {
		common.GStdout.Error("test param error")
		return
	}

	testtarget := str.GetString(1)
	accountStart := str.GetInt(2)
	accountEnd := accountStart + str.GetInt(3) - 1

	btree := GTestBTree.GetBTree(testtarget)
	if btree == nil {
		common.GStdout.Error("not found test :%v", testtarget)
		return nil
	}

	common.GStdout.Console("start [%v] player %v - %v ", testtarget, accountStart, accountEnd)

	for i := accountStart; i <= accountEnd; i++ {
		player := player.GPlayerMgr.CreatePlayer((uint64)(i))
		player.StartTest(btree)
		if i%30 == 0 {
			time.Sleep(5 * time.Second)
		}

		if i%100 == 0 {
			common.GStdout.Console("start player %v ", i)
		}
	}
	common.GStdout.Console("complete start [%v] player %v - %v ", testtarget, accountStart, accountEnd)
	return
}

func endTest(str *common.StringParse) (err error) {
	if str.Len() < 3 {
		common.GStdout.Error("test param error")
		return
	}

	accountStart := str.GetInt(1)
	accountEnd := accountStart + str.GetInt(2) - 1

	for i := accountStart; i <= accountEnd; i++ {
		player.GPlayerMgr.DelPlayer((uint64)(i))
	}
	return
}

// func startTest(str *common.StringParse) (err error) {
// 	if str.Len() < 2 {
// 		return nil
// 	}

// 	player := &player.Player{}
// 	runtime.SetFinalizer(player, PlayerRelease)
// 	player.Init(2, 2)
// 	player.Start()

// 	btree := GTestBTree.GetBTree(str.GetString(1))
// 	if btree == nil {
// 		common.GStdout.Error("not found tree :%v", str.GetString(1))
// 		return nil
// 	}

// 	player.StartTest(btree)
// 	time.Sleep(5 * time.Second)
// 	player.Release()
// 	return nil
// }

func startGC(str *common.StringParse) (err error) {
	runtime.GC()
	return nil
}

func showPlayerCount(str *common.StringParse) (err error) {
	count := player.GPlayerMgr.GetPlayerCount()
	common.GStdout.Console("player count : %v", count)
	return
}

func showTestRecord(str *common.StringParse) (err error) {
	BTree.PrintTaskRecord()
	return
}
