package dto

import "testing"

func Test_EncodeDecode_Query(t *testing.T) {
	expected := Query{1, []int32{1}, 2, []int32{0, 1}, 3, []int64{5, 7, 11, 17}, 0}
	t.Log("Input: ", expected)
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	t.Log("Encoded: ", buf)
	var actual Query
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
