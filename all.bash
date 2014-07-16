#!/bin/bash
if [ "$GOPATH" = "" ]; then
    echo "Missing GOPATH"
    exit 1
fi
if [ "$GOBIN" = "" ]; then
    echo "Missing GOBIN"
    exit 1 
fi

CWD=$(pwd)
if [ -f "$GOBIN/ws" ]; then
    rm "$GOBIN/ws"
fi
go get github.com/robertkrimen/otto
go install ws.go
if [ -f "$GOBIN/ws" ]; then
   echo "Installed in $GOBIN/ws"
else
   echo "Something went wrong. Missing $GOBIN/ws."
   exit 1
fi
cd $CWD
