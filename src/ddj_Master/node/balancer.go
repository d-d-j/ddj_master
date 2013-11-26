package node

import (
	log "code.google.com/p/log4go"
	"math/rand"
)

type LoadBalancer struct {
	CurrentInsertNodeId	int32
	CurrentInsertGpuId	int32
}

func NewLoadBalancer() *LoadBalancer {
	log.Debug("Load balancer constructor [START]")
	lb := new(LoadBalancer)
	lb.CurrentInsertNodeId = -1
	lb.CurrentInsertGpuId = -1
	lb.update(nil)
	log.Debug("Load balancer constructor [END]")
	return lb
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
	if(newInfo == nil) {
		/* Choose random node for a start */
		len := int32(len(NodeManager.nodes))
		if len > 0 {
			number := rand.Int31n(len)
			for k, n := range NodeManager.nodes {
				if k == number {
					lb.CurrentInsertNodeId = n.Id
					// TODO: Choose also GpuId
				}
			}
		}
	}
	// TODO: Write balance function for nodes
}



