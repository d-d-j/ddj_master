package node

import (
	log "code.google.com/p/log4go"
	"github.com/d-d-j/ddj_master/common"
	"github.com/d-d-j/ddj_master/dto"
	"net"
)

//Listener is responsible for accepting new connections.
type Listener struct {
	netListen   net.Listener
	idGenerator common.Int32Generator
}

//Listener's constructor take service on which listener should listen
func NewListener(service string) *Listener {
	_, error := net.ResolveTCPAddr("tcp", service)
	if error != nil {
		log.Critical("Error: Could not resolve address")
	}

	log.Info("Listening on: %v", service)
	netListen, error := net.Listen("tcp", service)
	if error != nil {
		log.Error(error)
	}

	l := new(Listener)
	l.netListen = netListen
	l.idGenerator = common.NewNodeIdGenerator()
	return l
}

//This method should be run as a gorutine it task channel  that is required by Node. When new node try to connect this method
//accept connection and create new Node using it's constructorand run it
func (l *Listener) WaitForNodes(getTaskChannel chan dto.GetTaskRequest) {
	for {
		log.Info("Waiting for nodes")
		connection, error := l.netListen.Accept()
		if error != nil {
			log.Error("node error: ", error)
		} else {
			log.Info("Accept node: ", connection.RemoteAddr())
			newNode := NewNode(l.idGenerator.GetId(), connection, getTaskChannel)
			NodeManager.AddChan <- newNode
			go newNode.StartWork(NodeManager.InfoChan)
		}
	}
}

//Stop listener
func (l *Listener) Close() {
	l.netListen.Close()
}
