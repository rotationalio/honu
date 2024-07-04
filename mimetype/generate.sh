#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

if [[ ! -d "./v1" ]]; then
    mkdir v1
fi

MODULE="github.com/rotationalio/honu/mimetype/v1"
MOD="github.com/rotationalio/honu/mimetype/v1;mimetype"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1 \
    --go_opt=module=${MODULE} \
    --go_opt=Mmimetype/v1/mimetype.proto="${MOD}" \
    --go_opt=Mmimetype/v1/charset.proto="${MOD}" \
    mimetype/v1/mimetype.proto \
    mimetype/v1/charset.proto