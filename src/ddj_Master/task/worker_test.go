package task

import (
	"ddj_Master/common"
	"ddj_Master/dto"
	"ddj_Master/node"
	"testing"
)

func Test_getNodeForInsert_with_no_nodes_should_return_nil_and_error_response_on_Response_channel(t *testing.T) {
	expected := dto.NewRestResponse("No node connected", common.TASK_UNINITIALIZED, nil)
	balancer := node.NewLoadBalancer(0, nil)
	var req dto.RestRequest
	req.Response = make(chan *dto.RestResponse, 1)
	ret := getNodeForInsert(req, balancer, nil)
	actual := <-req.Response
	if ret != nil {
		t.Error("Expected nil but got", ret)
	}
	if actual.String() != expected.String() {
		t.Error("Expected ", expected, " but got", actual)
	}
}

func Test_getNodeForInsert_should_return_node_provided_on_node_manager_channel(t *testing.T) {
	const NODE_ID = 1
	expected := node.NewNode(NODE_ID, nil, nil)
	nodes := make(map[int32]*node.Node)
	nodes[NODE_ID] = expected
	balancer := node.NewLoadBalancer(0, nodes)
	balancer.CurrentInsertNodeId = NODE_ID

	getNodeChan := make(chan node.GetNodeRequest, 1)
	go func() {
		nodeRequest := <-getNodeChan
		nodeRequest.BackChan <- nodes[nodeRequest.NodeId]
	}()

	var req dto.RestRequest
	actual := getNodeForInsert(req, balancer, getNodeChan)

	if actual.Id != expected.Id {
		t.Error("Expected nil but got", actual)
	}
}
