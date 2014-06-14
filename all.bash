#!/bin/bash
CWD=$(pwd)
go build
if [ -f ws ]; then
   mv ws $GOPATH/bin/
else
   echo "Something went wrong building ws."
   exit 1
fi
cd extra
go build
if [ -f ws-genecert ]; then
mv ws-gencert $GOPATH/bin/
else
    echo "Something went wrong building ws-gencert."
    exit 1
fi
cd $CWD
