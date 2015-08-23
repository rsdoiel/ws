#!/bin/bash
#
GO=$(which go)
if [ "$GO" = "" ]; then
    echo "Must install Golang first"
    echo "See http://golang.org for instructions"
    exit 1
fi
# Install dependent libraries
go get github.com/robertkrimen/otto
# go get github.com/knieriem/markdown
# go get github.com/zhemao/glisp
# go get github.com/kedebug/LispEx
make
make test
