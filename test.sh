#!/bin/bash

usage() { echo "Usage: $0 [-n repetition number] [-q call query]" 1>&2; exit 1; }

while getopts "n:q" o; do
    case "${o}" in
        n)
            s=${OPTARG}
            x=1;
            while [ $x -le $s ] ;
				do
					value=$[100 + (RANDOM % 100)]$[1000 + (RANDOM % 1000)]
					value=$[RANDOM % 15].${value:1:2}${value:4:3}
					series=$[(RANDOM % 10)]
					tag=$[(RANDOM % 10)]
					curl -X POST -d "{\"series\":$series,\"tag\":$tag,\"time\":`date -u +%s`,\"value\":$value}" http://localhost:8888/data --header "Content-Type:application/json"
					x=$[x + 1]
				done
            ;;
        q)
            curl http://localhost:8888/data
            ;;
        *)
            usage
            ;;
    esac
done
