package node

import (
	"ddj_Master/common"
	"testing"
	"ddj_Master/dto"
)

func AreEqual(expected, actual *LoadBalancer, t *testing.T) {
	if expected.CurrentInsertNodeId != actual.CurrentInsertNodeId {
		t.Error("CurrentInsertNodeId Got: ", actual.CurrentInsertNodeId, " when expected ", expected.CurrentInsertNodeId)
	}
}

func Test_NewLoadBalancer(t *testing.T) {
	expected := new(LoadBalancer)
	expected.CurrentInsertNodeId = common.CONST_UNINITIALIZED
	actual := NewLoadBalancer(nil)
	AreEqual(expected, actual, t)
}

func Test_Update_With_Nil_Entry_Cause_Balancer_Reset(t *testing.T) {
	actual := NewLoadBalancer(nil)
	expected := NewLoadBalancer(nil)
	actual.CurrentInsertNodeId = 1
	actual.update(nil)
	AreEqual(expected, actual, t)
}

func Test_Update_With_No_Nodes_Cause_Balancer_Reset(t *testing.T) {
	nodes := make(map[int32]*Node)
	info := []*dto.Info{&dto.Info{1, dto.MemoryInfo{GpuId: 1, MemoryTotal: 1, MemoryFree: 1, GpuMemoryTotal:  1, GpuMemoryFree: 1, DBMemoryFree: 2}}}
	actual := NewLoadBalancer(nodes)
	expected := NewLoadBalancer(nil)
	actual.update(info)
	AreEqual(expected, actual, t)
}

func Test_Update_With_One_Node_With_One_GPU_Node_And_GPU_Are_Set(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0}

	lb := NewLoadBalancer(nodes)
	lb.CurrentInsertNodeId = 0
	nodes[1].PreferredDeviceId = 0
	info := []*dto.Info{&dto.Info{1, dto.MemoryInfo{GpuId: 1, MemoryTotal: 1, MemoryFree: 1, GpuMemoryTotal:  1, GpuMemoryFree: 1, DBMemoryFree: 2}}}


	lb.update(info)
	if lb.CurrentInsertNodeId != 1 {
		t.Errorf("Wrong node selected. Expected #%d but get #%d", 1, lb.CurrentInsertNodeId)
	}
	if nodes[1].PreferredDeviceId != int32(1) {
		t.Errorf("Wrong card selected. Expected #%d but get #%d", 1, nodes[1].PreferredDeviceId)
	}

}

func Test_Update_With_One_Node_With_Two_GPUs_Cause_Changing_GPU(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1, 2}

	lb := NewLoadBalancer(nodes)
	lb.CurrentInsertNodeId = -1
	nodes[1].PreferredDeviceId = -1
	info := []*dto.Info{&dto.Info{1, dto.MemoryInfo{GpuId: 0, MemoryTotal: 1, MemoryFree: 1, GpuMemoryTotal:  1, GpuMemoryFree: 1, DBMemoryFree: 112}}, &dto.Info{1, dto.MemoryInfo{GpuId: 1, MemoryTotal: 1, MemoryFree: 1, GpuMemoryTotal:  1, GpuMemoryFree: 1, DBMemoryFree: 22222222}}}

	lb.update(info)
	if lb.CurrentInsertNodeId != 1 {
		t.Errorf("Wrong node selected. Expected #%d but get #%d", 1, lb.CurrentInsertNodeId)
	}
	if nodes[1].PreferredDeviceId != 0 {
		t.Errorf("Wrong card selected. Expected #%d but get #%d", 1, nodes[1].PreferredDeviceId)
	}
}

