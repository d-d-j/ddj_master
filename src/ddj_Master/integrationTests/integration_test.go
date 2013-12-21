package main

import (
	"testing"
	"net/http"
	"strings"
)

func Test_Insert_One_Value(t *testing.T) {

	json_data := "{\"tag\":1,	\"metric\":2,	\"time\":1383501407,	\"value\":0.5}"

	b := strings.NewReader(json_data)


	response, err := http.Post("http://localhost:8888/data", "application/json", b)

	t.Log("Status Code:",response.StatusCode,  err)

	if response.Status != 202 {
		t.Fail()
	}
}
