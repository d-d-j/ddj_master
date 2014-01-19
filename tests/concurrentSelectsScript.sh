#!/bin/bash

python concurrentSelects.py -s -nv 10000 -ns 1000 -mt 20 > concurrentSelectsResults10000.txt
python concurrentSelects.py -s -nv 100000 -ns 1000 -mt 20 > concurrentSelectsResults100000.txt
python concurrentSelects.py -s -nv 1000000 -ns 1000 -mt 20 > concurrentSelectsResults1000000.txt
