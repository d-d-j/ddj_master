package node

import (
	log "code.google.com/p/log4go"
	"ddj_Master/common"
	"ddj_Master/dto"
)

type LoadBalancer struct {
	CurrentInsertNodeId int32
	nodes               map[int32]*Node
}

func NewLoadBalancer(nodes map[int32]*Node) *LoadBalancer {
	lb := new(LoadBalancer)
	lb.reset()
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

	bestNodeId := this.chooseTheBestNode(newInfos)

	this.CurrentInsertNodeId = int32(bestNodeId)

	log.Info("Insert Node Id is now %d", this.CurrentInsertNodeId)
}

func (this *LoadBalancer) chooseTheBestNode(nodeInfos []*dto.Info) int {
	bestNodeId := common.CONST_UNINITIALIZED
	bestRank := common.CONST_UINT64_MAX_VALUE

	for _, node := range this.nodes {
		rank := this.calculateNodeRank(node, nodeInfos)

		if rank < bestRank {
			bestNodeId = int(node.Id)
			bestRank = rank
			log.Debug("calculated best rank: %d,  node id: %d", bestRank, node.Id)
		}
	}

	return bestNodeId
}

func (this *LoadBalancer) calculateNodeRank(node *Node, nodeInfos []*dto.Info) uint64 {
	nodeRank := uint64(0)
	bestGpuRank := common.CONST_UINT64_MAX_VALUE;
	for _, info := range nodeInfos {
		if info.NodeId != node.Id {
			continue
		}
		gpuRank := info.MemoryInfo.DBMemoryFree

		if gpuRank < bestGpuRank {
			bestGpuRank = gpuRank
			node.PreferredDeviceId = info.GpuId
			nodeRank = gpuRank
		}
	}

	log.Debug("calculated rank: %d, for node: %d, changed deviceId to %d", nodeRank, node.Id, node.PreferredDeviceId)
	return nodeRank
}

func (this *LoadBalancer) IsUnitialized() bool {
	return this.CurrentInsertNodeId == common.CONST_UNINITIALIZED
}
