package task

import log "code.google.com/p/log4go"

// TODO: get rid of this global variable
var TaskManager = NewManager()

type GetTaskRequest struct {
	TaskId		int64
	BackChan	chan<- *Task
}

type Manager struct {
	tasks		map[int64]*Task
	AddChan		chan *Task
	GetChan		chan GetTaskRequest
	DelChan		chan int64
	QuitChan	chan bool
}

func NewManager() *Manager {
	m := new(Manager)
	m.tasks = make(map[int64]*Task)
	m.AddChan = make(chan *Task)
	m.GetChan = make(chan GetTaskRequest)
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
				m.tasks[add.Id] = add
			case del := <-m.DelChan:
				delete(m.tasks, del)
			case <-m.QuitChan:
				log.Info("Task manager stopped managing")
				return
		}
	}
}