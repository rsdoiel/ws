#!/bin/bash
#
# Run the demo using shell environment variables to configure.
#
export WS_HOST=localhost
export WS_PORT=8000
export WS_DOCROOT=$(pwd)/demo/static
export WS_OTTO=true
export WS_OTTO_PATH=$(pwd)/demo/dynamic

echo "ws will listen for $WS_HOST:$WS_PORT"
echo "Static document root $WS_DOCROOT"
echo "Otto is $WS_OTTO, otto scripts $WS_OTTO_PATH"

echo "Ready to run Demo Y/N "
read YES_NO
if [ "$YES_NO" = "Y" ] || [ "$YES_NO" = "y" ]; then
    echo "Building ws webserver"
    make
    echo "Point your browser at http://$WS_HOST:$WS_PORT and click around to test."
    echo "Press ctrl-c to kill web server"
    ./ws
fi


