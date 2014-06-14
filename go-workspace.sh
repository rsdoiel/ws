#!/bin/bash

if [ "$GOROOT" = "" ]; then
    echo "Missing Golang or GOROOT env."
    exit 1
fi
echo "Found Golang, setting workspace to $(pwd)"
export GOPATH=$(pwd)
