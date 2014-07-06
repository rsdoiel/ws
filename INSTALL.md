
# Installation

If you already have your golang workspace setup you can skip steps 1.
I have assumed your running some type of *NIX with Bash or Sh.

1. Setup our Go Workspace
    a. Create your work space directory, PROJECT_NAME_HERE would be the name of your project. E.g.
        - mkdir -p $HOME/gocode/src
    b. Set GOPATH environment variable. E.g.
        - export GOPATH=$HOME/gocode
    c. Set GOBIN environment variable. E.g.
        - export GOBIN=$GOPATH
    d. Make sure Otto is available
2. Checkout _ws_ from github
    a. Change to $GOPATH/src
        - cd $GOPATH/src
    b. Clone repository
        - git clone git@github.com:rsdoiel/ws.git
2. Change your to _ws_ directory. E.g.
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
    a. Start the _ws_ webserver
        - ./ws
    b. Try the following URLs in your favorite web browser
        - http://localhost:8000
        - http://localhost:8000/helloworld
        - http://localhost:8000/test
        - http://localhost:8000/json
        - http://localhost:8000/html
7. If everything builds OK
    - ./all.bash
8. Generate your self-signed SSL cert if needed. E.g.
    - sudo ws -keygen
9. Add the default "WS_*" environment variables to your login scripts. E.g.
    - export WS_HOST="localhost"
    - export WS_PORT="8000"
    - export WS_DOCROOT="static"
    - export WS_TLS=false
    - export WS_CERT=$HOME/etc/ws/cert.pem
    - export WS_KEY=$HOME/etc/ws/key.pem
    - export WS_OTTO=false
    - export WS_OTTO_PATH="otto"
10. Make sure GOBIN is in your path. E.g.
    - export PATH="$GOBIN:$PATH"
11. Your runable binaries should be in $GOBIN. Try
    - ls -l $GOBIN


