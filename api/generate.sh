#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

if [[ ! -d "./v1" ]]; then
    mkdir v1
fi

MODULE="github.com/rotationalio/honu/api/v1"
APIMOD="github.com/rotationalio/honu/api/v1;api"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 --go-grpc_out=./v1 \
    --go_opt=module=${MODULE} \
    --go-grpc_opt=module=${MODULE} \
    --go_opt=Mhonu/v1/honu.proto="${APIMOD}" \
    --go-grpc_opt=Mhonu/v1/honu.proto="${APIMOD}" \
    honu/v1/honu.proto