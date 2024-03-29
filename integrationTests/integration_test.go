//Here are integration tests for project. Following test load some data and then do multiple selects to check
//if every integration is working
package integrationTests

import (
	"github.com/d-d-j/ddj_master/dto"
	"encoding/json"
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

//This is response element for simple select without aggregation. Returned data are simple Elements
type RestResponse struct {
	Error  string
	TaskId int64
	Data   []dto.Element
}

//This is response element for simple aggregation (not series aggregation) that return array of values
type RestForAggregation struct {
	Error  string
	TaskId int64
	Data   []dto.Value
}

//This is response elment for Histogram aggregation
type RestForHistogram struct {
	Error  string
	TaskId int64
	Data   []dto.Histogram
}

//This is response element for series aggregation (with interpolation)
type RestForSeries struct {
	Error  string
	TaskId int64
	Data   []dto.InterpolateElement
}

//Test set up constants
const (
	NUMBER_OF_TAGS_PER_METRICS int     = 4
	NUMBER_OF_METRICS          int     = 4
	INSERTED_DATA              int     = 10000
	HOST                       string  = "http://localhost:8888/data"
	eps                        float64 = 0.01
)

var (
	insert_completed bool
	data             []string
	expected         []dto.Element
	random           *rand.Rand
)

//This method generate data and load data to DB. It will populate expected slice with data uploaded to nodes.
func SetUp(b *testing.B) {

	if insert_completed != true {
		random = rand.New(rand.NewSource(99))
		data = make([]string, INSERTED_DATA)
		expected = make([]dto.Element, INSERTED_DATA)
		for i := 0; i < INSERTED_DATA; i++ {
			value := dto.Value(math.Log(float64(i + 1)))
			e := *dto.NewElement(int32(i%NUMBER_OF_TAGS_PER_METRICS), int32(i%NUMBER_OF_METRICS), int64(i), value)
			data[i] = fmt.Sprintf("{\"tag\":%d, \"metric\":%d, \"time\":%d, \"value\":%f}", e.Metric, e.Tag, e.Time, value)
			expected[i] = e
		}

		client := &http.Client{}
		for i := 0; i < INSERTED_DATA; i++ {

			req, err := http.NewRequest("POST", HOST, strings.NewReader(data[i]))
			if err != nil {
				b.Log("Error occurred: ", err)
				b.FailNow()
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Connection", "Keep-Alive")
			response, err := client.Do(req)

			if err != nil {
				b.Error("Error: ", err)
				b.FailNow()
			}

			if response.StatusCode != 202 {
				b.Log("Status Code:", response.StatusCode, err)
				b.FailNow()
			}
			response.Body.Close()
		}

		insert_completed = true
		b.Log("Inserted ", INSERTED_DATA, " elements")
		Flush(b)
	}
}

func Benchmark_Select_First_Value(b *testing.B) {

	SetUp(b)
	queryString := "/metric/0/tag/0/time/0-1/aggregation/none"
	response := Select(queryString, b)

	Assert(response, expected, queryString, b)
}

func Benchmark_Select_Values_From_All_Tags_And_Metrics(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)

	from := expected[f].Time
	to := expected[t].Time

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, nil, nil)
		}).Results()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	queryString := fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/none", from, to)
	response := Select(queryString, b)

	expectedElements := make([]dto.Element, 0, len(exp))
	for _, e := range exp {
		expectedElements = append(expectedElements, e.(dto.Element))
	}

	Assert(response, expectedElements, queryString, b)
}

func Benchmark_Select_All(b *testing.B) {

	SetUp(b)

	queryString := fmt.Sprintf("/metric/all/tag/all/time/all/aggregation/none")
	response := Select(queryString, b)

	Assert(response, expected, queryString, b)
}

func Flush(b *testing.B) {
	client := &http.Client{}
	flushUrl := fmt.Sprintf("%s/%s", HOST, "flush")
	req, err := http.NewRequest("POST", flushUrl, strings.NewReader(""))
	if err != nil {
		b.Log("Error occurred: ", err)
		b.FailNow()
	}
	response, err := client.Do(req)

	if response.StatusCode != 202 {
		b.Log("Status Code:", response.StatusCode, err)
		b.FailNow()
	}
	response.Body.Close()
	b.Log("Flushed buffers")
}

