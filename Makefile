all:
	go build ddj_Master
get:
	export GOPATH=`pwd`
	go get code.google.com/p/gorest
	go get code.google.com/p/gcfg
	go get code.google.com/p/log4go
debug:
	go build -ldflags "-s" ddj_Master

run:
	./ddj_Master

clean:
	rm ddj_Master
