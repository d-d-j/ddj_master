package main

import (
	"testing"
	"net/http"
	"strings"
	"fmt"
	"math/rand"
)


func Benchmark_Insert_One_Value(b *testing.B) {

	json_data := "{\"tag\":1,	\"metric\":2,	\"time\":1383501407,	\"value\":0.5}"

	json_data_reader  := strings.NewReader(json_data)
	response, err := http.Post("http://localhost:8888/data", "application/json", json_data_reader)

	if response.StatusCode != 202 {
		b.Log("Status Code:",response.StatusCode,  err)
		b.Fail()
	}
}

func Benchmark_1000_Values(b *testing.B) {
	client := &http.Client{}

	for i := 0; i < 1000; i++ {
		data := fmt.Sprintf("{\"tag\":%d, \"metric\":%d, \"time\":%d, \"value\":%f}", rand.Int31n(2) + 1, rand.Int31n(4) + 1, i, float32(i)+2.5)

		req, err := http.NewRequest("POST", "http://localhost:8888/data", strings.NewReader(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Connection", "Keep-Alive")
		response, err := client.Do(req)

		if response.StatusCode != 202 {
			b.Log("Status Code:",response.StatusCode,  err)
			b.Fail()
		}
		response.Body.Close()
	}
}