func Select(query string, b *testing.B) RestResponse {
	b.ResetTimer()

	selectUrl := fmt.Sprintf("%s/%s", HOST, query)
	req, err := http.Get(selectUrl)
	defer req.Body.Close()
	if err != nil {
		b.Log("Error occurred: ", err)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		b.Log("Error occurred: ", err)
		b.FailNow()
	}

	b.StopTimer()

	response := RestResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		b.Error("Error occurred: ", err)
		b.FailNow()
	}

	return response
}

func SelectAggr(query string, b *testing.B) RestForAggregation {
	b.ResetTimer()

	selectUrl := fmt.Sprintf("%s/%s", HOST, query)
	req, err := http.Get(selectUrl)
	defer req.Body.Close()
	if err != nil {
		b.Log("Error occurred: ", err)
		b.Log(query)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		b.Log("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}

	b.StopTimer()

	response := RestForAggregation{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		b.Error("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}

	return response
}

func SelectHistogram(query string, b *testing.B) RestForHistogram {
	b.ResetTimer()

	selectUrl := fmt.Sprintf("%s/%s", HOST, query)
	req, err := http.Get(selectUrl)
	defer req.Body.Close()
	if err != nil {
		b.Log("Error occurred: ", err)
		b.Log(query)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		b.Log("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}

	b.StopTimer()

	response := RestForHistogram{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		b.Error("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}

	return response
}

func SelectSeries(query string, b *testing.B) RestForSeries {
	b.ResetTimer()

	selectUrl := fmt.Sprintf("%s/%s", HOST, query)
	req, err := http.Get(selectUrl)
	if err != nil {
		b.Log("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		b.Log("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}

	b.StopTimer()

	response := RestForSeries{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		b.Error("Error occurred: ", err)
		b.Log(query)
		b.FailNow()
	}

	return response
}

func Benchmark_Select_Min(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)

	from := expected[f].Time
	to := expected[t].Time
	metrics := make([]int32, 0, NUMBER_OF_METRICS)
	for metric := 0; metric <= NUMBER_OF_METRICS; metric++ {
		tags := make([]int32, 0, NUMBER_OF_TAGS_PER_METRICS)
		for tag := 0; tag <= NUMBER_OF_TAGS_PER_METRICS; tag++ {
			Select_Min(b, from, to, tags, metrics)
			tags = append(tags, int32(tag))
		}
		metrics = append(metrics, int32(metric))
	}
}

func Benchmark_Select_Max(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)

	from := expected[f].Time
	to := expected[t].Time
	metrics := make([]int32, 0, NUMBER_OF_METRICS)
	for metric := 0; metric <= NUMBER_OF_METRICS; metric++ {
		tags := make([]int32, 0, NUMBER_OF_TAGS_PER_METRICS)
		for tag := 0; tag <= NUMBER_OF_TAGS_PER_METRICS; tag++ {
			Select_Max(b, from, to, tags, metrics)
			tags = append(tags, int32(tag))
		}
		metrics = append(metrics, int32(metric))
	}
}

func Benchmark_Select_Sum(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)

	from := expected[f].Time
	to := expected[t].Time
	metrics := make([]int32, 0, NUMBER_OF_METRICS)
	for metric := 0; metric <= NUMBER_OF_METRICS; metric++ {
		tags := make([]int32, 0, NUMBER_OF_TAGS_PER_METRICS)
		for tag := 0; tag <= NUMBER_OF_TAGS_PER_METRICS; tag++ {
			Select_Sum(b, from, to, tags, metrics)
			tags = append(tags, int32(tag))
		}
		metrics = append(metrics, int32(metric))
	}
}

func Benchmark_Select_Avg(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)

	from := expected[f].Time
	to := expected[t].Time
	metrics := make([]int32, 0, NUMBER_OF_METRICS)
	for metric := 0; metric <= NUMBER_OF_METRICS; metric++ {
		tags := make([]int32, 0, NUMBER_OF_TAGS_PER_METRICS)
		for tag := 0; tag <= NUMBER_OF_TAGS_PER_METRICS; tag++ {
			Select_Avg(b, from, to, tags, metrics)
			tags = append(tags, int32(tag))
		}
		metrics = append(metrics, int32(metric))
	}
}

func Benchmark_Select_Hist(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)

	from := expected[f].Time
	to := expected[t].Time
	metrics := make([]int32, 0, NUMBER_OF_METRICS)
	for metric := 0; metric <= NUMBER_OF_METRICS; metric++ {
		tags := make([]int32, 0, NUMBER_OF_TAGS_PER_METRICS)
		for tag := 0; tag <= NUMBER_OF_TAGS_PER_METRICS; tag++ {
			Select_Histogram_With_One_Bucket(b, from, to, tags, metrics)
			tags = append(tags, int32(tag))
		}
		metrics = append(metrics, int32(metric))
	}
}

func Benchmark_Select_Series_Sum(b *testing.B) {

	SetUp(b)

	t := random.Intn(INSERTED_DATA) + 1
	f := random.Intn(t - 1)
	from := expected[f].Time
	to := expected[t].Time
	tag := expected[f].Tag
	metric := expected[f].Metric

	Select_Series_Sum_With_One_Sample(b, from, to, []int32{tag}, []int32{metric})
}

func ElementInRange(element dto.Element, from, to int64, tags []int32, metrics []int32) (bool, error) {
	var inMetric, inTag, inTime bool
	inMetric = (len(metrics) == 0)
	for _, metric := range metrics {
		if element.Metric == metric {
			inMetric = true
			break
		}
	}
	inTag = len(tags) == 0
	for _, tag := range tags {
		if element.Tag == tag {
			inTag = true
			break
		}
	}
	inTime = element.Time <= to && element.Time >= from

	return inTime && inTag && inMetric, nil
}

func prepareTagsOrMetrics(input []int32) string {
	var str string
	str = ""
	for _, x := range input {
		str += fmt.Sprintf("%d,", x)
	}
	if len(input) == 0 {
		str = "all"
	} else {
		str = strings.TrimSuffix(str, ",")
	}
	return str
}

func Select_Histogram_With_One_Bucket(b *testing.B, from, to int64, tags []int32, metrics []int32) {

	metricsStr := prepareTagsOrMetrics(metrics)
	tagsStr := prepareTagsOrMetrics(tags)
	queryString := fmt.Sprintf("/metric/%s/tag/%s/time/%d-%d/aggregation/histogramByTime/from/%d/to/%d/buckets/1",
		metricsStr, tagsStr, from, to, from, to+1)

	response := SelectHistogram(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, tags, metrics)
		}).Count()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	if len(response.Data[0].Data) != 1 || int(response.Data[0].Data[0]) != exp {
		b.Error("Got ", response.Data, " when expected ", exp)
		b.Log(queryString)
	}
}

func Select_Series_Sum_With_One_Sample(b *testing.B, from, to int64, tags []int32, metrics []int32) {

	metricsStr := prepareTagsOrMetrics(metrics)
	tagsStr := prepareTagsOrMetrics(tags)
	queryString := fmt.Sprintf("/metric/%s/tag/%s/time/from/%d/to/%d/aggregation/series/sum/samples/1",
		metricsStr, tagsStr, from, to)

	response := SelectSeries(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.Log(queryString)
		b.FailNow()
	}

	first, found, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, tags, metrics)
		}).Select(Value).First()

	if err != nil || found == false {
		b.Error("Error: ", err)
		b.FailNow()
	}

	exp := first.(float64)

	if len(response.Data[0].Data) != 1 || math.Abs(float64(response.Data[0].Data[0])-exp) > eps {
		b.Error("Got ", response.Data, " when expected ", exp)
		b.Log(queryString)
	}
}

