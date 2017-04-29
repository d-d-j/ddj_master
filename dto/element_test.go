package dto

import "testing"

func TestEncodeDecode(t *testing.T) {
	expected := NewElement(1, 2, 0, 0.33)
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual Element
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Element_Size(t *testing.T) {
	actual := NewElement(1, 2, 0, 0.33).Size()
	expected := 24
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
