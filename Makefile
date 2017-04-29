all: test ddj_master

ddj_master:
	go build
debug:
	go build -ldflags "-s"

run: ddj_master
	./ddj_master

test: ddj_master
	go test $(shell go list ./... | grep -v /vendor/) -cover	

integrationTest: test
	go test ./integrationTests/ -bench="."

clean:
	rm -f ddj_master
	go clean ./
