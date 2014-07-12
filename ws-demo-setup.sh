#!/bin/bash
export WS_HOST=localhost
export WS_DOCROOT=$(pwd)/demo/static
export WS_OTTO=true
export WS_OTTO_PATH=$(pwd)/demo/dynamic

echo "ws will listen for $WS_HOST"
echo "static documnet roo $WS_DOCROOT"
echo "Otto is $WS_OTTO, otto scripts $WS_OTTO_PATH"
if [ -f ./ws ]; then
    ./ws
else
    ws
fi
