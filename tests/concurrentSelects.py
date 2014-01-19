__author__ = 'dud'

#!/usr/bin/python

import sys
import json
import math
import requests
import threading
import datetime


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

        if i % 100 == 0:
            sys.stdout.write(str(i) + " ")
            sys.stdout.flush()

    print ""

    s.close()
    requests.post(url + '/flush')

    print "flushed buffers"


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
    # print "finished running " + str(numberOfSelects) + " integrals"

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

    return {'integrals': integralsTime, 'sums': sumsTime}


def main():
    numberOfValues = 10000
    numberOfSelects = 1000
    postValues(numberOfValues, 1, 1, 'http://localhost:8888/data', math.sin)

    results = {}
    for i in range(1, 16):
        results[i] = runSelectsWithNInsertThreads(i, numberOfSelects, numberOfValues)
        print str(i) + " threads FINISHED"

    print "finished testing for " + str(numberOfSelects) + "selects on " + str(numberOfValues) + "values"
    print "threads" + "\t" + "integrals" + "\t" + "sums"
    for key, value in results.items():
        print str(key) + "\t" + str(value['integrals']) + '\t' + str(value['sums'])


if __name__ == "__main__":
    main()