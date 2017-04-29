all:
	go build
debug:
	go build -ldflags "-s"

run: all
	./ddj_master

test: all
	go test $(shell go list ./... | grep -v /vendor/) -cover	

integrationTest: all
	go test ./integrationTests/ -bench="."

clean:
	rm -f ddj_master
	go clean ./
