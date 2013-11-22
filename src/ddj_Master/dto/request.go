package dto

import (
	"fmt"
	"bytes"
	"encoding/binary"
)

type Request struct {
	Header
	Data []byte
}

func (r *Request) String() string {
	return fmt.Sprintf("#%d Code: %d", r.Id, r.Code)
}

func (r *Request) EncodeHeader() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, r.Header)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Request) DecodeHeader(buf []byte) error {

	return r.Header.Decode(buf)

}
