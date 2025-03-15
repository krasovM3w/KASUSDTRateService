#!/bin/bash
set -e

PROTO_DIR=./proto
PROTO_FILE=currency.proto
GEN_DIR=./internal/api/grpc/gen

mkdir -p ${GEN_DIR}
protoc --go_out=${GEN_DIR} --go_opt=paths=source_relative \
  --go-grpc_out=${GEN_DIR} --go-grpc_opt=paths=source_relative \
  -I ${PROTO_DIR} ${PROTO_DIR}/currency/${PROTO_FILE}n_DIR}