package task

import (
	log "code.google.com/p/log4go"
	"ddj_Master/dto"
)

type GetTaskRequest struct {
	TaskId   int64
	BackChan chan *dto.Task
}

// TODO: get rid of this global variable
var TaskManager = NewManager()

type Manager struct {
	tasks    map[int64]*dto.Task
	AddChan  chan *dto.Task
	GetChan  chan dto.GetTaskRequest
	DelChan  chan int64
	QuitChan chan bool
}

func NewManager() *Manager {
	m := new(Manager)
	m.tasks = make(map[int64]*dto.Task)
	m.AddChan = make(chan *dto.Task)
	m.GetChan = make(chan dto.GetTaskRequest)
	m.DelChan = make(chan int64)
	m.QuitChan = make(chan bool)
	return m
}

func (m *Manager) Manage() {
	log.Info("Task manager started managing")
	for {
		select {
		case get := <-m.GetChan:
			get.BackChan <- m.tasks[get.TaskId]
		case add := <-m.AddChan:
			log.Finest("Add new task: %s", add)
			m.tasks[add.Id] = add
		case del := <-m.DelChan:
			delete(m.tasks, del)
		case <-m.QuitChan:
			log.Info("Task manager stopped managing")
			return
		}
	}
}