func Select_Min(b *testing.B, from, to int64, tags []int32, metrics []int32) {

	metricsStr := prepareTagsOrMetrics(metrics)
	tagsStr := prepareTagsOrMetrics(tags)
	queryString := fmt.Sprintf("/metric/%s/tag/%s/time/%d-%d/aggregation/min",
		metricsStr, tagsStr, from, to)

	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, tags, metrics)
		}).Select(Value).MinFloat64()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	if math.Abs(float64(response.Data[0])-exp) > eps {
		b.Error("Got ", response.Data, " when expected ", exp)
		b.Log(queryString)
	}
}

func Select_Max(b *testing.B, from, to int64, tags []int32, metrics []int32) {

	metricsStr := prepareTagsOrMetrics(metrics)
	tagsStr := prepareTagsOrMetrics(tags)
	queryString := fmt.Sprintf("/metric/%s/tag/%s/time/%d-%d/aggregation/max",
		metricsStr, tagsStr, from, to)

	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, tags, metrics)
		}).Select(Value).MaxFloat64()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	if math.Abs(float64(response.Data[0])-exp) > eps {
		b.Error("Got ", response.Data, " when expected ", exp)
		b.Log(queryString)
	}
}

func Select_Sum(b *testing.B, from, to int64, tags []int32, metrics []int32) {

	metricsStr := prepareTagsOrMetrics(metrics)
	tagsStr := prepareTagsOrMetrics(tags)
	queryString := fmt.Sprintf("/metric/%s/tag/%s/time/%d-%d/aggregation/sum",
		metricsStr, tagsStr, from, to)

	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, tags, metrics)
		}).Select(Value).Sum()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	if math.Abs(float64(response.Data[0])-exp) > eps {
		b.Error("Got ", response.Data, " when expected ", exp)
		b.Log(queryString)
	}
}

