#!/bin/bash
while :
do
	value=$[100 + (RANDOM % 100)]$[1000 + (RANDOM % 1000)]
	value=$[RANDOM % 15].${value:1:2}${value:4:3}
	series=$[(RANDOM % 10)]
	tag=$[(RANDOM % 10)]
	curl -X POST -d "{\"series\":$series,\"tag\":$tag,\"time\":`date -u +%s`,\"value\":$value}" http://localhost:8081/series/id/4/data --header "Content-Type:application/json"
done
