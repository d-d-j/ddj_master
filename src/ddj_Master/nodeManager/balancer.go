package nodeManager

import (
	"container/list"
	log "code.google.com/p/log4go"
)

type LoadBalancer struct {
	Nodes 				list.List
	CurrentInsertNode	*Node
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

}

func (lb *LoadBalancer) AddNode(node *Node) {

	lb.Nodes.PushBack(node)
}

func (lb *LoadBalancer) RemoveNode(node *Node) {

	lb.Nodes.Remove(node)
}



