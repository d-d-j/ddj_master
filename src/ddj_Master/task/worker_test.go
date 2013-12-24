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
	ret := getNodeForInsert(req, balancer)
	actual := <-req.Response
	if ret != nil {
		t.Error("Expected nil but got", ret)
	}
	if actual.String() != expected.String() {
		t.Error("Expected ", expected, " but got", actual)
	}
}
