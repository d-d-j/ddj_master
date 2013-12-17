DDJ_Master
==========

## Set up
1. Install go
	- version 1.1.2 or above required

    wget http://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.2.linux-amd64.tar.gz

append `export PATH=$PATH:/usr/local/go/bin` to `/etc/profile`

2. Build `make`
3. Run `./DDJ_Master [-port=<port>]`

## Sample query

1. Insert

		curl -X POST -d "{\"series\":7,\"tag\":2,\"time\":`date -u +%s`,\"value\":0.5}" http://localhost:8888/data --header "Content-Type:application/json"

2. Select All

		curl http://localhost:8888/data