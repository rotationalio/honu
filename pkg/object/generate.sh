#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/honu/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

MODULE="github.com/rotationalio/honu/pkg/object/v1"
MOD="github.com/rotationalio/honu/pkg/object/v1;object"
OUT="./v1"

if [[ ! -d $OUT ]]; then
    mkdir $OUT
fi

protoc -I=${PROTOS} \
    --go_out=${OUT} \
    --go_opt=module="${MODULE}" \
    --go_opt=Mobject/v1/object.proto="${MOD}" \
    object/v1/object.proto