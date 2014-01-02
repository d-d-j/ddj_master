package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
	"net"
)

type Listener struct {
	netListen    net.Listener
	idGenerator  common.Int32Generator
}

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

func (l *Listener) Close() {
	l.netListen.Close()
}
