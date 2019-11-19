#! /bin/bash

mkdir /results
mkdir /source/test
cp -r /data/* /source/test/
mocha --recursive --reporter=mocha_reporter -O fileName=/results/result.json