package dto

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"ddj_Master/dto"
	log "code.google.com/p/log4go"
)

type Request struct {
	Header
	Data []byte
}

func NewRequest(id int64, ttype int32, size int32, data dto.Dto) *Request {
	r := new(Request)
	r.Header = Header{id, ttype, size}
	var err error
	r.Data, err = data.Encode()
	if err != nil {
		log.Error(err)
	}
	return r
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

func (r *Request) Encode() ([]byte, error) {

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
