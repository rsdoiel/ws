
# Installation

1. Set your GOPATH to the root project directory
2. Create a _bin_ directory in your root project directory
3. cd src/ws
4. go build ws.go
5. If everything builds OK
    - go install ws.go
6. cd extra
7. go build ws-gencert.go
8. If everything OK
    - mv ws-gencert $GOPATH/bin

Your runable binaries should be in $GOPATH/bin.

