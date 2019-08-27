#!/usr/bin/env bash

set -x

ROOT_DIR=`pwd`
echo ${ROOT_DIR}/bin

if [ ! -d ${ROOT_DIR}/bin ]; then
  mkdir bin
fi

cd ${ROOT_DIR}/demo/connstat && GOOS=linux GOARCH=amd64 go build -o ${ROOT_DIR}/bin/connstat_linux
