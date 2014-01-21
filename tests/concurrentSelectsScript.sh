#!/bin/bash

python concurrentSelects.py -s -nv 1000 -ns 1000 -mt 25> concurrentSelectsResults1000.txt
python concurrentSelects.py -s -nv 10000 -ns 1000 -mt 25> concurrentSelectsResults10000.txt
python concurrentSelects.py -s -nv 100000 -ns 1000 -mt 25 > concurrentSelectsResults100000.txt
python concurrentSelects.py -s -nv 1000000 -ns 1000 -mt 25 > concurrentSelectsResults1000000.txt
