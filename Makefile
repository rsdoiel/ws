#
# Biuld the project.
#
build: bin/ws bin/wsjs bin/wskeygen bin/wsinit bin/slugify bin/unslugify bin/wsmarkdown

bin/ws: cmd/ws/ws.go cfg/cfg.go fsengine/fsengine.go ottoengine/ottoengine.go cli/cli.go wslog/wslog.go
	go build -o bin/ws cmd/ws/ws.go

bin/wsjs: cmd/wsjs/wsjs.go cfg/cfg.go fsengine/fsengine.go ottoengine/ottoengine.go cli/cli.go wslog/wslog.go
	go build -o bin/wsjs cmd/wsjs/wsjs.go

bin/wskeygen: cmd/wskeygen/wskeygen.go cfg/cfg.go keygen/keygen.go cli/cli.go
	go build -o bin/wskeygen cmd/wskeygen/wskeygen.go

bin/wsinit: cmd/wsinit/wsinit.go cfg/cfg.go keygen/keygen.go cli/cli.go
	go build -o bin/wsinit cmd/wsinit/wsinit.go

bin/slugify: cmd/slugify/slugify.go slugify/slugify.go cli/cli.go
	go build -o bin/slugify cmd/slugify/slugify.go

bin/unslugify: cmd/unslugify/unslugify.go slugify/slugify.go cli/cli.go
	go build -o bin/unslugify cmd/unslugify/unslugify.go


bin/wsmarkdown: cmd/wsmarkdown/wsmarkdown.go
	go build -o bin/wsmarkdown cmd/wsmarkdown/wsmarkdown.go

lint:
	gofmt -w cfg/cfg.go && golint cfg/cfg.go
	gofmt -w cli/cli.go && golint cli/cli.go
	gofmt -w fsengine/fsengine.go && golint fsengine/fsengine.go
	gofmt -w keygen/keygen.go && golint keygen/keygen.go
	gofmt -w ottoengine/ottoengine.go && golint ottoengine/ottoengine.go
	gofmt -w slugify/slugify.go && golint slugify/slugify.go
	gofmt -w wslog/wslog.go && golint wslog/wslog.go
	gofmt -w prompt/prompt.go && golint prompt/prompt.go
	gofmt -w cmd/ws/ws.go && golint cmd/ws/ws.go
	gofmt -w cmd/wsjs/wsjs.go && golint cmd/wsjs/wsjs.go
	gofmt -w cmd/wsinit/wsinit.go && golint cmd/wsinit/wsinit.go
	gofmt -w cmd/wskeygen/wskeygen.go && golint cmd/wskeygen/wskeygen.go
	gofmt -w cmd/wsmarkdown/wsmarkdown.go && golint cmd/wsmarkdown/wsmarkdown.go
	gofmt -w cmd/slugify/slugify.go && golint cmd/slugify/slugify.go
	gofmt -w cmd/unslugify/unslugify.go && golint cmd/unslugify/unslugify.go


install: bin/ws bin/wsjs bin/wskeygen bin/wsinit bin/slugify bin/unslugify
	go install cmd/ws/ws.go
	go install cmd/wsjs/wsjs.go
	go install cmd/wskeygen/wskeygen.go
	go install cmd/wsinit/wsinit.go
	go install cmd/slugify/slugify.go
	go install cmd/unslugify/unslugify.go
	go install cmd/wsmarkdown/wsmarkdown.go

clean: 
	if [ -f bin/ws ]; then rm bin/ws; fi
	if [ -f bin/wsjs ]; then rm bin/wsjs; fi
	if [ -f bin/slugify ]; then rm bin/slugify; fi
	if [ -f bin/unslugify ]; then rm bin/unslugify; fi
	if [ -f bin/wskeygen ]; then rm bin/wskeygen; fi
	if [ -f bin/wsinit ]; then rm bin/wsinit; fi
	if [ -f bin/wsmarkdown ]; then rm bin/wsmarkdown; fi

test:
	cd slugify && go test

