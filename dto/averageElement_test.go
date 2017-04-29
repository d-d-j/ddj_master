package dto

import (
	"bytes"
	"testing"
)

func Test_Encode_Variance(t *testing.T) {
	expected := []byte{2, 0, 0, 0, 0, 0, 128, 63}
	actual, err := (&VarianceElement{2, 1, 0}).Encode()
	if err != nil {
		t.Error(err)
	}
	if bytes.Equal(actual, expected) {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Decode_Variance(t *testing.T) {
	expected := &VarianceElement{2, 1, 0}
	buf := []byte{2, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 0}
	var actual VarianceElement
	err := actual.Decode(buf)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}

	if expected.Size() != 12 {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
