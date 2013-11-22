package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	log "code.google.com/p/log4go"
)

type Result struct {
	Header
	Data []byte
}

func NewResult(id int64, ttype int32, size int32, data []byte) *Result {
	r := new(Result)
	r.Header = Header{id, ttype, size}
	r.Data = data
	return r
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

func (r *Result) Encode() ([]byte, error) {

	var (
		buf       []byte
		headerBuf []byte
		err       error
	)

	// Encode header nad data
	if r.Data != nil {
		buf, err = r.Data.Encode()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	headerBuf, err = r.Header.Encode()
	if err != nil {
		log.Error(err)
		continue
	}

	// Merge header and data to one []byte buffer
	// TODO: CHANGE 100 to real data size
	complete := make([]byte, 100)
	copy(complete, headerBuf)
	copy(complete[len(headerBuf):], buf)

	return complete, err
}

func (r *Result) Decode(buf []byte) error {

	return r.Header.Decode(buf)

}
