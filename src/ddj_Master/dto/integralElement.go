package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type IntegralElement struct {
	Integral Value
	// left store element
	LeftValue Value
	LeftTime  int64
	// right store element
	RightValue Value
	RightTime  int64
}

type ByLeftTime []*IntegralElement

func (this ByLeftTime) Len() int      { return len(this) }
func (this ByLeftTime) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this ByLeftTime) Less(i, j int) bool {
	return this[i].LeftTime < this[j].LeftTime
}

func (this *IntegralElement) String() string {
	return fmt.Sprintf("Value: %f Left: [Value: %f Time: %d] Right: [Value: %f Time: %d]",
		this.Integral, this.LeftValue, this.LeftTime, this.RightValue, this.RightTime)
}

func (this *IntegralElement) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *IntegralElement) Decode(buf []byte) error {

	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, this)
}

func (this *IntegralElement) Size() int {

	return binary.Size(this)
}
