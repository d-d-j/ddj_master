all:
	go build ddj_Master
debug:
	go build -ldflags "-s" ddj_Master

run: all
	./ddj_Master

test: all
	go test ./...

clean:
	rm ddj_Master
