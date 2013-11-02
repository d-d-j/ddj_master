get:
	export GOPATH=`pwd`
	go get code.google.com/p/gorest
	go get code.google.com/p/gcfg
all:
	go build DDJ_Master
run: all
	./DDJ_Master
clean:
	rm DDJ_Master