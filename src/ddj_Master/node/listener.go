package node

import (
	"net"
	"ddj_Master/dto"
	"ddj_Master/node"
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
	list.idGenerator = common.NodeIdGenerator{}
	return list
}


func (list *Listener) WaitForNodes(in chan dto.Result) {
	for {
		log.Info("Waiting for nodes")
		connection, error := list.netListen.Accept()
		if error != nil {
			log.Error("node error: ", error)
		} else {
			log.Info("Accept node: ", connection.RemoteAddr())
			// TODO: Instead of 0 there should be slice of GPUIds of new Node
			NodeManager.AddChan <- node.NewNode(list.idGenerator.getId(), 0, connection)
		}
	}
}
