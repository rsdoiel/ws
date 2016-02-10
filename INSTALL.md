
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
3. Install the other Go packages needed
    - go get github.com/robertkrimen/otto
4. Compile _ws_, and _wsint_
    - go build
    - go build cmds/ws/ws.go
    - go test
5. Setup for testing
    - export WS_URL=http://localhost:8000
    - export WS_HTDOCS=$(pwd)/htdocs
    - export WS_JSDOCS=$(pwd)/jsdocs
    - export WS_SSL_KEY=$(pwd)/etc/site.key
    - export WS_SSL_PEM=$(pwd)/etc/site.pem
6. Test web server with JavaScript support
    + Start the _ws_ web server
        - ./ws
    + Try the following URLs in your favorite web browser
        - http://localhost:8000
        - http://localhost:8000/helloworld
    + Kill the web server with Ctrl-C
7. Experiment adding content, server side JavaScript, and restarting

