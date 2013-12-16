all:
	export GOPATH=`pwd`
	go build ddj_Master
debug:
	go build -ldflags "-s" ddj_Master

run:
	./ddj_Master

clean:
	rm ddj_Master
