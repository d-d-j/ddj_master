package node

import (
	"net"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
)

type Listener struct {
	netListen	net.Listener
	idGenerator common.Int32Generator
}

func NewListener(service string) *Listener {
	tcpAddr, error := net.ResolveTCPAddr("tcp", service)
	if error != nil {
		log.Critical("Error: Could not resolve address")
	}

	log.Info("Listening on: ", tcpAddr.String())
	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if error != nil {
		log.Error(error)
	}

	list := new(Listener)
	list.netListen = netListen
	list.idGenerator = common.NewNodeIdGenerator()
	return list
}


func (list *Listener) WaitForNodes() {
	for {
		log.Info("Waiting for nodes")
		connection, error := list.netListen.Accept()
		if error != nil {
			log.Error("node error: ", error)
		} else {
			log.Info("Accept node: ", connection.RemoteAddr())
			newNode := NewNode(list.idGenerator.GetId(), connection)
			NodeManager.AddChan <- newNode
			go newNode.StartWork()
		}
	}
}

func (list *Listener) Close() {
	list.netListen.Close()
}
