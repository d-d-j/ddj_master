package dto

import (
	"fmt"
	"github.com/ugorji/go/codec"
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
	var (
		buf []byte
		mh  codec.MsgpackHandle
	)
	enc := codec.NewEncoderBytes(&buf, &mh)
	err := enc.Encode(e)
	return buf, err
}

func (e *Element) Decode(buf []byte) error {
	var mh codec.MsgpackHandle
	dec := codec.NewDecoderBytes(buf, &mh)
	err := dec.Decode(e)
	return err
}
