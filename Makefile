#
# Biuld the project.
#
build: bin/ws

bin/ws: cmds/ws/ws.go ws.go
	go build -o bin/ws cmds/ws/ws.go

lint:
	gofmt -w ws.go && golint ws.go
	gofmt -w cmds/ws/ws.go && golint cmds/ws/ws.go

install: bin/ws ws.go
	env GOBIN=${HOME}/bin go install
	env GOBIN=${HOME}/bin go install cmds/ws/ws.go

clean: 
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

test:
	go test
	gocyclo -over 15 .

release:
	./mk-release.sh
