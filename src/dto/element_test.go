package dto

import "testing"

func TestNewElement(t *testing.T) {
	var expected Element
	expected.series = 1
	expected.tag = 2
	expected.time = 0
	expected.value = 0.33
	var actual = NewElement(1, 2, 0, 0.33)
	if !expected.Equal(actual) {
		t.Error()
	}
}

func TestString(t *testing.T) {
	expected := "1#2 1970-01-01 01:00:00 +0100 CET [0.330000]"
	var actual = NewElement(1, 2, 0, 0.33).String()
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func TestEncodeDecode(t *testing.T) {
	expected := NewElement(1, 2, 0, 0.33)
	buf, err := expected.GobEncode()
	if err != nil {
		t.Error("Problem with encoding")
	}
	var actual Element
	err = actual.GobDecode(buf)
	if err != nil {
		t.Error("Problem with decoding")
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
