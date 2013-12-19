package dto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Query struct {
	MetricsCount    int32
	Metrics         []int32
	TagsCount       int32
	Tags            []int32
	TimeSpansCount  int32
	TimeSpans       []int64
	AggregationType int32
}

func (q *Query) String() string {
	output := ""
	output += fmt.Sprintf("Metrics [%d]: ", q.MetricsCount)
	var i int32
	for i = 0; i < q.MetricsCount; i++ {
		output += fmt.Sprintf("%d, ", q.Metrics[i])
	}
	output += fmt.Sprintf("\tTags [%d]: ", q.TagsCount)
	for i = 0; i < q.TagsCount; i++ {
		output += fmt.Sprintf("%d", q.Tags[i])
	}
	output += fmt.Sprintf("\tTimes [%d]: ", q.TimeSpansCount)
	for i = 0; i < q.TagsCount*2; i += 2 {
		from := time.Unix(q.TimeSpans[i], 0)
		to := time.Unix(q.TimeSpans[i+1], 0)
		output += fmt.Sprintf("%s-%s ", from, to)
	}
	return output

}

func (q *Query) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, q.MetricsCount)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (q *Query) Decode(buf []byte) error {

	buffer := bytes.NewBuffer(buf)
	return binary.Read(buffer, binary.LittleEndian, q)
}

func (q *Query) Size() int {

	return binary.Size(q)
}
