package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
)

type LoadBalancer struct {
	CurrentInsertNodeId int32
	CurrentInsertGpuId  int32
	timeout             int32
}

func NewLoadBalancer(timeout int32) *LoadBalancer {
	log.Debug("Load balancer constructor [START]")
	lb := new(LoadBalancer)
	lb.reset()
	lb.timeout = timeout
	lb.update(nil)
	log.Debug("Load balancer constructor [END]")
	return lb
}

func (lb *LoadBalancer) reset() {
	lb.CurrentInsertNodeId = common.CONST_UNINITIALIZED
	lb.CurrentInsertGpuId = common.CONST_UNINITIALIZED
}

func (lb *LoadBalancer) Balance(info <-chan Info) {

	log.Info("Node manager balancer started")
	for {
		select {
		case nodeInfo := <-info:
			lb.update(&nodeInfo)
		}
	}
}

func (lb *LoadBalancer) update(newInfo *Info) {
	if newInfo == nil {
		lb.reset()
		return
	}

	// DIRECTS ALL INCOMING DATA TO NEW NODE's FIRST GPU
	ch := make(chan *Node)
	NodeManager.GetChan <- GetNodeRequest{newInfo.nodeId, ch}
	var n *Node = <-ch
	lb.CurrentInsertNodeId = n.Id
	lb.CurrentInsertGpuId = n.GpuIds[0]
	log.Debug("Node with id ", n.Id, " and GPU ", n.GpuIds[0], " set to current")
	// TODO: Write balance function for nodes
}
