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
	actual := NewLoadBalancer(0, nil)
	AreEqual(expected, actual, t)
}

func Test_Update_With_Nil_Entry_Cause_Balancer_Reset(t *testing.T) {
	actual := NewLoadBalancer(0, nil)
	expected := NewLoadBalancer(0, nil)
	actual.CurrentInsertNodeId = 1
	actual.CurrentInsertGpuId = 2
	actual.update(nil)
	AreEqual(expected, actual, t)
}

func Test_Update_With_No_Nodes_Cause_Balancer_Reset(t *testing.T) {
	nodes := make(map[int32]*Node)
	info := &Info{1}
	actual := NewLoadBalancer(0, nodes)
	expected := NewLoadBalancer(0, nil)
	actual.update(info)
	AreEqual(expected, actual, t)
}

func Test_Update_Unitialized_Balancer_Set_It_To_First_Node_And_GPU(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil)
	nodes[1].GpuIds = []int32{0}
	nodes[2] = NewNode(2, nil)
	nodes[2].GpuIds = []int32{0}

	lb := NewLoadBalancer(0, nodes)
	t.Log(nodes[1].GpuIds)
	info := &Info{1}
	lb.update(info)
	if lb.CurrentInsertNodeId != 1 || lb.CurrentInsertGpuId != 0 {
		t.Errorf("Wrong card selected. Expected #%d:%d but get #%d:%d", 1, 0, lb.CurrentInsertNodeId, lb.CurrentInsertGpuId)
	}
	t.Log(nodes[1].GpuIds)
}

func Test_Update_With_One_Node_With_Two_GPUs_Cause_Changing_GPU_After_Timeout(t *testing.T) {
	t.Error("Not implemented")
}
