#!/bin/bash
#
GO=$(which go)
if [ "$GO" = "" ]; then
    echo "Must install Golang first"
    echo "See http://golang.org for instructions"
    exit 1
fi
# Install dependent libraries
# Add from shorthand for generating website
go get github.com/rsdoiel/shorthand
## Used by ws js support
go get github.com/caltechlibrary/otto
make
make test
