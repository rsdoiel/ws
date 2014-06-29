#!/bin/bash
CWD=$(pwd)
mkdir -p $GOPATH/bin
go build ws.go
if [ -f ws ]; then
   mv ws $GOPATH/bin/
else
   echo "Something went wrong building ws."
   exit 1
fi
cd $CWD
