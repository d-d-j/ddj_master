package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//This is structure that will be receive from node for histogram aggregation
type Histogram struct {
	Data []int32
}

func (this Histogram) String() string {
	return fmt.Sprintf("Histogram: %v", this.Data)

}

//This method is used only for tests and to fulfill dto interface. It shouldn't be used
func (this *Histogram) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *Histogram) Decode(buf []byte) error {
	length := len(buf) / 4
	if length < 1 {
		return fmt.Errorf("Empty buffer. ")
	}
	buffer := bytes.NewBuffer(buf)
	newHistogram := make([]int32, length)
	for i := 0; i < length; i++ {
		var value int32
		err := binary.Read(buffer, binary.LittleEndian, &value)
		newHistogram[i] = value
		if err != nil {
			return err
		}
	}
	this.Data = newHistogram
	return nil
}

func (this Histogram) Size() int {

	return len(this.Data) * 4
}
