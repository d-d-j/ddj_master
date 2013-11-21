#!/bin/bash

rm threads1.temp
rm threads10.temp
rm threads20.temp
rm threads50.temp
rm threads100.temp

ab -n 100000 -c 1 -g threads1.temp -p insert.json -T "'application/x-www-form-urlencoded'" http://127.0.0.1:8888/data
ab -n 100000 -c 10 -g threads10.temp -p insert.json -T "'application/x-www-form-urlencoded'" http://127.0.0.1:8888/data
ab -n 100000 -c 20 -g threads20.temp -p insert.json -T "'application/x-www-form-urlencoded'" http://127.0.0.1:8888/data
ab -n 100000 -c 50 -g threads50.temp -p insert.json -T "'application/x-www-form-urlencoded'" http://127.0.0.1:8888/data
ab -n 100000 -c 100 -g threads100.temp -p insert.json -T "'application/x-www-form-urlencoded'" http://127.0.0.1:8888/data


