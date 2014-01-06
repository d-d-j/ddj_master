package dto

import "testing"

func Test_Integral_EncodeDecode(t *testing.T) {
	expected := IntegralElement{0.5, 0.1, 1, 0.2, 2}
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual IntegralElement
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
