package dto

import (
	"fmt"
)

type EmptyElement struct{}

func (this *EmptyElement) String() string {
	return fmt.Sprintf("Empty DTO used for INFO ")
}

func (this *EmptyElement) Encode() ([]byte, error) {
	return make([]byte, 0, 0), nil
}

func (this *EmptyElement) Decode(buf []byte) error {
	if len(buf) == 0 {
		this = new(EmptyElement)
	}
	return fmt.Errorf("Expected empty slice (length = 0) but get  %d", len(buf))
}

func (this *EmptyElement) Size() int {

	return 0
}
