package dto

import "testing"

func TestNewStatus(t *testing.T) {
	expected := Status{1024, 512, 54}
	var actual = NewStatus(1024, 512, 54)
	if !expected.Equal(actual) {
		t.Error()
	}
}

func TestStatusString(t *testing.T) {
	expected := "Ram: 512/1024 Temp: 54℃"
	var actual = NewStatus(1024, 512, 54).String()
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func TestStatusEncodeDecode(t *testing.T) {
	expected := NewStatus(1024, 512, 54)
	buf, err := expected.GobEncode()
	if err != nil {
		t.Error("Problem with encoding")
	}
	var actual Status
	err = actual.GobDecode(buf)
	if err != nil {
		t.Error("Problem with decoding")
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
