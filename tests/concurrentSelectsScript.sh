#!/bin/bash

python concurrenSelects.py -s -nv 10000 -ns 1000 -nt 20 > concurrentSelectsResults10000.txt
python concurrenSelects.py -s -nv 100000 -ns 1000 -nt 20 > concurrentSelectsResults100000.txt
python concurrenSelects.py -s -nv 1000000 -ns 1000 -nt 20 > concurrentSelectsResults1000000.txt
