package main

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/node"
	"ddj_Master/reduce"
	"ddj_Master/restApi"
	"ddj_Master/task"
	"fmt"
	"runtime"
)

func loadMasterConfiguration() *common.Config {
	log.Debug("Load configuration")
	cfg, err := common.LoadConfig()
	if err != nil {
		log.Critical("Problem with configuration: ", err)
		panic("Wrong configuration")
	}
	log.Info(cfg)
	return cfg
}

// Main: Starts a TCP server and waits infinitely for connections
func main() {
	log.LoadConfiguration("log.cfg")
	log.Info("Start Master")

	cfg := loadMasterConfiguration()

	log.Info("Setting go cpu number to ", cfg.Constants.CpuNumber, " success: ", runtime.GOMAXPROCS(cfg.Constants.CpuNumber))

	// Start rest api server with tcp services for inserts and selects
	portNum := fmt.Sprintf(":%d", cfg.Ports.RestApi)
	var server = restApi.Server{Port: portNum}
	chanReq := server.StartApi()

	// Initialize node manager
	log.Info("Initialize node manager")
	go node.NodeManager.Manage()
	nodeBal := node.NewLoadBalancer(node.NodeManager.GetNodes())
	go nodeBal.Balance(node.NodeManager.InfoChan)

	// Initialize reduce factory
	log.Info("Initialize reduce factory")
	reduce.Initialize()

	// Initialize task manager (balancer)
	log.Info("Initialize task manager")
	go task.TaskManager.Manage()
	taskBal := task.NewBalancer(cfg.Constants.WorkersCount, cfg.Constants.JobForWorkerCount, nodeBal)
	go taskBal.Balance(chanReq, cfg.Balancer.Timeout)

	// Initialize node listener
	log.Info("Initialize node listener")
	service := fmt.Sprintf(":%d", cfg.Ports.NodeCommunication)
	log.Debug(service)
	list := node.NewListener(service)
	go list.WaitForNodes(task.TaskManager.GetChan)
	defer list.Close() // fire netListen.Close() when program ends

	// TODO: Wait for console instructions (q - quit for example)
	// Wait for some input end exit (only for now)
	var i int
	fmt.Scanf("%d", &i)
}
