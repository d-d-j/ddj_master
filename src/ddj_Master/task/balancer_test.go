package task

import (
	"ddj_Master/restApi"
	"testing"
)

func Test_Update_Called_10_Times_Will_Fire_All_10_Workers(t *testing.T) {
	const (
		SIZE  int = 10
		TIMES int = 3
	)
	done := make(chan Worker)
	balancer := Balancer{}
	pool := MockWorkersPool(SIZE, 1, done, nil)
	balancer.pool = pool
	requestChan := make(chan restApi.RestRequest)
	go balancer.Balance(requestChan)
	go func() {
		for i := 0; i < SIZE*TIMES; i++ {
			go func() {
				requestChan <- restApi.RestRequest{}
			}()
		}
	}()
	finished := [SIZE]int{}
	for i := 0; i < SIZE*TIMES; i++ {
		w := <-done
		finished[w.Id()]++
	}
	for i := 0; i < SIZE; i++ {
		if finished[i] != TIMES {
			t.Error("Worker #", i, " was dispatched ", finished[i], " times but expected ", TIMES)
		}
	}
}
