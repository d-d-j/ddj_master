package tests

import "testing"
import "ddj_Master/dto"


func TestNewElement(t *testing.T) {
	var expected dto.Element
	expected.Series = 1
	expected.Tag = 2
	expected.Time = 0
	expected.Value = 0.33
	var actual = dto.NewElement(1, 2, 0, 0.33)
	if expected.String() != actual.String() {
		t.Error()
	}
}

func TestString(t *testing.T) {
	expected := "1#2 1970-01-01 01:00:00 +0100 CET [0.330000]"
	var actual = dto.NewElement(1, 2, 0, 0.33).String()
	if expected != actual {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}

func TestEncodeDecode(t *testing.T) {
	expected := dto.NewElement(1, 2, 0, 0.33)
	buf, err := expected.Encode()
	if err != nil {
		t.Error(err)
	}
	var actual dto.Element
	err = actual.Decode(buf)
	if err != nil {
		t.Error(err)
	}
	if expected.String() != actual.String() {
		t.Error("Got: ", actual, " when expected ", expected)
	}
}
