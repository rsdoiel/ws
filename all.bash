#!/bin/bash
CWD=$(pwd)
mkdir -p $GOPATH/bin
GOBIN=$GOPATH/bin
go install ws.go
if [ -f $GOBIN/ws ]; then
   echo "ws installed in $GOBIN"
else
   echo "Something went wrong building ws."
   exit 1
fi
cd $CWD
