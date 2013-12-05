package node

import (
	"net"
	log "code.google.com/p/log4go"
	"ddj_Master/common"
)

type Listener struct {
	netListen		net.Listener
	idGenerator 	common.Int32Generator
	balancerChan	chan<- Info
}

func NewListener(service string, nodeInfoChannel chan<- Info) *Listener {
	tcpAddr, error := net.ResolveTCPAddr("tcp", service)
	if error != nil {
		log.Critical("Error: Could not resolve address")
	}

	log.Info("Listening on: ", tcpAddr.String())
	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if error != nil {
		log.Error(error)
	}

	l := new(Listener)
	l.netListen = netListen
	l.idGenerator = common.NewNodeIdGenerator()
	l.balancerChan = nodeInfoChannel
	return l
}


func (l *Listener) WaitForNodes() {
	for {
		log.Info("Waiting for nodes")
		connection, error := l.netListen.Accept()
		if error != nil {
			log.Error("node error: ", error)
		} else {
			log.Info("Accept node: ", connection.RemoteAddr())
			newNode := NewNode(l.idGenerator.GetId(), connection)
			NodeManager.AddChan <- newNode
			go newNode.StartWork(l.balancerChan)
		}
	}
}

func (l *Listener) Close() {
	l.netListen.Close()
}