func Select_Avg(b *testing.B, from, to int64, tags []int32, metrics []int32) {

	metricsStr := prepareTagsOrMetrics(metrics)
	tagsStr := prepareTagsOrMetrics(tags)
	queryString := fmt.Sprintf("/metric/%s/tag/%s/time/%d-%d/aggregation/avg",
		metricsStr, tagsStr, from, to)

	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, tags, metrics)
		}).Select(Value).Average()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	if math.Abs(float64(response.Data[0])-exp) > eps {
		b.Error("Got ", response.Data, " when expected ", exp)
		b.Log(queryString)
	}
}

func Benchmark_Select_Std(b *testing.B) {

	SetUp(b)
	from := expected[0].Time
	to := expected[b.N%INSERTED_DATA].Time + 1
	queryString := fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/std", from, to)
	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	mean, err := From(expected).Where(
		func(element T) (bool, error) {
			elem := element.(dto.Element).Time
			return elem <= to && elem >= from, nil
		}).Select(Value).Average()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	var μ float64
	for i := from; i <= to; i++ {
		v := float64(expected[i].Value)
		μ += (mean - v) * (mean - v)
	}
	σ := dto.Value(math.Sqrt(μ / float64(to-from)))

	if math.Abs(float64(response.Data[0]-σ)) > eps {
		b.Error("Got ", response.Data, " when expected ", σ)
		b.Log(queryString)
	}
}

func Benchmark_Select_Var(b *testing.B) {

	SetUp(b)
	from := expected[0].Time
	to := expected[b.N%INSERTED_DATA].Time + 1
	queryString := fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/var", from, to)
	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	mean, err := From(expected).Where(
		func(element T) (bool, error) {
			elem := element.(dto.Element).Time
			return elem <= to && elem >= from, nil
		}).Select(Value).Average()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	var μ float64
	for i := from; i <= to; i++ {
		v := float64(expected[i].Value)
		μ += (mean - v) * (mean - v)
	}
	σ := dto.Value(μ / float64(to-from))

	if math.Abs(float64(response.Data[0]-σ)) > eps {
		b.Error("Got ", response.Data, " when expected ", σ)
		b.Log(queryString)
	}
}

func Benchmark_Select_Int(b *testing.B) {

	SetUp(b)
	from := expected[0].Time
	to := expected[b.N%INSERTED_DATA].Time
	tag := 0
	metric := 0
	queryString := fmt.Sprintf("/metric/%d/tag/%d/time/%d-%d/aggregation/int", metric, tag, from, to)
	response := SelectAggr(queryString, b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	exp, err := From(expected).Where(
		func(e T) (bool, error) {
			element := e.(dto.Element)
			return ElementInRange(element, from, to, []int32{0}, []int32{0})
		}).Results()

	if err != nil {
		b.Error("Error: ", err)
		b.FailNow()
	}

	integral := dto.Value(0.0)
	for i := 1; i < len(exp); i++ {
		integral += (exp[i].(dto.Element).Value + exp[i-1].(dto.Element).Value) * dto.Value(exp[i].(dto.Element).Time-exp[i-1].(dto.Element).Time) / 2
	}

	if math.Abs(float64(response.Data[0]-integral)) > eps {
		b.Error("Got ", response.Data, " when expected ", integral)
		b.Log(queryString)
	}
}

func Assert(response RestResponse, expectedValues []dto.Element, query string, b *testing.B) {
	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.Log(query)
		b.FailNow()
	}
	if len(response.Data) > len(expectedValues) {
		b.Error("Too many values returned. Expected less or equal than ", len(expectedValues), " but got ", len(response.Data))
		b.Log("Did you forget to restart server?")
		b.Log(query)
		b.FailNow()
	}
	ok := true
	for index, element := range response.Data {
		if expectedValues[index].String() != element.String() {
			b.Error("Expected ", expectedValues[index], " but got ", element)
			ok = false
		}
	}
	if !ok {
		b.Log(query)
	}
}

func Value(element T) (T, error) {
	return float64(element.(dto.Element).Value), nil
}
