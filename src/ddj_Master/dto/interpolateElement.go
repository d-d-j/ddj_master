package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//This structure is revived for series operation
type InterpolateElement struct {
	Data []Value
}

func (this InterpolateElement) String() string {
	return fmt.Sprintf("InterpolateElement: %v", this.Data)

}

func (this *InterpolateElement) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *InterpolateElement) Decode(buf []byte) error {
	length := len(buf) / 4
	if length < 1 {
		return fmt.Errorf("Empty buffer. ")
	}
	buffer := bytes.NewBuffer(buf)
	newInterpolateElement := make([]Value, length)
	for i := 0; i < length; i++ {
		var value Value
		err := binary.Read(buffer, binary.LittleEndian, &value)
		newInterpolateElement[i] = value
		if err != nil {
			return err
		}
	}
	this.Data = newInterpolateElement
	return nil
}

func (this InterpolateElement) Size() int {

	return len(this.Data) * 4
}
