build: ws wsjs wskeygen wsinit slugify unslugify range reldate shorthand

ws: cmds/ws/ws.go src/cfg/cfg.go src/fsengine/fsengine.go src/ottoengine/ottoengine.go src/cli/cli.go src/wslog/wslog.go
	go build cmds/ws/ws.go

wsjs: cmds/wsjs/wsjs.go src/cfg/cfg.go src/fsengine/fsengine.go src/ottoengine/ottoengine.go src/cli/cli.go src/wslog/wslog.go
	go build cmds/wsjs/wsjs.go

wskeygen: cmds/wskeygen/wskeygen.go src/cfg/cfg.go src/keygen/keygen.go src/cli/cli.go
	go build cmds/wskeygen/wskeygen.go

wsinit: cmds/wsinit/wsinit.go src/cfg/cfg.go src/keygen/keygen.go src/cli/cli.go
	go build cmds/wsinit/wsinit.go

slugify: cmds/slugify/slugify.go src/slugify/slugify.go src/cli/cli.go
	go build cmds/slugify/slugify.go

unslugify: cmds/unslugify/unslugify.go src/slugify/slugify.go src/cli/cli.go
	go build cmds/unslugify/unslugify.go

range: cmds/range/range.go
	go build cmds/range/range.go

reldate: cmds/reldate/reldate.go
	go build cmds/reldate/reldate.go

shorthand: cmds/shorthand/shorthand.go src/shorthand/shorthand.go
	go build cmds/shorthand/shorthand.go

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

clean: 
	if [ -f ws ]; then rm ws; fi
	if [ -f wsjs ]; then rm wsjs; fi
	if [ -f slugify ]; then rm slugify; fi
	if [ -f unslugify ]; then rm unslugify; fi
	if [ -f wskeygen ]; then rm wskeygen; fi
	if [ -f wsinit ]; then rm wsinit; fi
	if [ -f range ]; then rm range; fi
	if [ -f reldate ]; then rm reldate; fi
	if [ -f shorthand ]; then rm shorthand; fi


test:
	cd src/slugify && go test
	cd src/shorthand && go test
