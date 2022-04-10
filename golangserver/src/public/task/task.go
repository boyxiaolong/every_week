package task

import (
	"sync"
)

func MakeTask(count uint32) *Task {
	return &Task{
		RunTask: make(chan func(), count),
		//stopChan:  make(chan struct{}),
		IsRunning: false,
	}
}

type Task struct {
	RunTask   chan func()
	stopChan  chan struct{}
	stopWait  sync.WaitGroup
	IsRunning bool
	once      sync.Once
}

func (m *Task) Run() {
	for {
		select {
		case call := <-m.RunTask:
			call()
		case <-m.stopChan:
			m.stopWait.Done()
			//common.GStdout.Console("Run ============================ done")
			return
		}
	}
}

func (m *Task) AddTask(test func()) {
	if !m.IsRunning {
		return
	}
	m.RunTask <- test
}

func (m *Task) Start() {
	if m.IsRunning {
		return
	}

	m.stopChan = make(chan struct{})
	m.IsRunning = true
	go m.Run()
}

func (m *Task) StopOnce() {
	m.once.Do(m.Stop)
}

func (m *Task) Stop() {
	m.IsRunning = false
	m.stopWait.Add(1)
	close(m.stopChan)
	m.stopWait.Wait()
}
