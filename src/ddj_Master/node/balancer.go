package node

import (
	log "code.google.com/p/log4go"
)

type LoadBalancer struct {
	CurrentInsertNodeId	int32
	CurrentInsertGpuId	int32
}

func NewLoadBalancer() *LoadBalancer {
	lb := new(LoadBalancer)
	lb.update(nil)
	return lb
}

func (lb *LoadBalancer) balance(info <-chan Info) {

	log.Info("Node manager balancer started")
	for {
		select {
		case nodeInfo := <-info:
			lb.update(nodeInfo)
		}
	}
}

func (lb *LoadBalancer) update(newInfo Info) {
	if(newInfo == nil) {
		// TODO: Reset balance (choose random node ang his GPU if any)
	}
	// TODO: Write balance function for nodes
}



