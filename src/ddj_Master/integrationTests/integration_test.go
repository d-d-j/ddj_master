package integrationTests

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

func Benchmark_Insert_One_Value(b *testing.B) {

	json_data := "{\"tag\":1,	\"metric\":2,	\"time\":1383501407,	\"value\":0.5}"

	json_data_reader := strings.NewReader(json_data)

	// TODO: use cfg values instead of hardcoded port
	response, err := http.Post("http://localhost:8888/data", "application/json", json_data_reader)

	if response.StatusCode != 202 {
		b.Log("Status Code:", response.StatusCode, err)
		b.Fail()
	}
}

func Benchmark_Insert_N_Values(b *testing.B) {
	client := &http.Client{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		// TODO: cleanup the object declarations before resetTimer to get reliable execution times
		data := fmt.Sprintf("{\"tag\":%d, \"metric\":%d, \"time\":%d, \"value\":%f}", rand.Int31n(2)+1, rand.Int31n(4)+1, i, float32(i)+2.5)

		// TODO: use cfg values instead of hardcoded port
		req, err := http.NewRequest("POST", "http://localhost:8888/data", strings.NewReader(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Connection", "Keep-Alive")
		response, err := client.Do(req)

		if response.StatusCode != 202 {
			b.Log("Status Code:", response.StatusCode, err)
			b.Fail()
		}
		response.Body.Close()
	}
}
