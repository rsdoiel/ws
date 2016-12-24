#
# Biuld the project.
#

PROJECT = ws

PROG = ws

VERSION = $(shell grep -m1 'Version = ' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep "* " | cut -d\  -f 2)


build:
	env CGO_ENABLED=0 go build -o bin/$(PROG) cmds/$(PROG)/$(PROG).go

lint:
	gofmt -w ws.go && golint ws.go
	gofmt -w cmds/ws/ws.go && golint cmds/ws/ws.go

install:
	env GOBIN=$(HOME)/bin go install cmds/$(PROG)/$(PROG).go

clean: 
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

test:
	go test
	gocyclo -over 15 $(PROG).go
	gocyclo -over 15 cmds/$(PROG)/$(PROG).go

save:
	git commit -am "Quick save"
	git push origin $(BRANCH)

release:
	./mk-release.bash

website:
	./mk-website.bash

publish:
	./mk-website.bash
	./publish.bash
