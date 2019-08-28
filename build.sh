#!/usr/bin/env bash

set -x

ROOT_DIR=`pwd`
echo ${ROOT_DIR}/bin

if [ ! -d ${ROOT_DIR}/bin ]; then
  mkdir bin
fi

cd ${ROOT_DIR}/demo/connstat && go build -o ${ROOT_DIR}/bin/connstat