package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
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

func (this *LoadBalancer) Balance(info <-chan []*dto.Info) {

	log.Info("Node manager balancer started")
	for {
		select {
		case nodeInfo := <-info:
			log.Debug("Balancer got nodeInfos: %v", nodeInfo)
			this.update(nodeInfo)
		}
	}
}

func (this *LoadBalancer) update(newInfos []*dto.Info) {
	if newInfos == nil {
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


	bestNodeId := this.chooseTheBestNode(newInfos)

	this.CurrentInsertNodeId = int32(bestNodeId)
}

func (this *LoadBalancer) chooseTheBestNode(nodeInfos []*dto.Info) int {
	bestNodeId := common.CONST_UNINITIALIZED
	bestRank := common.CONST_INT_MIN_VALUE

	for id, node := range this.nodes {

		rank := this.calculateNodeRank(node)

		if rank > bestRank {
			bestNodeId = int(id)
			bestRank = rank
			log.Debug(bestRank, id)
		}
	}

	return bestNodeId
}

func (this *LoadBalancer) calculateNodeRank(node *Node) int {
	rank := 0
	const (
		CurrentNodePenalty = 10
	)
	if this.CurrentInsertNodeId == node.Id {
		rank -= CurrentNodePenalty
		for gpuId := range node.GpuIds {
			if node.PreferredDeviceId != int32(gpuId) {
				node.PreferredDeviceId = int32(gpuId)
				break
			}
		}
	}
	return rank
}

func (this *LoadBalancer) IsUnitialized() bool {
	return this.CurrentInsertNodeId == common.CONST_UNINITIALIZED
}
