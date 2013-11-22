package main

import (
	log "code.google.com/p/log4go"
	"fmt"
	"net"
	"ddj_Master/restApi"
	"ddj_Master/common"
	"ddj_Master/task"
	"ddj_Master/node"
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

	// Initialize task manager (balancer)
	go task.TaskManager.Manage()
	taskBal := task.NewBalancer(cfg.Constants.WorkersCount)
	go taskBal.balance(chanReq)

	// Initialize node manager
	go node.NodeManager.Manage()
	infoChan := make(chan Info)
	nodeBal := node.NewLoadBalancer()
	go nodeBal.balance(infoChan)

	// Initialize node listener
	list := node.NewListener(service)
	defer list.NetListen.Close()	// fire netListen.Close() when program ends

	// TODO: Wait for console instructions (q - quit for example)
}
