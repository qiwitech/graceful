#!/bin/bash

set -e
cd `dirname $0`

# Rebuild generator
go install github.com/qiwitech/graceful/protoc-gen-httpgw

for GENERATOR in gogo httpgw; do
  protoc \
    -I${GOPATH}/src \
    -I. \
    --${GENERATOR}_out=:./ \
    srv.proto
done

mockgen -package testservice -source srv.httpgw.go > gomock_test.go

#    -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
