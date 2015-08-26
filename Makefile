build: ws wsjs wskeygen wsinit slugify unslugify range reldate shorthand wsmarkdown
	mkdir -p bin

ws: cmds/ws/ws.go src/cfg/cfg.go src/fsengine/fsengine.go src/ottoengine/ottoengine.go src/cli/cli.go src/wslog/wslog.go
	go build -o bin/ws cmds/ws/ws.go

wsjs: cmds/wsjs/wsjs.go src/cfg/cfg.go src/fsengine/fsengine.go src/ottoengine/ottoengine.go src/cli/cli.go src/wslog/wslog.go
	go build -o bin/wsjs cmds/wsjs/wsjs.go

wskeygen: cmds/wskeygen/wskeygen.go src/cfg/cfg.go src/keygen/keygen.go src/cli/cli.go
	go build -o bin/wskeygen cmds/wskeygen/wskeygen.go

wsinit: cmds/wsinit/wsinit.go src/cfg/cfg.go src/keygen/keygen.go src/cli/cli.go
	go build -o bin/wsinit cmds/wsinit/wsinit.go

slugify: cmds/slugify/slugify.go src/slugify/slugify.go src/cli/cli.go
	go build -o bin/slugify cmds/slugify/slugify.go

unslugify: cmds/unslugify/unslugify.go src/slugify/slugify.go src/cli/cli.go
	go build -o bin/unslugify cmds/unslugify/unslugify.go

range: cmds/range/range.go
	go build -o bin/range cmds/range/range.go

reldate: cmds/reldate/reldate.go
	go build -o bin/reldate cmds/reldate/reldate.go

shorthand: cmds/shorthand/shorthand.go src/shorthand/shorthand.go
	go build -o bin/shorthand cmds/shorthand/shorthand.go

wsmarkdown: cmds/wsmarkdown/wsmarkdown.go
	go build -o bin/wsmarkdown cmds/wsmarkdown/wsmarkdown.go

install: ws wsjs wskeygen wsinit slugify unslugify range reldate shorthand
	go install cmds/ws/ws.go
	go install cmds/wsjs/wsjs.go
	go install cmds/wskeygen/wskeygen.go
	go install cmds/wsinit/wsinit.go
	go install cmds/slugify/slugify.go
	go install cmds/unslugify/unslugify.go
	go install cmds/range/range.go
	go install cmds/reldate/reldate.go
	go install cmds/shorthand/shorthand.go
	go install cmds/wsmarkdown/wsmarkdown.go

clean: 
	if [ -f bin/ws ]; then rm bin/ws; fi
	if [ -f bin/wsjs ]; then rm bin/wsjs; fi
	if [ -f bin/slugify ]; then rm bin/slugify; fi
	if [ -f bin/unslugify ]; then rm bin/unslugify; fi
	if [ -f bin/wskeygen ]; then rm bin/wskeygen; fi
	if [ -f bin/wsinit ]; then rm bin/wsinit; fi
	if [ -f bin/range ]; then rm bin/range; fi
	if [ -f bin/reldate ]; then rm bin/reldate; fi
	if [ -f bin/shorthand ]; then rm bin/shorthand; fi
	if [ -f bin/wsmarkdown ]; then rm bin/wsmarkdown; fi


test:
	cd src/slugify && go test
	cd src/shorthand && go test

