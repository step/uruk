#! /bin/bash

mkdir /results
mocha --recursive --reporter=json > /results/result.json