package node

type LoadBalancer interface {
	GetBestNodeForSeries(series int) *Node
}

