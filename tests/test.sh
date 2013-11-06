#!/bin/bash

#!/bin/bash

usage() { cat << EOF
Usage: $0
OPTIONS
    -n Number of requests to perform
    -c Number of multiple requests to make
    -P POST
    -G GET
EOF
}

requests_count=
concuret_request_count=
post=
get=

while getopts "n:c:PG" o; do
    case "${o}" in
        n)
            requests_count=${OPTARG}
            ;;
        c)
            concuret_request_count=${OPTARG}
            ;;
        P)
            post=true
            get = false
            ;;
        G)
            get=true
            post = false
            ;;
        \? ) echo "Unknown option: -$OPTARG" >&2; exit 1;;
        :  ) echo "Missing option argument for -$OPTARG" >&2; exit 1;;
        *)
            usage
            ;;
    esac
done

if $get
    then
    ab -n $requests_count -c $concuret_request_count http://localhost:8888/data
fi

if $post
    then
    ab -n $requests_count -c $concuret_request_count -p insert.json -T "'application/x-www-form-urlencoded'"  http://localhost:8888/data
fi