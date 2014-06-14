#!/bin/bash
go install ws
cd extra
go build ws-gencert.go
mv ws-gencert $GOPATH/bin/
cd ..
