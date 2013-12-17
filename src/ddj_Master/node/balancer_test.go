package node

import (
	"ddj_Master/common"
	"testing"
)

func TestNewLoadBalancer(t *testing.T) {
	expected := new(LoadBalancer)
	expected.CurrentInsertNodeId = common.CONST_UNINITIALIZED
	expected.CurrentInsertGpuId = common.CONST_UNINITIALIZED
	actual := NewLoadBalancer(0)
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
