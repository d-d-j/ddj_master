package main

import (
	log "code.google.com/p/log4go"
	"container/list"
	"fmt"
	"net"
	"ddj_Master/dto"
	"ddj_Master/restApi"
	"ddj_Master/taskManager"
	"ddj_Master/common"
)

// Main: Starts a TCP server and waits infinitely for connections
func main() {

	log.Info("Start Master")

	// Load master configuration
	log.Debug("Load configuration")
	cfg, err := common.LoadConfig()
	if err != nil {
		log.Critical("Problem with configuration: ", err)
	}
	log.Info(cfg)
	log.LoadConfiguration("log.cfg")

	// Start rest api server with tcp services for inserts and selects
	portNum := fmt.Sprintf(":%d", cfg.Ports.RestApi)
	var server = restApi.Server{portNum}
	chanReq := server.StartApi()
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
	defer netListen.Close()	// fire netListen.Close() when program ends

	// Initialize task manager (balancer)
	bal := taskManager.NewBalancer(cfg.Constants.WorkersCount)
	go bal.balance(chanReq)

	// TODO: Initialize node manager
	nodeList := list.New()
	WaitForNodes(netListen, nodeList, in)

	// TODO: Wait for console instructions (q - quit for example)
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
