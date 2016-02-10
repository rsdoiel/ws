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
	go install
	go install cmds/ws/ws.go

clean: 
	if [ -f bin/ws ]; then rm bin/ws; fi

test:
	go test

