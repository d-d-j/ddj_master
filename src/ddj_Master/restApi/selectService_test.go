package restApi

import (
	"ddj_Master/dto"
	"testing"
)

func Test_Preapare_Query_For_Valid_Arguments(t *testing.T) {
	metrics := "1,2,3"
	tags := "4,5,6,7"
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{3, []int32{1, 2, 3}, 4, []int32{4, 5, 6, 7}, 3, []int64{10, 20, 30, 40, 50, 60}, 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}
