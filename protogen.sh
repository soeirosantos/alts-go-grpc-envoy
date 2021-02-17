#!/bin/bash
#
# Generate Go code from the protobuf definition
# It assumes protoc is located at `./.protoc` as installed by `install_protoc_grpc_go.sh`
#
# Usage: ./protogen.sh

set -e

export PATH="$PATH:$(go env GOPATH)/bin:$PWD/.protoc/bin"

declare -a services=("client" "checkout" "shipping" "payment")

for service in "${services[@]}"; do
  rm -rf "$service"/proto
  mkdir -p "$service"/proto

  protoc --go_out="$service"/proto \
         --go-grpc_out="$service"/proto \
         pb/checkout.proto
done
