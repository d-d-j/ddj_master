package dto

import (
	"testing"
)

func Test_Dtos_Encode(t *testing.T) {
	expected := Dtos{NewElement(1, 2, 3, 0.33), NewElement(4, 5, 6, 0.66), NewElement(7, 8, 9, 0.99)}
	buffer, err := expected.Encode()
	if err != nil {
		t.Error("Error occurred: ", err)
	}
	length := len(expected) / expected[0].Size()
	actual := make([]Dto, length)

	for i := 0; i < length; i++ {
		var e Element
		err = e.Decode(buffer[i*e.Size():])
		if err != nil {
			t.Error(err)
			continue
		}
		actual[i] = &e
	}

	for i := 0; i < length; i++ {
		if expected[i].String() != actual[i].String() {
			t.Error("Expected: ", expected[i], " but got: ", actual[i])
		}
	}

}
