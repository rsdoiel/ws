#!/bin/bash
#
GO=$(which go)
if [ "$GO" = "" ]; then
    echo "Must install Golang first"
    echo "See http://golang.org for instructions"
    exit 1
fi
# Install dependent libraries
## Used by wsjs
go get github.com/robertkrimen/otto
## Use by wsmarkdown
go get github.com/russross/blackfriday
go get github.com/microcosm-cc/bluemonday
# go get github.com/knieriem/markdown
# go get github.com/zhemao/glisp
# go get github.com/kedebug/LispEx
make
make test
