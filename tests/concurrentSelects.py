__author__ = 'dud'

#!/usr/bin/python

import sys
import json
import math
import requests
import threading
import datetime
import argparse


def postValues(number, tag, metric, url, func):
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

        # if i % 100 == 0:
        #     sys.stdout.write(str(i) + " ")
        #     sys.stdout.flush()

    # print ""

    s.close()
    requests.post(url + '/flush')

    # print "flushed buffers"


def insertThread(metric, tag, stop_event):
    s = requests.Session()

    headers = {'content-type': 'application/json'}
    req = requests.Request('POST', "http://localhost:8888/data",
                           headers=headers)
    i = 0
    while not stop_event.is_set():
        payload = {"tag": tag, "metric": metric, "time": i, "value": math.sin(0.01 * i)}
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


def runSelects(aggregationType, numberOfSelects, numberOfValues):
    s = requests.Session()
    headers = {'content-type': 'application/json'}
    req = requests.Request('GET', "http://localhost:8888/data/metric/1/tag/1/time/0-" + str(numberOfValues) + "/aggregation/" + aggregationType,
                           headers=headers)

    start = datetime.datetime.now()
    for i in range(0, numberOfSelects):
        prepped = req.prepare()
        try:
            response = s.send(prepped)
            if response.status_code != 200:
                print "Got response: ", response.status_code
                sys.exit(1)

        except requests.exceptions.RequestException as e:
            print e
            sys.exit(1)

    end = datetime.datetime.now()

    s.close()

    return end-start


def runSelectsWithNInsertThreads(n, numberOfSelects, numberOfValues):
    threadStops = []
    for i in range(0, n):
        stop = threading.Event()
        threadStops.append(stop)
        thread = threading.Thread(target=insertThread, args=(1, i, stop))
        thread.start()

    integralsTime = runSelects('int', numberOfSelects, numberOfValues)
    sumsTime = runSelects('sum', numberOfSelects, numberOfValues)

    for stop in threadStops:
        stop.set()

    return {'integrals': integralsTime/numberOfSelects, 'sums': sumsTime/numberOfSelects}


def main():
    parser = argparse.ArgumentParser(description='Post some data into DB.')
    parser.add_argument('-nv', '--number-of-values', type=int, dest='numberOfValues', default=1000,
                        help='Number of values in one select')
    parser.add_argument('-ns', '--number-of-selects', type=int, dest='numberOfSelects', default=100,
                        help='Number of selects for each test')
    parser.add_argument('-mt', '--max-threads', type=int, dest='maxThreads', default=15,
                    help='maximum number of threads')


    parser.add_argument('-s', '--silent', help='use to disable progress logging', dest='silent', action='store_true')
    args = parser.parse_args()

    numberOfValues = args.numberOfValues
    numberOfSelects = args.numberOfSelects
    maxThreads = args.maxThreads
    silent = False
    if args.silent:
        silent = True
    postValues(numberOfValues, 1, 1, 'http://localhost:8888/data', math.sin)

    results = {}
    for i in range(0, maxThreads + 1):
        results[i] = runSelectsWithNInsertThreads(i, numberOfSelects, numberOfValues)
        if not silent:
            print str(i) + " threads FINISHED"

    print "finished testing for " + str(numberOfSelects) + "selects on " + str(numberOfValues) + "values"
    print "threads" + "\t" + "integrals" + "\t" + "sums"
    for key, value in results.items():
        print str(key) + "\t" + str(value['integrals']) + '\t' + str(value['sums'])


if __name__ == "__main__":
    main()