package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type VarianceElement struct {
	Count int32
	Mean  Value
	M2    Value
}

func (this *VarianceElement) String() string {
	return fmt.Sprintf("Mean: %f, Count: %d", this.Mean, this.Count)

}

func (this *VarianceElement) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *VarianceElement) Decode(buf []byte) error {

	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, this)
}

func (this *VarianceElement) Size() int {

	return binary.Size(this)
}
