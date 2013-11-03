package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Element struct {
	Series, Tag int32
	Time        int64
	Value       float32
}

func NewElement(Series int32, Tag int32, Time int64, Value float32) *Element {
	e := new(Element)
	e.Series = Series
	e.Tag = Tag
	e.Time = Time
	e.Value = Value
	return e
}

func (e *Element) String() string {
	t := time.Unix(e.Time, 0)
	return fmt.Sprintf("%d#%d %s [%f]", e.Series, e.Tag, t, e.Value)

}

func (e *Element) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, e)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *Element) Decode(buf []byte) error {

	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, e)
}

func (e *Element) Size() int {

	return binary.Size(e)
}
