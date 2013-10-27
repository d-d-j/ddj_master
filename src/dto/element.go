package dto

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func (e *Element) Equal(other *Element) bool {
	left, err := e.GobEncode()
	if err != nil {
		return false
	}
	right, err := e.GobEncode()
	if err != nil {
		return false
	}
	if bytes.Equal(left, right) {
		return true
	}
	return false
}

func (e *Element) String() string {
	return fmt.Sprintf("%d#%d %d [%f]", e.Series, e.Tag, e.Time, e.Value)
}

func (e *Element) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(e.Series)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(e.Tag)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(e.Time)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(e.Value)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e *Element) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&e.Series)
	if err != nil {
		return err
	}
	err = decoder.Decode(&e.Tag)
	if err != nil {
		return err
	}
	err = decoder.Decode(&e.Time)
	if err != nil {
		return err
	}
	return decoder.Decode(&e.Value)
}
