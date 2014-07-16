
# Installation

## General go workspace setup

If you already have your golang workspace setup you can skip steps 1.
I have assumed your running some type of *NIX with Bash or Sh.

1. Setup our Go Workspace
    + Create your work space directory, PROJECT_NAME_HERE would be the name of your project. E.g.
        - mkdir -p $HOME/gocode/src
    + Set GOPATH environment variable. E.g.
        - export GOPATH=$HOME/gocode
    + Set GOBIN environment variable. E.g.
        - export GOBIN=$GOPATH
    + Make sure Otto is available

## Setup to build and test _ws_

2. Checkout _ws_ from github
    + Change to $GOPATH/src
        - cd $GOPATH/src
    + Clone repository
        - git clone git@github.com:rsdoiel/ws.git
2. Change to your _ws_ directory. E.g.
    - cd $GOPATH/src/ws
3. Make sure [otto](https://github.com/robertkrimen/otto) is available
    - go get github.com/robertkrimen/otto
4. Compile _ws_ web server. E.g.
    - go build ws.go
5. Setup for testing
    - export WS_HOST=localhost
    - export WS_PORT=8000
    - export WS_DOCROOT=$(pwd)
    - export WS_OTTO=true
    - export WS_OTTO_PATH=$(pwd)/demo
6. Test
    + Start the _ws_ webserver
        - ./ws
    + Try the following URLs in your favorite web browser
        - http://localhost:8000
        - http://localhost:8000/helloworld
        - http://localhost:8000/test
        - http://localhost:8000/json
        - http://localhost:8000/html
7. If everything builds OK
    - ./all.bash
8. Generate your self-signed SSL cert for testing (accepting defaults is fine). E.g.
    - ws -keygen
9. Add the default "WS_*" environment variables to your login scripts. E.g.
    - export WS_HOST="localhost"
    - export WS_PORT="8000"
    - export WS_DOCROOT=$(pwd)/demo/static
    - export WS_TLS=true
    - export WS_CERT=$(pwd)/demo/etc/ws/cert.pem
    - export WS_KEY=$(pwd)/demo/etc/ws/key.pem
    - export WS_OTTO=true
    - export WS_OTTO_PATH=$(pwd)/demo/dynamic
10. Make sure GOBIN is in your path. E.g.
    - export PATH="$GOBIN:$PATH"
11. Double checkt to make sure your runable _ws_ binary is in $GOBIN. Try
    - ls -l $GOBIN
12. Run _ws_ without SSL
    - ws
    - You should see some startup information on the console
    - Point your browser at http://localhost:8000 and run through your test URLs
    - Ctrl-c shuts down _ws_
13. Run _ws_ again with SSL support and test
    - ws -tls
14. If all has worked go through the Tutorial.md file and experiment.

