package task

import (
	log "code.google.com/p/log4go"
	"github.com/d-d-j/ddj_master/dto"
)

//This struct is used to get task with specific Id
type GetTaskRequest struct {
	TaskId   int64
	BackChan chan *dto.Task
}

var TaskManager = NewManager()

//Task Manager is similar to Node Manager and do same things but operates on tasks
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

//This methods handle request that came on Manager channels. It should be run as a gorutine
func (m *Manager) Manage() {
	log.Info("Task manager started managing")
	for {
		select {
		case get := <-m.GetChan:
			t := m.tasks[get.TaskId]
			log.Finest("Get task: %s", t)
			get.BackChan <- t
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
