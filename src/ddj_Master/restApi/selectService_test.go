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
	expected := dto.Query{MetricsCount: 3, Metrics: []int32{1, 2, 3}, TagsCount: 4, Tags: []int32{4, 5, 6, 7}, TimeSpansCount: 3, TimeSpans: []int64{10, 20, 30, 40, 50, 60}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_All_Metrics_With_Specified_Tags(t *testing.T) {
	metrics := "all"
	tags := "4,5,6,7"
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{MetricsCount: 0, Metrics: []int32{}, TagsCount: 4, Tags: []int32{4, 5, 6, 7}, TimeSpansCount: 3, TimeSpans: []int64{10, 20, 30, 40, 50, 60}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_All_Tags_Wiht_Specified_Metrics(t *testing.T) {
	metrics := "1,2,3"
	tags := "all"
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{ MetricsCount: 3, Metrics:  []int32{1, 2, 3}, TagsCount: 0, Tags: []int32{}, TimeSpansCount: 3, TimeSpans: []int64{10, 20, 30, 40, 50, 60}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_All_Tags_And_All_Metrics(t *testing.T) {
	metrics := "all"
	tags := "all"
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{MetricsCount: 0, Metrics:  []int32{}, TagsCount: 0, Tags: []int32{}, TimeSpansCount: 3, TimeSpans: []int64{10, 20, 30, 40, 50, 60}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_Invalid_Arguments(t *testing.T) {
	metrics := "invalid"
	tags := "invalid"
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err == nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}
