package node

import (
	"testing"
	"ddj_Master/dto"
)

func Test_MemoryInfo_EncodeDecode(t *testing.T) {

	expected := &dto.MemoryInfo{1, 5, 4, 3, 2}
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
