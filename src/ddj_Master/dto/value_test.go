package dto

import "testing"

func Test_EncodeDecode_Value(t *testing.T) {
	expected := Value(0.99)
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual Value
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Value_Size(t *testing.T) {
	actual := Value(0).Size()
	expected := 4
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
