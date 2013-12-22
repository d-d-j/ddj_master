package main

import (
	"testing"
	"net/http"
	"strings"
)

func Benchmark_Insert_One_Value(b *testing.B) {

	json_data := "{\"tag\":1,	\"metric\":2,	\"time\":1383501407,	\"value\":0.5}"

	json_data_reader  := strings.NewReader(json_data)


	response, err := http.Post("http://localhost:8888/data", "application/json", json_data_reader)

	b.Log("Status Code:",response.StatusCode,  err)

	if response.StatusCode != 202 {
		b.Fail()
	}
}

func Benchmark_Insert_One_Value_And_Select(b *testing.B) {

	data := "{\"tag\":2,	\"metric\":3,	\"time\":1383501407,	\"value\":0.5}"

	json_data := strings.NewReader(data)

	http.Post("http://localhost:8888/data", "application/json", json_data)

//	response, err := http.Get("http://localhost:8888/data/metric/2/tag/3/time/10-20,30-60/aggregation/none")

//	b.Log("Status Code:",response.StatusCode,  err)
}
