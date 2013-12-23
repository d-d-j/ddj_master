package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
)

type LoadBalancer struct {
	CurrentInsertNodeId int32
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

	if this.IsUnitialized() {
		for _, node := range this.nodes {
			this.CurrentInsertNodeId = node.Id
			node.PreferredDeviceId = node.GpuIds[0]
			break
		}
		return
	}

	const (
		CurrentNodePenalty = 10
	)

	bestNodeId := common.CONST_UNINITIALIZED
	bestRank := -(int(^uint(0) >> 1))

	//TODO: Calculate full rank when data will be available
	//Now we have no info about card load, ram, proc etc
	for id, node := range this.nodes {
		rank := 0
		if this.CurrentInsertNodeId == node.Id {
			rank -= CurrentNodePenalty
			for gpuId := range node.GpuIds {
				if node.PreferredDeviceId != int32(gpuId) {
					node.PreferredDeviceId = int32(gpuId)
					break
				}
			}
		}

		if rank > bestRank {
			bestNodeId = int(id)
			bestRank = rank
			log.Debug(bestRank, id)
		}
	}

	this.CurrentInsertNodeId = int32(bestNodeId)
}

func (this *LoadBalancer) IsUnitialized() bool {
	return this.CurrentInsertNodeId == common.CONST_UNINITIALIZED
}
