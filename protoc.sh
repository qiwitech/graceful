#!/bin/bash

set -e
cd `dirname $0`

services="./systest/api_v1 ./systest/pluto"

if [ $# -gt 0 ]; then
	services=$@
fi

# Rebuild custom generator
go install github.com/qiwitech/graceful/protoc-gen-httpgw

for PKG in $services; do
for PROTO in ${PKG}/*.proto; do
  for TYPE in gogo httpgw; do
    protoc \
      -I${GOPATH}/src \
      -I${PKG} \
      --${TYPE}_out=:${PKG}/ \
      ${PROTO}
  done
#  for INTERFACE in ${PKG}/*.httpgw.go; do
#    mockgen -package $(basename $PKG) -source ${INTERFACE} > ${PKG}/$(basename ${INTERFACE})_mocks_test.go
#  done
done
done

## Generate HTTP-client (github.com/go-swagger/go-swagger)
#cd $GOPATH/src/github.com/qiwitech/graceful/test
#swagger generate client -f ../v1/api.swagger.json
#go get ./client/...
#cd -
