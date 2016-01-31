#!/bin/bash
#
GO=$(which go)
if [ "$GO" = "" ]; then
    echo "Must install Golang first"
    echo "See http://golang.org for instructions"
    exit 1
fi
# Install dependent libraries
go get -u github.com/rsdoiel/otto
# Now run tests
go test
# Now build everything
mkdir -p bin
go build -o bin/ws cmds/ws/ws.go
go build -o bin/wsinit cmds/wsinit/wsinit.go
echo "Binaries written to $(pwd)/bin"
