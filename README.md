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

## Sample queries

1. Insert single element
        Send JSON using method POST to address: http://localhost:8888/data
        JSON STRUCTURE:
           {
           "tag":1,                   <- int
           "metric":2,                <- int
           "time":1383501407,         <- int64 (unsigned)
           "value":0.5                <- float32
           }

        EXAMPLE:
        element="{\"series\":7,\"tag\":2,\"time\":`date -u +%s`,\"value\":0.5}"
		curl -X POST -d $element http://localhost:8888/data --header "Content-Type:application/json"

2. Select

		curl -G 'http://localhost:8888/data/metric/1,2,3/tag/3,4,5/time/10-20,30-60/aggregation/none'

3. SelectAll

        curl -G 'http://localhost:8888/data/metric/all/tag/all/time/0-100000000/aggregation/none'

4. FLush - some data can be still hold in buffers, use it to load them all to store

		curl -X POST 'http://localhost:8888/data/flush'





## Integration tests
1. `run master`
2. `run node`
3. `make integrationTest`

## Value Aggregation Types

* none
* sum
* max - maximum
* min - minimum
* avg - average
* std - standard deviation
* var - variance
* int - integral

## Series Aggregation Types

* sum     (only sum now supported - more will come soon...)

