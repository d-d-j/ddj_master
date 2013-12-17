package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
)

type LoadBalancer struct {
	CurrentInsertNodeId int32
	CurrentInsertGpuId  int32
	timeout             int32
	nodes               map[int32]*Node
}

func NewLoadBalancer(timeout int32, nodes map[int32]*Node) *LoadBalancer {
	log.Debug("Load balancer constructor [START]")
	lb := new(LoadBalancer)
	lb.reset()
	lb.timeout = timeout
	lb.nodes = nodes
	lb.update(nil)
	log.Debug("Load balancer constructor [END]")
	return lb
}

func (this *LoadBalancer) reset() {
	this.CurrentInsertNodeId = common.CONST_UNINITIALIZED
	this.CurrentInsertGpuId = common.CONST_UNINITIALIZED
}

func (this *LoadBalancer) Balance(info <-chan Info) {

	log.Info("Node manager balancer started")
	for {
		select {
		case nodeInfo := <-info:
			this.update(&nodeInfo)
		}
	}
}

func (this *LoadBalancer) update(newInfo *Info) {
	if newInfo == nil {
		this.reset()
		return
	}

}
