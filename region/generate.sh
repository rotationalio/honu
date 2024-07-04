#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

if [[ ! -d "./v1" ]]; then
    mkdir v1
fi

MODULE="github.com/rotationalio/honu/region/v1"
APIMOD="github.com/rotationalio/honu/region/v1;region"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 --go-grpc_out=./v1 \
    --go_opt=module=${MODULE} \
    --go-grpc_opt=module=${MODULE} \
    --go_opt=Mregion/v1/region.proto="${APIMOD}" \
    --go-grpc_opt=Mregion/v1/region.proto="${APIMOD}" \
    region/v1/region.proto