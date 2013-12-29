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

2. Select

		curl -G 'http://localhost:8888/data/metric/1,2,3/tag/3,4,5/time/10-20,30-60/aggregation/none'

3. SelectAll

        curl -G 'http://localhost:8888/data/metric/all/tag/all/time/10-20,30-60/aggregation/none'

4. FLush

		curl -X POST 'http://localhost:8888/data/flush'



## Integration tests
1. run master
2. run node
3. 'make integrationTest'