func Test_Update_With_Two_Nodes_Two_GPUs_Each_Cause_Changing_Node(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1}
	nodes[2] = NewNode(2, nil, nil)
	nodes[2].GpuIds = []int32{0, 1}


	lb := NewLoadBalancer(nodes)
	lb.CurrentInsertNodeId = -1
	nodes[1].PreferredDeviceId = -1
	info := []*dto.Info{&dto.Info{1, dto.MemoryInfo{GpuId:0, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 10, DBMemoryFree:  1112312312}}, &dto.Info{1, dto.MemoryInfo{ GpuId:1, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 1, DBMemoryFree: 11231}},
		&dto.Info{2, dto.MemoryInfo{GpuId:0, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   11, GpuMemoryFree: 1, DBMemoryFree: 12}}, &dto.Info{2, dto.MemoryInfo{GpuId:1, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal: 14, GpuMemoryFree: 1, DBMemoryFree:112222222}}}

	lb.update(info)
	if lb.CurrentInsertNodeId != 2 {
		t.Errorf("Wrong node selected. Expected #%d but get #%d", 2, lb.CurrentInsertNodeId)
	}

	if nodes[1].PreferredDeviceId != 1 {
		t.Errorf("Wrong card selected. Expected #%d but get #%d", 0, nodes[1].PreferredDeviceId)
	}

	if nodes[2].PreferredDeviceId != 0 {
		t.Errorf("Wrong card selected. Expected #%d but get #%d", 0, nodes[2].PreferredDeviceId)
	}

}


func Test_CalculateNodeRank(t *testing.T) {
	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1}
	nodes[2] = NewNode(2, nil, nil)
	nodes[2].GpuIds = []int32{0, 1}

	lb := NewLoadBalancer(nodes)

	nodes[1].PreferredDeviceId = -1
	nodes[2].PreferredDeviceId = -1
	info := []*dto.Info{&dto.Info{1, dto.MemoryInfo{GpuId:0, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 14, DBMemoryFree:  1123123123}}, &dto.Info{1, dto.MemoryInfo{GpuId:1, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 14, DBMemoryFree: 4545531}},
		&dto.Info{2, dto.MemoryInfo{GpuId:0, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 10, DBMemoryFree:  34534534}}, &dto.Info{2, dto.MemoryInfo{GpuId:1, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 11, DBMemoryFree: 435}}}


	ranks := []uint64{4545531, 435}

	for id, node := range nodes {
		rank := lb.calculateNodeRank(node, info)
		if rank != ranks[id - 1] {
			t.Errorf("Wrong rank calculated. Expected #%d but get #%d", ranks[id], rank)
		}
	}

	if nodes[1].PreferredDeviceId != 1 {
		t.Errorf("Wrong card selected. Expected #%d but get #%d", 0, nodes[1].PreferredDeviceId)
	}

	if nodes[2].PreferredDeviceId != 1 {
		t.Errorf("Wrong card selected. Expected #%d but get #%d", 0, nodes[2].PreferredDeviceId)
	}
}

func Test_RemoveDeadNodes(t *testing.T) {

	nodes := make(map[int32]*Node)
	nodes[1] = NewNode(1, nil, nil)
	nodes[1].GpuIds = []int32{0, 1}
	nodes[2] = NewNode(2, nil, nil)
	nodes[2].GpuIds = []int32{0, 1}

	go NodeManager.Manage()
	NodeManager.nodes = nodes
	lb:=NewLoadBalancer(NodeManager.GetNodes())
	info := []*dto.Info{&dto.Info{1, dto.MemoryInfo{GpuId:0, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 14, DBMemoryFree:  1123123123}}, &dto.Info{1, dto.MemoryInfo{GpuId:1, MemoryTotal:  1, MemoryFree:  1, GpuMemoryTotal:   1, GpuMemoryFree: 14, DBMemoryFree: 4545531}}}

	lb.removeDeadNodes(info)
	if len(lb.nodes) != 1 || lb.nodes[2] != nil {
		t.Errorf("Node #%d should have been removed from Balancer, but wasn't.", 2)
	}
	if NodeManager.GetNodesLen() != 1 || NodeManager.GetNodes()[2]!= nil {
		t.Errorf("Node #%d should have been removed from Manager, but wasn't.", 2)
	}

}
