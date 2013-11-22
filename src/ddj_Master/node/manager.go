package node

import log "code.google.com/p/log4go"

// TODO: get rid of this global variable
var NodeManager = NewManager()

type GetNodeRequest struct {
	NodeId		int32
	BackChan	chan<- *Node
}

type Manager struct {
	nodes		map[int32]Node
	AddChan		<-chan *Node
	GetChan		<-chan GetNodeRequest
	DelChan		<-chan int64
	QuitChan	<-chan bool
}

func NewManager() *Manager {
	m := new(Manager)
	m.tasks = make(map[int64]*Node)
	m.AddChan = make(<-chan *Node)
	m.GetChan = make(<-chan GetNodeRequest)
	m.DelChan = make(<-chan int64)
	m.QuitChan = make(<-chan bool)
	return m
}

func (m *Manager) Manage() {
	log.Info("Node manager started managing")
	for {
		select {
		case get := <-m.GetChan:
			if val, ok := m.nodes["route"]; ok {
				get.BackChan <- m.nodes[get.NodeId]
			} else {
				get.BackChan <- nil
			}
		case add := <-m.AddChan:
			m.nodes[add.Id] = add
		case del := <-m.DelChan:
			delete(m.nodes, del)
		case q := <-m.QuitChan:
			log.Info("Node manager stopped managing")
			return
		}
	}
}

