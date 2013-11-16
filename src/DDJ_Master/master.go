package main

import (
	log "code.google.com/p/log4go"
	"config"
	"container/list"
	"dto"
	"fmt"
	"net"
	"node"
	"rest"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {

	log.Info("Start Master")

	nodeList := list.New()

	in := make(chan dto.Result)
	go node.IOHandler(rest.Channel.QueryChannel(), in, nodeList)

	log.Debug("Load configuration")
	cfg, err := config.Load()
	if err != nil {
		log.Critical("Problem with configuration: ", err)
	}
	log.Info(cfg)

	log.LoadConfiguration("log.cfg")

	portApi := fmt.Sprintf(":%d", cfg.Ports.Api)
	rest.StartApi(portApi)

	service := fmt.Sprintf("127.0.0.1:%d", cfg.Ports.NodeCommunication)
	tcpAddr, error := net.ResolveTCPAddr("tcp", service)
	if error != nil {
		log.Critical("Error: Could not resolve address")
	}

	log.Info("Listening on: ", tcpAddr.String())
	netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if error != nil {
		log.Error(error)
	}

	defer netListen.Close()

	WaitForNodes(netListen, nodeList, in)

}

func WaitForNodes(netListen net.Listener, nodes *list.List, in chan dto.Result) {
	for {
		log.Info("Waiting for nodes")
		connection, error := netListen.Accept()
		if error != nil {
			log.Error("node error: ", error)
		} else {
			log.Info("Accept node: ", connection.RemoteAddr())
			go node.NodeHandler(connection, in, nodes)
		}
	}
}
