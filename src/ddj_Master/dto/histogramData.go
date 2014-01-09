package dto

import (
	"bytes"
	"encoding/binary"
)

type AggregationData interface {
	Encoder
	Size() int
	GetBucketCount() int32
}

type HistogramValueData struct {
	Min         float32
	Max         float32
	BucketCount int32
}

type HistogramTimeData struct {
	Min         int64
	Max         int64
	BucketCount int32
}

func (this HistogramTimeData) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this HistogramTimeData) GetBucketCount() int32 {
	return this.BucketCount
}

func (this HistogramTimeData) Size() int {
	return binary.Size(this)
}

func (this HistogramValueData) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this HistogramValueData) GetBucketCount() int32 {
	return this.BucketCount
}

func (this HistogramValueData) Size() int {
	return binary.Size(this)
}
