all:
	go build ddj_Master
debug:
	go build -ldflags "-s" ddj_Master

run: all
	./ddj_Master

test: all
	go test -cover ./...

integrationTest: all
	go test ./src/ddj_Master/integrationTests/ -bench="."

clean:
	rm -f ddj_Master
	go clean ddj_Master