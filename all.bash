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
cd extra
go build ws-gencert.go
if [ -f ws-gencert ]; then
mv ws-gencert $GOPATH/bin/
else
    echo "Something went wrong building ws-gencert."
    exit 1
fi
cd $CWD
