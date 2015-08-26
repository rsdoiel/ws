#
# Biuld the project.
#
build: bin/ws bin/wsjs bin/wskeygen bin/wsinit bin/slugify bin/unslugify bin/shorthand bin/wsmarkdown

bin/ws: cmds/ws/ws.go cfg/cfg.go fsengine/fsengine.go ottoengine/ottoengine.go cli/cli.go wslog/wslog.go
	go build -o bin/ws cmds/ws/ws.go

bin/wsjs: cmds/wsjs/wsjs.go cfg/cfg.go fsengine/fsengine.go ottoengine/ottoengine.go cli/cli.go wslog/wslog.go
	go build -o bin/wsjs cmds/wsjs/wsjs.go

bin/wskeygen: cmds/wskeygen/wskeygen.go cfg/cfg.go keygen/keygen.go cli/cli.go
	go build -o bin/wskeygen cmds/wskeygen/wskeygen.go

bin/wsinit: cmds/wsinit/wsinit.go cfg/cfg.go keygen/keygen.go cli/cli.go
	go build -o bin/wsinit cmds/wsinit/wsinit.go

bin/slugify: cmds/slugify/slugify.go slugify/slugify.go cli/cli.go
	go build -o bin/slugify cmds/slugify/slugify.go

bin/unslugify: cmds/unslugify/unslugify.go slugify/slugify.go cli/cli.go
	go build -o bin/unslugify cmds/unslugify/unslugify.go

bin/shorthand: cmds/shorthand/shorthand.go shorthand/shorthand.go
	go build -o bin/shorthand cmds/shorthand/shorthand.go

bin/wsmarkdown: cmds/wsmarkdown/wsmarkdown.go
	go build -o bin/wsmarkdown cmds/wsmarkdown/wsmarkdown.go


install: bin/ws bin/wsjs bin/wskeygen bin/wsinit bin/slugify bin/unslugify bin/shorthand
	go install cmds/ws/ws.go
	go install cmds/wsjs/wsjs.go
	go install cmds/wskeygen/wskeygen.go
	go install cmds/wsinit/wsinit.go
	go install cmds/slugify/slugify.go
	go install cmds/unslugify/unslugify.go
	go install cmds/shorthand/shorthand.go
	go install cmds/wsmarkdown/wsmarkdown.go

clean: 
	if [ -f bin/ws ]; then rm bin/ws; fi
	if [ -f bin/wsjs ]; then rm bin/wsjs; fi
	if [ -f bin/slugify ]; then rm bin/slugify; fi
	if [ -f bin/unslugify ]; then rm bin/unslugify; fi
	if [ -f bin/wskeygen ]; then rm bin/wskeygen; fi
	if [ -f bin/wsinit ]; then rm bin/wsinit; fi
	if [ -f bin/shorthand ]; then rm bin/shorthand; fi
	if [ -f bin/wsmarkdown ]; then rm bin/wsmarkdown; fi

test:
	cd slugify && go test
	cd shorthand && go test

