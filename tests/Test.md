# Test using ap

Require `sudo apt-get install apache2-utils`


It will make 10K request in 4 threads



### Get all data

		ab -n 10000 -c 4 -p http://localhost:8888/data

### Insert data

		ab -n 10000 -c 4 -p insert.json -T "'application/x-www-form-urlencoded'"  http://localhost:8888/data

