package task

import (
	"ddj_Master/node"
	"testing"
)

func Test_getNodeForInsert_with_no_nodes_should_return_nil_and_error(t *testing.T) {

	balancer := node.NewLoadBalancer(nil)

	var worker TaskWorker
	worker.balancer = balancer
	node, actual := worker.getNodeForInsert()

	if node != nil {
		t.Error("Expected nil but got", node)
	}
	if actual == nil {
		t.Error("Expected not nil but got", actual)
	}
}

func Test_getNodeForInsert_should_return_node_provided_on_node_manager_channel(t *testing.T) {
	const NODE_ID = 1
	expected := node.NewNode(NODE_ID, nil, nil)
	nodes := make(map[int32]*node.Node)
	nodes[NODE_ID] = expected
	balancer := node.NewLoadBalancer(nodes)
	balancer.CurrentInsertNodeId = NODE_ID
	getNodeChan := make(chan node.GetNodeRequest, 1)

	var worker TaskWorker
	worker.balancer = balancer
	worker.getNodeChan = getNodeChan

	go func() {
		getNodeRequest := <-getNodeChan
		getNodeRequest.BackChan <- nodes[getNodeRequest.NodeId]
	}()

	actual, err := worker.getNodeForInsert()

	if err != nil {
		t.Error("Expected nil but got", err)
	}
	if actual.Id != expected.Id {
		t.Error("Expected ", expected, " but got", actual)
	}
}

func Test_getNodeForInsert_should_return_error_when_node_is_nil(t *testing.T) {
	const NODE_ID = 1
	expected := node.NewNode(NODE_ID, nil, nil)
	nodes := make(map[int32]*node.Node)
	nodes[NODE_ID] = expected
	balancer := node.NewLoadBalancer(nodes)
	balancer.CurrentInsertNodeId = NODE_ID
	getNodeChan := make(chan node.GetNodeRequest, 1)

	var worker TaskWorker
	worker.balancer = balancer
	worker.getNodeChan = getNodeChan

	go func() {
		getNodeRequest := <-getNodeChan
		getNodeRequest.BackChan <- nil
	}()

	actual, err := worker.getNodeForInsert()

	if actual != nil {
		t.Error("Expected nil but got", actual)
	}
	if err == nil {
		t.Error("Expected not nil but got", err)
	}
}
