#!/bin/bash

# update dependencies
#docker run --rm -it -v $PWD:/go/src/github.com/xanderstrike/goplaxt -w /go/src/github.com/xanderstrike/goplaxt treeder/glide update

# build binary
docker run --rm -v "$PWD":/go/src/github.com/xanderstrike/goplaxt -w /go/src/github.com/xanderstrike/goplaxt iron/go:dev go build -o goplaxt-docker

# build docker image
docker build -t xanderstrike/goplaxt .
