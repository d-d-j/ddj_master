__author__ = 'dud'

#!/usr/bin/python

import sys
import json
import math
import argparse
import requests


def postValues(number, tag, metric, url, func):
    funcString = "sin()"
    if func == math.cos:
        funcString = "cos()"

    print "About to post", number, funcString, "values", "to", url
    print "\ttag:", tag
    print "\tmetric:", metric

    s = requests.Session()

    headers = {'content-type': 'application/json'}
    req = requests.Request('POST', url,
                           headers=headers)
    for i in range(1, number + 1):
        payload = {"tag": tag, "metric": metric, "time": i, "value": func(0.01 * i)}
        req.data = json.dumps(payload)

        prepped = req.prepare()

        try:
            response = s.send(prepped)
            if response.status_code != 202:
                print "Got response: ", response.status_code
                sys.exit(1)

        except requests.exceptions.RequestException as e:    # This is the correct syntax
            print e
            sys.exit(1)

        if i % 100 == 0:
            sys.stdout.write(str(i) + " ")
            sys.stdout.flush()


def main():
    func = math.sin

    parser = argparse.ArgumentParser(description='Post some data into DB.')
    parser.add_argument('-u', '--url', dest='url', default='http://localhost:8888/data',
                        help='url to post data (default: http://localhost:8888/data)')
    parser.add_argument('-c', '--cos', help='use cosine values instead of sine', action='store_true')
    parser.add_argument('-n', type=int, dest='n', default=1000, help='number of records to post (default: 1000)')
    parser.add_argument('-t', '--tag', type=int, dest='tag', default=1, help='tag of data to post (default: 1)')
    parser.add_argument('-m', '--metric', type=int, dest='metric', default=1,
                        help='metric of data to post (default: 1)')

    args = parser.parse_args()

    if args.n:
        n = args.n
    if args.tag:
        tag = args.tag
    if args.metric:
        metric = args.metric
    if args.url:
        url = args.url
    if args.cos:
        func = math.cos

    postValues(n, tag, metric, url, func)


if __name__ == "__main__":
    main()