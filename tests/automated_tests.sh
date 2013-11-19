#!/bin/bash

rm test.temp
rm newtest.tst
echo Requests per second: 0 [#/sec] (mean) concurency 0 >> newtest.tst

./test.sh -n 20000 -c 1 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 1 >> newtest.tst

./test.sh -n 20000 -c 2 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 2 >> newtest.tst

./test.sh -n 20000 -c 3 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 3 >> newtest.tst

./test.sh -n 20000 -c 4 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 4 >> newtest.tst

./test.sh -n 20000 -c 5 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 5 >> newtest.tst

./test.sh -n 20000 -c 6 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 6 >> newtest.tst

./test.sh -n 20000 -c 8 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 8 >> newtest.tst

./test.sh -n 20000 -c 10 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 10 >> newtest.tst

./test.sh -n 20000 -c 15 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 15 >> newtest.tst

./test.sh -n 20000 -c 30 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 30 >> newtest.tst

./test.sh -n 20000 -c 40 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 40 >> newtest.tst

./test.sh -n 20000 -c 60 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 60 >> newtest.tst

./test.sh -n 20000 -c 100 -P > test.temp
record=$(grep '#/sec' test.temp)
echo $record concurency 100 >> newtest.tst
