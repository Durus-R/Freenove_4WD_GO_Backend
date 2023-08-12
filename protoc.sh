#! /bin/bash
mkdir -p dist
protoc --go_out=dist/ \
  --go_opt=paths=source_relative \
  --go-grpc_out=dist/ \
  --go-grpc_opt=paths=source_relative proto/*.proto

