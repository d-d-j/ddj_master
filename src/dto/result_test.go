package dto

import "testing"

func TestResultString(t *testing.T) {
	expected := "#1 Code: 2"
	r := Result{1, 2}
	actual := r.String()
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func TestResultEncodeDecode(t *testing.T) {
	expected := &Result{1, 2}
	buf, err := expected.GobEncode()
	if err != nil {
		t.Error("Problem with encoding")
	}
	var actual Result
	err = actual.GobDecode(buf)
	if err != nil {
		t.Error("Problem with decoding")
	}
	if !actual.Equal(expected) {
		t.Error("Got: ", actual.String(), " when expected ", expected.String())
	}
}
