#!/bin/bash

set -e

cd $(dirname $0)

CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

strip main

docker build -t testcmd .

