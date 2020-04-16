#!/usr/bin/env bash

cd adapter/handler
find . -name main.go -type f \
 | xargs -n 1 dirname \
 | xargs -n 1 -I@ bash -c "CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -o ../../bin/@/main ./@"