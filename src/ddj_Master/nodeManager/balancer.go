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

// Creates a new Node object for each new connection using the Id sent by the Node,
// then starts the NodeSender and NodeReader goroutines to handle the IO
func (lb *LoadBalancer) AddNode(node *Node) {

	lb.Nodes.PushBack(node)
}



