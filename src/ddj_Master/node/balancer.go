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
	lb := new(LoadBalancer)
	lb.reset()
	lb.timeout = timeout
	lb.nodes = nodes
	lb.update(nil)
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

	if this.CurrentInsertGpuId == common.CONST_UNINITIALIZED || this.CurrentInsertNodeId == common.CONST_UNINITIALIZED {
		for _, node := range this.nodes {
			this.CurrentInsertNodeId = node.Id
			this.CurrentInsertGpuId = node.GpuIds[0]
			break
		}
	}
}
