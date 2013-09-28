all:
	export GOPATH=`pwd`
	@echo 'GOPATH = $(GOPATH)'
	go build DDJ_Master
clean:
	rm DDJ_Master