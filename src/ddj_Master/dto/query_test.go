package dto

import "testing"

func Test_EncodeDecode_Query(t *testing.T) {

	t.Skip("Method not implemented")

	expected := Query{1, []int32{1}, 2, []int32{0, 1}, 4, []int64{5, 7, 11, 21}, 0}
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual Query
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Query_Size(t *testing.T) {
	actual := (&Query{1, []int32{1}, 2, []int32{0, 1}, 4, []int64{5, 7, 11, 21}, 0}).Size()
	expected := (7*32 + 4*64) / 8
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Query_Size_For_Empty_Query(t *testing.T) {
	actual := (&Query{}).Size()
	expected := (4 * 32) / 8
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
