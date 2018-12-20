#!/bin/bash

set -e
cd `dirname $0`

DIR=systest/integration/fuzz

rm -fr ${DIR}_var
mkdir -p ${DIR}_var/corpus
cp ${DIR}_sample/*.ase ${DIR}_var/corpus/

echo go-fuzz-build ...
go-fuzz-build -o ${DIR}_var/fuzz-build.zip github.com/qiwitech/graceful/systest/integration

echo go-fuzz ...
go-fuzz -bin=S{DIR}_var/fuzz-build.zip -workdir=${DIR}_var