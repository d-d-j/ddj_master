package dto

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Query struct {
	id   int64
	code int32
	load []byte
}

func (q *Query) String() string {
	return fmt.Sprintf("#%d Code: %d [%X]", q.id, q.code, q.load)
}

func (q *Query) Equal(other *Query) bool {
	return q.id == other.id && q.code == other.code && bytes.Equal(q.load, other.load)
}

func (q *Query) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	err := encoder.Encode(q.id)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(q.code)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(q.load)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (q *Query) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)
	err := decoder.Decode(&q.id)
	if err != nil {
		return err
	}
	err = decoder.Decode(&q.code)
	if err != nil {
		return err
	}
	return decoder.Decode(&q.load)
}
