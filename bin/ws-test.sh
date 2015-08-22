#!/bin/bash
#
# Run the contents of this repository with ws
#
if [ -f "cmds/ws/ws.go" ]; then
    echo "Running test from "$(pwd)
else
    echo "Can't find cmds/ws/ws.go"
    exit 1
fi
if [ -f "./ws" ]; then
    echo "Stale ws found. Removing."
    /bin/rm ./ws
fi
# build ws
echo "Building cmds/ws/ws.go"
go build cmds/ws/ws.go
# Try to run ws in this directory.
echo "Try to run ws"
if [ -f "./ws" ]; then
    echo "Build OK, starting with otto and TLS false"
    ./ws -otto=false -tls=false -D doc
else
    echo "Something went wrong building ws."
fi
