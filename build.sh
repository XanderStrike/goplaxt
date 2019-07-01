#!/bin/bash

# build binary
docker run --rm -e GO111MODULE=on -v "$PWD":/go/src/github.com/xanderstrike/goplaxt -w /go/src/github.com/xanderstrike/goplaxt iron/go:dev go build -o goplaxt-docker

# build docker image
docker build -t xanderstrike/goplaxt .

rm goplaxt-docker
