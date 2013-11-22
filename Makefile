get:
	export GOPATH=`pwd`
	go get code.google.com/p/gorest
	go get code.google.com/p/gcfg
	go get code.google.com/p/log4go
all:
	go build ddj_Master
run: all
	./ddj_Master
clean:
	rm ddj_Master
