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
	info := &Info{1, MemoryInfo{1, 1, 1, 1}}
	actual := NewLoadBalancer(0, nodes)
	expected := NewLoadBalancer(0, nil)
	actual.update(info)
	AreEqual(expected, actual, t)
}

func Test_Update_Unitialized_Balancer_Set_It_To_First_Node_And_GPU(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1, 2}
	nodes[2] = NewNode(2, nil, nil)
	nodes[2].GpuIds = []int32{0, 1, 2}

	lb := NewLoadBalancer(0, nodes)
	info := &Info{1, MemoryInfo{1, 1, 1, 1}}
	lb.update(info)
	if lb.CurrentInsertNodeId != 1 || lb.CurrentInsertGpuId != 0 {
		t.Errorf("Wrong card selected. Expected #%d:%d but get #%d:%d", 1, 0, lb.CurrentInsertNodeId, lb.CurrentInsertGpuId)
	}
}

func Test_Update_With_One_Node_With_Two_GPUs_Cause_Changing_GPU(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1, 2}

	lb := NewLoadBalancer(0, nodes)
	lb.CurrentInsertNodeId = 1
	lb.CurrentInsertGpuId = 0
	info := &Info{1, MemoryInfo{1, 1, 1, 1}}

	for i := 0; i < 3; i++ {
		lb.update(info)
		if lb.CurrentInsertGpuId == int32(i) {
			t.Errorf("Wrong card selected. Expected #%d but get #%d", i, lb.CurrentInsertGpuId)
		}
	}
}

func Test_Update_With_Two_Nodes_Cause_Changing_Node(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1, 2}
	nodes[2] = NewNode(2, nil, nil)
	nodes[2].GpuIds = []int32{0, 1, 2}

	lb := NewLoadBalancer(0, nodes)
	lb.CurrentInsertNodeId = 1
	lb.CurrentInsertGpuId = 0
	info := &Info{1, MemoryInfo{1, 1, 1, 1}}

	for i := 0; i < 4; i++ {
		nodeId := i%2 + 1
		lb.update(info)
		if lb.CurrentInsertNodeId == int32(nodeId) {
			t.Errorf("Wrong card selected. Expected #%d but get #%d", nodeId, lb.CurrentInsertNodeId)
		}
	}
}

func Test_Update_Called_3_Times_Will_Fire_All_3_Nodes_And_GPUs(t *testing.T) {
	t.Skip("Now we don't support this feature. We need to think how node balancer should work")
	nodes := make(map[int32]*Node)
	for i := int32(0); i < 3; i++ {
		nodes[i] = NewNode(i, nil, nil)
		nodes[i].GpuIds = []int32{0, 1, 2}
	}

	lb := NewLoadBalancer(0, nodes)
	info := &Info{1, MemoryInfo{1, 1, 1, 1}}

	actual := [9]int32{}
	expected := [9]int32{0, 1, 2, 10, 11, 12, 20, 21, 22}
	for i := 0; i < 9; i++ {
		lb.update(info)
		actual[i] = lb.CurrentInsertNodeId*10 + lb.CurrentInsertGpuId
	}
	if expected != actual {
		t.Errorf("Expected: ", expected, "but got: ", actual)
	}
}
