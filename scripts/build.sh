#!/usr/bin/env sh

set -e

docker build -t 117503445/staticdns-builder -f Dockerfile.builder .
docker run --rm -v "$(pwd):/workspace" -v "GOCACHE:/root/.cache/go-build" -v "GOMODCACHE:/go/pkg/mod" 117503445/staticdns-builder