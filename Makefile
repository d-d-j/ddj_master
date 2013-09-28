all:
	GOPATH=`pwd` go build DDJ_Master
run:
	GOPATH=`pwd` go run src/DDJ_Master/master.go
clean:
	rm DDJ_Master