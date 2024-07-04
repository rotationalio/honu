#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

if [[ ! -d "./v1" ]]; then
    mkdir v1
fi

MODULE="github.com/rotationalio/honu/replica/v1"
APIMOD="github.com/rotationalio/honu/replica/v1;replica"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 --go-grpc_out=./v1 \
    --go_opt=module=${MODULE} \
    --go-grpc_opt=module=${MODULE} \
    --go_opt=Mreplica/v1/replica.proto="${APIMOD}" \
    --go-grpc_opt=Mreplica/v1/replica.proto="${APIMOD}" \
    replica/v1/replica.proto