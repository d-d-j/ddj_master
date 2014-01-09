package dto

import "testing"

func Test_EncodeDecode_Histogram(t *testing.T) {
	expected := Histogram{[]int32{0, 1, 1, 5, 7, 2}}
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	t.Log(buf)
	var actual Histogram
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Histogram_Size(t *testing.T) {
	actual := Histogram{[]int32{0, 1, 1, 5, 7, 2}}.Size()
	expected := 24
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Histogram_String(t *testing.T) {
	actual := Histogram{[]int32{0, 1, 1, 5, 7, 2}}.String()
	expected := "Histogram: [0 1 1 5 7 2]"
	if actual != expected {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
