#! /bin/bash
mkdir -p dist
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/go/bin
protoc --go_out=dist/ \
  --go_opt=paths=source_relative \
  --go-grpc_out=dist/ \
  --go-grpc_opt=paths=source_relative proto/*.proto

