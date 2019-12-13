#! /bin/bash

mkdir /results
cp -r /data/.eslintrc.json /source/
cd /source
eslint . -f json -o /results/result.json