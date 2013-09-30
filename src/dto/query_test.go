package dto

import "testing"

func TestQueryString(t *testing.T) {
	expected := "#1 Code: 2 [0000000000]"
	load := []byte{0, 0, 0, 0, 0}
	q := Query{1, 2, load}
	actual := q.String()
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func TestQueryEncodeDecode(t *testing.T) {
	load := []byte{0, 0, 0, 0, 0}
	expected := &Query{1, 2, load}
	buf, err := expected.GobEncode()
	if err != nil {
		t.Error("Problem with encoding")
	}
	var actual Query
	err = actual.GobDecode(buf)
	if err != nil {
		t.Error("Problem with decoding")
	}
	if !actual.Equal(expected) {
		t.Error("Got: ", actual.String(), " when expected ", expected.String())
	}
}
