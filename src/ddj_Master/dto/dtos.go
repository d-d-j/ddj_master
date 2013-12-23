package dto

import (
	"bytes"
)

//From: http://golang.org/pkg/sort/#example_
type ByTime []Dto

func (this ByTime) Len() int      { return len(this) }
func (this ByTime) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this ByTime) Less(i, j int) bool {
	a, okA := (this[i]).(*Element)
	b, okB := (this[j]).(*Element)
	if okA && okB {
		return a.Time < b.Time
	} else {
		panic("Cannot use this method on non Element")
	}
}

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
