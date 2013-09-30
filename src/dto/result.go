package dto

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Result struct {
	id   int64
	code int32
}

func (r *Result) Equal(other *Result) bool {
	if r.id == other.id && r.code == other.code {
		return true
	}
	return false
}

func (r *Result) String() string {
	return fmt.Sprintf("#%d Code: %d", r.id, r.code)
}

func (r *Result) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	err := encoder.Encode(r.id)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(r.code)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (r *Result) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)
	err := decoder.Decode(&r.id)
	if err != nil {
		return err
	}
	return decoder.Decode(&r.code)

}
