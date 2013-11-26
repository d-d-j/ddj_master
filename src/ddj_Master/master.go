package main

import (
	log "code.google.com/p/log4go"
	"fmt"
	"ddj_Master/restApi"
	"ddj_Master/common"
	"ddj_Master/task"
	"ddj_Master/node"
)

func loadMasterConfiguration() *common.Config {
	log.Debug("Load configuration")
	cfg, err := common.LoadConfig()
	if err != nil {
		log.Critical("Problem with configuration: ", err)
	}
	log.Info(cfg)
	log.LoadConfiguration("log.cfg")
	return cfg
}

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.Info("Start Master")

	cfg := loadMasterConfiguration()

	// Start rest api server with tcp services for inserts and selects
	portNum := fmt.Sprintf("%d", cfg.Ports.RestApi)
	var server = restApi.Server{portNum}
	chanReq := server.StartApi()

	// Initialize task manager (balancer)
	go task.TaskManager.Manage()
	taskBal := task.NewBalancer(cfg.Constants.WorkersCount, cfg.Constants.JobForWorkerCount)
	go taskBal.Balance(chanReq)

	// Initialize node manager
	go node.NodeManager.Manage()
	infoChan := make(chan node.Info)
	nodeBal := node.NewLoadBalancer()
	go nodeBal.Balance(infoChan)

	// Initialize node listener
	service := fmt.Sprintf("127.0.0.1:%d", cfg.Ports.NodeCommunication)
	list := node.NewListener(service)
	defer list.Close()	// fire netListen.Close() when program ends

	// TODO: Wait for console instructions (q - quit for example)
}
