package dto

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Query struct {
	Id   int64
	Code int32
	Load string
}

func (q *Query) String() string {
	return fmt.Sprintf("#%d Code: %d [%s]", q.Id, q.Code, q.Load)
}

func (q *Query) Equal(other *Query) bool {
	return q.Id == other.Id && q.Code == other.Code && q.Load == other.Load
}

func (q *Query) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	err := encoder.Encode(q.Id)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(q.Code)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(q.Load)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (q *Query) GobDecode(buf []byte) error {
	b := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(b)
	err := decoder.Decode(&q.Id)
	if err != nil {
		return err
	}
	err = decoder.Decode(&q.Code)
	if err != nil {
		return err
	}
	return decoder.Decode(&q.Load)
}
