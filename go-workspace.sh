#!/bin/bash

if [ "$GOROOT" = "" ]; then
    echo "Missing Golang or GOROOT env."
    exit 1
fi
WORK_PATH=$(pwd)
cd ../..
export GOPATH=$(pwd)
export GOBIN=$GOPATH/bin
mkdir -p $GOBIN
cd $WORK_PATH
echo "Found Golang, set workspace to $GOPATH and $GOBIN"
