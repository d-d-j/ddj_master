all:
	go build ddj_Master
debug:
	go build -ldflags "-s" ddj_Master

run:
	./ddj_Master

test: all
	go test ./...

clean:
	rm ddj_Master
