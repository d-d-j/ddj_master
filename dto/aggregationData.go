package dto

import (
	"bytes"
	"encoding/binary"
)

//This interface is used for additional data that could be send with some aggregations
type AggregationData interface {
	Encoder
	Size() int
	GetBucketCount() int32
}

//This structure keeps additional data for histogram by value aggregation
type HistogramValueData struct {
	Min         float32
	Max         float32
	BucketCount int32
}

//This structure keeps additional data for histogram by time aggregation
type HistogramTimeData struct {
	Min         int64
	Max         int64
	BucketCount int32
}

//This structure keeps additional data for series aggregation with interpolation
type InterpolatedData int32

//This method do binary encoding and return result as a slice of bytes
func (this InterpolatedData) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this InterpolatedData) GetBucketCount() int32 {
	return int32(this)
}

func (this InterpolatedData) Size() int {
	return 4
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
