package dto

import (
	"bytes"
)

type Dtos []Dto

func (this Dtos) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	for _, dto := range this {
		buf, err := dto.Encode()
		if err != nil {
			return nil, err
		}
		buffer.Write(buf)
	}
	return buffer.Bytes(), nil
}

func (this Dtos) Size() int {
	size := 0
	for _, dto := range this {
		size += dto.Size()
	}
	return size
}

func (this Dtos) String() string {
	str := ""
	for _, dto := range this {
		str += dto.String() + "; "
	}
	return str
}
