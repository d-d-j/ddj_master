package node

import log "code.google.com/p/log4go"
import "ddj_Master/dto"

// TODO: get rid of this global variable
var NodeManager = NewManager()

type GetNodeRequest struct {
	NodeId   int32
	BackChan chan<- *Node
}

type Manager struct {
	nodes    map[int32]*Node
	AddChan  chan *Node
	GetChan  chan GetNodeRequest
	DelChan  chan int32
	QuitChan chan bool
	InfoChan chan []*dto.Info
}

func NewManager() *Manager {
	m := new(Manager)
	m.nodes = make(map[int32]*Node)
	m.AddChan = make(chan *Node)
	m.GetChan = make(chan GetNodeRequest)
	m.DelChan = make(chan int32)
	m.QuitChan = make(chan bool)
	m.InfoChan = make(chan []*dto.Info)
	return m
}

func (this *Manager) GetNodes() map[int32]*Node {
	return this.nodes
}

func (this *Manager) GetNodesLen() int {
	return len(this.nodes)
}

func (m *Manager) Manage() {
	log.Info("Node manager started managing")
	for {
		select {
		case get := <-m.GetChan:
			if node, ok := m.nodes[get.NodeId]; ok {
				get.BackChan <- node
			} else {
				panic("Node not found")
			}
		case newNode := <-m.AddChan:
			m.nodes[newNode.Id] = newNode
		case closedNodeId := <-m.DelChan:
			delete(m.nodes, closedNodeId)
			log.Info("Node manager deleted Node #%d from nodes", closedNodeId)
		case <-m.QuitChan:
			log.Info("Node manager stopped managing")
			return
		}
	}
}

func (this *Manager) SendToAllNodes(message []byte) {
	log.Debug("Sending message to all %d", this.GetNodesLen(), " nodes")
	// SEND MESSAGE TO ALL NODES
	for _, n := range this.nodes {
		log.Finest("Sending message to node #%d", n.Id)
		n.Incoming <- message
	}
}
