package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Result struct {
	Header
	Data []Dto
}

func (r *Result) String() string {
	return fmt.Sprintf("#%d Code: %d", r.Id, r.Code)
}

func (r *Result) EncodeHeader() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, r.Header)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Result) DecodeHeader(buf []byte) error {

	return r.Header.Decode(buf)

}
