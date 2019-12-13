#! /bin/bash

mkdir /results
mkdir /source/__test
cp -r /data/* /source/__test/
cd /source
npm install
npm install shelljs
mocha --recursive --reporter=mocha_reporter -O fileName=/results/result.json __test