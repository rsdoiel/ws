#
# Biuld the project.
#

PROG = ws

build: cmd/$(PROG)/$(PROG).go $(PROG).go

$(PROG).go:
	env CGO_ENABLED=0 go build

cmd/$(PROG)/$(PROG).go:
	env CGO_ENABLED=0 go build -o bin/$(PROG) cmds/$(PROG)/$(PROG).go

lint:
	gofmt -w ws.go && golint ws.go
	gofmt -w cmds/ws/ws.go && golint cmds/ws/ws.go

install:
	env GOBIN=$(HOME)/bin go install cmds/$(PROG)/$(PROG).go

clean: 
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f ws-release.zip ]; then /bin/rm ws-release.zip; fi

test:
	go test
	gocyclo -over 15 $(PROG).go
	gocyclo -over 15 cmds/$(PROG)/$(PROG).go

save:
	git commit -am "Quick save"
	git push origin master

release:
	./mk-release.bash

website:
	./mk-website.bash

publish:
	./mk-website.bash
	./publish.bash
