#!/bin/bash
#
# Run the demo site in the demo folder with _ws_
#
go run ws.go \
    -docroot="demo/static" \
    -otto=true \
    -otto-path="demo/dynamic"
