#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

if [[ ! -d "./v1" ]]; then
    mkdir v1
fi

MODULE="github.com/rotationalio/honu/pagination/v1"
MOD="github.com/rotationalio/honu/pagination/v1;pagination"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 \
    --go_opt=module=${MODULE} \
    --go_opt=Mpagination/v1/pagination.proto="${MOD}" \
    pagination/v1/pagination.proto