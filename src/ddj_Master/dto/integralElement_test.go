package dto

import "testing"

func Test_Integral_EncodeDecode(t *testing.T) {
	expected := &IntegralElement{0.5, 0.1, 1, 0.2, 2}
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

func Test_Integral_Decode(t *testing.T) {
	expected := &IntegralElement{0.1, -0.5, 5, 0.5, 10}
	buf := []byte{205, 204, 204, 61, 0, 0, 0, 191, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 255, 127, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0}

	var actual IntegralElement
	err := actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Integral_Size(t *testing.T) {
	expected := 32
	actual := (IntegralElement{}).Size()

	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
