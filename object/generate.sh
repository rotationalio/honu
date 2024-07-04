#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

if [[ ! -d "./v1" ]]; then
    mkdir v1
fi

MODULE="github.com/rotationalio/honu/object/v1"
MOD="github.com/rotationalio/honu/object/v1;object"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 \
    --go_opt=module=${MODULE} \
    --go_opt=Mobject/v1/object.proto="${MOD}" \
    object/v1/object.proto