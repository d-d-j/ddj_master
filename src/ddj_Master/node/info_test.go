package node

import (
	"testing"
	"ddj_Master/dto"
)

func Test_MemoryInfo_EncodeDecode(t *testing.T) {

	expected := &dto.MemoryInfo{GpuId:0, MemoryTotal:  1, MemoryFree:  6, GpuMemoryTotal:   5, GpuMemoryFree: 3, DBMemoryFree:  2}
	actual := &dto.MemoryInfo{}
	buf, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred: ", err)
	}
	err = actual.Decode(buf)
	if err != nil {
		t.Error("Error occurred: ", err)
	}
	if actual.String() != expected.String() {
		t.Error("Expected: ", expected, " but got: ", actual)
	}
}


func Test_MemoryInfo_DecodeBinary(t *testing.T) {

	expected := &dto.MemoryInfo{GpuId:0, MemoryTotal:  16248688, MemoryFree:  6844272, GpuMemoryTotal:   2147155968, GpuMemoryFree: 1585799168, DBMemoryFree:  141840}
	actual := &dto.MemoryInfo{}

	buf := []byte{0, 0, 0, 0, 112, 239, 247, 0, 112, 111, 104, 0, 0, 0, 251, 127, 0, 96, 133, 94, 0, 0, 0, 0, 16, 42, 2, 0, 0, 0, 0, 0}
	err := actual.Decode(buf)
	if err != nil {
		t.Error("Error occurred: ", err)
	}
	if actual.String() != expected.String() {
		t.Error("Expected: ", expected, " but got: ", actual)
	}
}
