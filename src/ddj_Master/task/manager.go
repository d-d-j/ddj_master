package task

import log "code.google.com/p/log4go"

// TODO: get rid of this global variable
var TaskManager = NewManager()

type GetTaskRequest struct {
	taskId		int64
	backChan	chan<- *Task
}

type Manager struct {
	tasks		map[int64]*Task
	addChan		<-chan *Task
	getChan		<-chan GetTaskRequest
	delChan		<-chan int64
	quitChan	<-chan bool
}

func NewManager() *Manager {
	m := new(Manager)
	m.tasks = make(map[int64]*Task)
	m.addChan = make(<-chan *Task)
	m.getChan = make(<-chan GetTaskRequest)
	m.delChan = make(<-chan int64)
	m.quitChan = make(<-chan bool)
	return m
}

func (m *Manager) Manage() {
	log.Info("Task manager started managing")
	for {
		select {
			case get := <-m.getChan:
				get.backChan <- m.tasks[get.taskId]
			case add := <-m.addChan:
				m.tasks[add.Id] = add
			case del := <-m.delChan:
				delete(m.tasks, del)
			case q := <-m.quitChan:
				log.Info("Task manager stopped managing")
				return
		}
	}
}
