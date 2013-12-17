package node

import (
	"ddj_Master/common"
	"testing"
)

func AreEqual(expected, actual *LoadBalancer, t *testing.T) {
	if expected.CurrentInsertGpuId != actual.CurrentInsertGpuId {
		t.Error("CurrentInsertGpuId Got: ", actual.CurrentInsertGpuId, " when expected ", expected.CurrentInsertNodeId)
	}
	if expected.CurrentInsertNodeId != actual.CurrentInsertNodeId {
		t.Error("CurrentInsertNodeId Got: ", actual.CurrentInsertNodeId, " when expected ", expected.CurrentInsertNodeId)
	}
	if expected.timeout != actual.timeout {
		t.Error("timeout Got: ", actual.timeout, " when expected ", expected.CurrentInsertNodeId)
	}
}

func Test_NewLoadBalancer(t *testing.T) {
	expected := new(LoadBalancer)
	expected.CurrentInsertNodeId = common.CONST_UNINITIALIZED
	expected.CurrentInsertGpuId = common.CONST_UNINITIALIZED
	actual := NewLoadBalancer(0, nil, nil)
	AreEqual(expected, actual, t)
}

func Test_Update_With_Nil_Entry(t *testing.T) {
	actual := NewLoadBalancer(0, nil, nil)
	expected := NewLoadBalancer(0, nil, nil)
	actual.CurrentInsertNodeId = 1
	actual.CurrentInsertGpuId = 2
	actual.update(nil)
	AreEqual(expected, actual, t)
}

func Test_Update_With_No_Nodes_Cause_Balancer_Reset(t *testing.T) {
	nodes := make(map[int32]*Node)
	info := Info{1}
	balancer := NewLoadBalancer(0, nodes)
	expected := NewLoadBalancer(0, nil)
	actual.update(info)
	AreEqual(expected, actual, t)
}

func Test_Update_With_Two_Nodes_Each_With_One_GPU_Cause_Changing_Current_Node_After_Timeout(t *testing.T) {
	t.Error("Not implemented")
}

func Test_Update_With_One_Node_With_Two_GPUs_Cause_Changing_GPU_After_Timeout(t *testing.T) {
	t.Error("Not implemented")
}
