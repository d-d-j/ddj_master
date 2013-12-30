package integrationTests

import (
	"ddj_Master/dto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type RestResponse struct {
	Error  string
	TaskId int64
	Data   []dto.Element
}

const (
	NUMBER_OF_TAGS_PER_METRICS int = 4
	NUMBER_OF_METRICS          int = 4
	INSERTED_DATA              int = 1000
	// TODO: use cfg values instead of hardcoded port
	HOST string = "http://localhost:8888/data"
)

var (
	insert_completed bool
	data             []string
	expected         []dto.Element
)

func SetUp(b *testing.B) {

	if insert_completed != true {
		client := &http.Client{}
		data = make([]string, INSERTED_DATA)
		expected = make([]dto.Element, INSERTED_DATA)
		for i := 0; i < INSERTED_DATA; i++ {
			e := *dto.NewElement(int32(i%NUMBER_OF_TAGS_PER_METRICS), int32(i%NUMBER_OF_METRICS), int64(i), 1.0)
			data[i] = fmt.Sprintf("{\"tag\":%d, \"metric\":%d, \"time\":%d, \"value\":%f}", e.Metric, e.Tag, e.Time, e.Value)
			expected[i] = e
		}

		for i := 0; i < INSERTED_DATA; i++ {

			req, err := http.NewRequest("POST", HOST, strings.NewReader(data[i]))
			if err != nil {
				b.Log("Error occurred: ", err)
				b.FailNow()
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Connection", "Keep-Alive")
			response, err := client.Do(req)

			if response.StatusCode != 202 {
				b.Log("Status Code:", response.StatusCode, err)
				b.FailNow()
			}
			response.Body.Close()
		}
		insert_completed = true
		b.Log("Inserted ", INSERTED_DATA, " elements")
	}
}

//TODO: Add more select test for filters and aggregations

func Benchmark_Select_First_Value(b *testing.B) {

	SetUp(b)

	response := Select("/metric/0/tag/0/time/0-0/aggregation/none", b)

	Assert(response, expected, b)
}

func Benchmark_Select_Ten_Inserted_Values_From_All_Tags_And_Metrics(b *testing.B) {

	SetUp(b)

	response := Select(fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/none", INSERTED_DATA-10, INSERTED_DATA), b)

	Assert(response, expected[INSERTED_DATA-10:INSERTED_DATA], b)
}

func Benchmark_Select_All(b *testing.B) {

	SetUp(b)

	response := Select(fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/none", 0, INSERTED_DATA), b)

	Assert(response, expected[:INSERTED_DATA], b)
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

func Benchmark_Select_All_Max(b *testing.B) {

	SetUp(b)

	response := Select(fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/max", 0, INSERTED_DATA), b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	if response.Data[0].String() != expected[0].Value.String() {
		b.Error("Got ", response.Data, " when expected ", expected[INSERTED_DATA].Value)
	}
}

func Benchmark_Select_All_Min(b *testing.B) {

	SetUp(b)

	response := Select(fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/min", 0, INSERTED_DATA), b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	if response.Data[0].String() != expected[0].Value.String() {
		b.Error("Got ", response.Data, " when expected ", expected[0].Value)
	}
}

func Benchmark_Select_Sum(b *testing.B) {

	SetUp(b)

	response := Select(fmt.Sprintf("/metric/all/tag/all/time/%d-%d/aggregation/min", 0, 1), b)

	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}

	if response.Data[0].String() != expected[0].Value.String() {
		b.Error("Got ", response.Data, " when expected ", expected[0].Value)
	}
}

func Assert(response RestResponse, expectedValues []dto.Element, b *testing.B) {
	if len(response.Data) < 1 {
		b.Log("Nothing returned")
		b.FailNow()
	}
	if len(response.Data) > len(expectedValues) {
		b.Error("Too many values returned. Expected less or equal than ", len(expectedValues), " but got ", len(response.Data))
		b.Log("Did you forget to restart server?")
		b.FailNow()
	}
	for index, element := range response.Data {
		if expectedValues[index].String() != element.String() {
			b.Error("Expected ", expectedValues[index], " but got ", element)
		}
	}
}
