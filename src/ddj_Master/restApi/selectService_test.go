package restApi

import (
	"ddj_Master/dto"
	"testing"
)

func Test_StartService(t *testing.T) {
	var server = Server{Port: ":6666"}
	chanReq := server.StartApi()
	if chanReq == nil {
		t.Error("Nil channel")
	}
}

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
	metrics := ALL
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
	tags := ALL
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{MetricsCount: 3, Metrics: []int32{1, 2, 3}, TagsCount: 0, Tags: []int32{}, TimeSpansCount: 3, TimeSpans: []int64{10, 20, 30, 40, 50, 60}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_All_Tags_And_All_Metrics(t *testing.T) {
	metrics := ALL
	tags := ALL
	times := "10-20,30-40,50-60"
	aggr := "none"
	expected := dto.Query{MetricsCount: 0, Metrics: []int32{}, TagsCount: 0, Tags: []int32{}, TimeSpansCount: 3, TimeSpans: []int64{10, 20, 30, 40, 50, 60}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_All_Tags_And_All_Metrics_And_All_Times(t *testing.T) {
	metrics := ALL
	tags := ALL
	times := ALL
	aggr := "none"
	expected := dto.Query{MetricsCount: 0, Metrics: []int32{}, TagsCount: 0, Tags: []int32{}, TimeSpansCount: 0, TimeSpans: []int64{}, AggregationType: 0}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err != nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}

func Test_Preapare_Query_For_All_Tags_And_All_Metrics_And_All_Times_For_Integral_Should_Return_Error(t *testing.T) {
	metrics := ALL
	tags := ALL
	times := ALL
	aggr := "int"
	_, err := prepareQuery(metrics, tags, times, aggr)
	if err == nil {
		t.Error("Error was expected")
	}
}

func Test_Preapare_Query_For_All_Tags_And_All_Metrics_And_Invalid_Time_Spans_Shoudl_Return_Error(t *testing.T) {
	metrics := ALL
	tags := ALL
	times := "10-5"
	aggr := "none"
	_, err := prepareQuery(metrics, tags, times, aggr)
	if err == nil {
		t.Error("Error was expected")
	}
}

func Test_prepareTagsOrMetrics_Should_Return_Sorted_And_Unique_Arraya_of_Numbers(t *testing.T) {
	input := "1,1,2,5,7,2,5,1"
	expected := []int32{1, 2, 5, 7}
	actual, _ := prepareTagsOrMetrics(input)
	if len(expected) != len(actual) {
		t.Error("Element count doesn't match.")
	}

	for index, e := range expected {
		if e != actual[index] {
			t.Error("Expected ", e, " but got ", actual[index])
		}
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

func Test_Preapare_Query_For_Invalid_Arguments_No_Times(t *testing.T) {
	metrics := ALL
	tags := ALL
	times := ""
	aggr := "none"
	expected := dto.Query{}
	actual, err := prepareQuery(metrics, tags, times, aggr)
	if expected.String() != actual.String() || err == nil {
		t.Error("Got: ", actual, "when expecting: ", expected)
		t.Error("error: ", err)
	}
}
