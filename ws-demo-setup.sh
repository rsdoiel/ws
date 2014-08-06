#!/bin/bash
#
export WS_HOST=localhost:8000
export WS_DOCROOT=$(pwd)/demo/static
export WS_OTTO=true
export WS_OTTO_PATH=$(pwd)/demo/dynamic

echo "ws will listen for $WS_HOST"
echo "Static document root $WS_DOCROOT"
echo "Otto is $WS_OTTO, otto scripts $WS_OTTO_PATH"
if [ -f ./ws ]; then
    ./ws
else
    ws
fi

