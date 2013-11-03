DDJ_Master
==========

## Set up
1. Install go
	- version 1.1.2 required
	- http://golang.org/doc/install
	- Install GoRest `go get code.google.com/p/gorest`
2. Build `make`
3. Run `./DDJ_Master [-port=<port>]`

## Sample query

1. Insert

		curl -X POST -d "{\"series\":7,\"tag\":2,\"time\":`date -u +%s`,\"value\":0.5}" http://localhost:8888/data --header "Content-Type:application/json"

2. Select All

		curl http://localhost:8888/data