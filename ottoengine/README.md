ottoengine
==========

A golang module that wraps the Otto JS VM by Robert Krimen to implement dynamic
route handlers in the _ws_ web server.

# How it works

If you use the command line options *-otto=true* and *-otto-path=SOME_DIRECTORY* or set
the environment variables *WS_OTTO=true* and *WS_OTTO_PATH* then those JavaScript files
found in the *WS_OTTO_PATH* will be used to create route handles in the _ws_ web server.

The return value of the JavaScript files renders as the content body sent back to the
web bowser.

