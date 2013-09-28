package dto

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Element struct {
	series, tag int32
	time        int64
	value       float32
}

func NewElement(series int32, tag int32, time int64, value float32) *Element {
	e := new(Element)
	e.series = series
	e.tag = tag
	e.time = time
	e.value = value
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
	t := time.Unix(e.time, 0)
	return fmt.Sprintf("%d#%d %s [%f]", e.series, e.tag, t, e.value)
}

func (e *Element) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(e.series)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(e.tag)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(e.time)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(e.value)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e *Element) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&e.series)
	if err != nil {
		return err
	}
	err = decoder.Decode(&e.tag)
	if err != nil {
		return err
	}
	err = decoder.Decode(&e.time)
	if err != nil {
		return err
	}
	return decoder.Decode(&e.value)
}
