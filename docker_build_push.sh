#!/bin/bash

set -e

service=$1
version=$2

cd $service

docker build -t soeirosantos/alts-"$service":"$version" .
docker push soeirosantos/alts-"$service":"$version"
