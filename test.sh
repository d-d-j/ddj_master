#!/bin/bash
while :
do
	curl -X POST -d "{\"series\":7,\"tag\":2,\"time\":`date -u +%s`,\"value\":0.5}" http://localhost:8081/series/id/4/data --header "Content-Type:application/json"
done
