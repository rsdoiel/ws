#!/bin/bash
#
GO=$(which go)
if [ "$GO" = "" ]; then
    echo 'Must install Golang first"
    echo "See http://golang.org for instructions"
    exit 1
fi
make
make test
