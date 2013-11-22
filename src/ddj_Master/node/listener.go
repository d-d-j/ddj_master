package node

import (
	"net"
	"ddj_Master/dto"
	"ddj_Master/node"
	log "code.google.com/p/log4go"
)

type Listener struct {
	NetListen	net.Listener
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
	list.NetListen = netListen
	return list
}


func (list *Listener) WaitForNodes(in chan dto.Result) {
	for {
		log.Info("Waiting for nodes")
		connection, error := list.NetListen.Accept()
		if error != nil {
			log.Error("node error: ", error)
		} else {
			log.Info("Accept node: ", connection.RemoteAddr())
			// TODO: Implement add new node
			go node.NodeHandler(connection, in, nodes)
		}
	}
}
