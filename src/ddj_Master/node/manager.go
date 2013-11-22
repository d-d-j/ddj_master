package node

import log "code.google.com/p/log4go"

// TODO: get rid of this global variable
var NodeManager = NewManager()

type GetNodeRequest struct {
	nodeId		int32
	backChan	chan<- *Node
}

type Manager struct {
	nodes		map[int32]Node
	addChan		<-chan *Node
	getChan		<-chan GetNodeRequest
	delChan		<-chan int64
	quitChan	<-chan bool
}

func NewManager() *Manager {
	m := new(Manager)
	m.tasks = make(map[int64]*Node)
	m.addChan = make(<-chan *Node)
	m.getChan = make(<-chan GetNodeRequest)
	m.delChan = make(<-chan int64)
	m.quitChan = make(<-chan bool)
	return m
}

func (m *Manager) Manage() {
	log.Info("Node manager started managing")
	for {
		select {
		case get := <-m.getChan:
			get.backChan <- m.nodes[get.nodeId]
		case add := <-m.addChan:
			m.nodes[add.Id] = add
		case del := <-m.delChan:
			delete(m.nodes, del)
		case q := <-m.quitChan:
			log.Info("Node manager stopped managing")
			return
		}
	}
}

