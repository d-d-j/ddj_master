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
	log.Debug("#", this.CurrentInsertNodeId)
	if this.CurrentInsertGpuId == common.CONST_UNINITIALIZED || this.CurrentInsertNodeId == common.CONST_UNINITIALIZED {
		for _, node := range this.nodes {
			this.CurrentInsertNodeId = node.Id
			this.CurrentInsertGpuId = node.GpuIds[0]
			break
		}
		return
	}

	const (
		CurrentNodePenalty = 10
	)
	bestNodeId := common.CONST_UNINITIALIZED
	bestGpuId := common.CONST_UNINITIALIZED
	bestRank := -(int(^uint(0) >> 1))
	for id, node := range this.nodes {
		rank := 0
		if this.CurrentInsertNodeId == node.Id {
			rank -= CurrentNodePenalty
			for gpuId := range node.GpuIds {
				if this.CurrentInsertGpuId != int32(gpuId) {
					bestGpuId = gpuId
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

	log.Debug(bestNodeId)
	this.CurrentInsertGpuId = int32(bestNodeId)
	if bestGpuId == common.CONST_UNINITIALIZED {
		this.CurrentInsertGpuId = this.nodes[int32(bestNodeId)].GpuIds[0]
	} else {
		this.CurrentInsertGpuId = int32(bestGpuId)
	}
}
