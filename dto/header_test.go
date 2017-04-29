package dto

import (
	"bytes"
	"testing"
)

func Test_EncodeDecode_Header(t *testing.T) {
	expected := Header{1, 4, 24, 0}
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual Header
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Encode_Header(t *testing.T) {
	expected := []byte{1, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 24, 0, 0, 0, 76, 199, 124, 0, 132, 14, 55, 0, 0, 0, 249, 63, 0, 0, 0, 0, 0, 0, 96, 9, 0, 0, 0, 0}
	actual, err := (&Header{1, 4, 24, 0}).Encode()
	if err != nil {
		t.Error(err)
	}
	if bytes.Equal(actual, expected) {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func Test_Decode_Header(t *testing.T) {
	expected := Header{1, 4, 24, 1}
	buf := []byte{1, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 24, 0, 0, 0, 1, 0, 0, 0, 76, 199, 124, 0, 132, 14, 55, 0, 0, 0, 249, 63, 0, 0, 0, 0, 0, 0, 96, 9, 0, 0, 0, 0}
	var actual Header
	err := actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
