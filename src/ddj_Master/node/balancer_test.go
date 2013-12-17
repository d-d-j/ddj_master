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

func Test_Update_With_Nil_Entry(t *testing.T) {
	actual := NewLoadBalancer(0, nil)
	expected := NewLoadBalancer(0, nil)
	actual.CurrentInsertNodeId = 1
	actual.CurrentInsertGpuId = 2
	actual.update(nil)
	AreEqual(expected, actual, t)
}
