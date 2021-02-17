#!/bin/bash
#
# Install protoc and the Go protoc grpc plugin
# It keeps the binaries in a local folder called .protoc
# Use `protogen.sh` to generate the stubs based on this installation

set -e

PB_REL="https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-osx-x86_64.zip"

rm -rf $PWD/.protoc

curl -L $PB_REL -o protoc.zip
mkdir -p $PWD/.protoc
unzip protoc.zip -d $PWD/.protoc

export GO111MODULE=on
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc

rm protoc.zip