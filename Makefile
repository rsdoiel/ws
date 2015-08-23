build: ws wsjs wskeygen wsinit slugify unslugify 

ws: cmds/ws/ws.go cfg/cfg.go fsengine/fsengine.go ottoengine/ottoengine.go
	go build cmds/ws/ws.go

wsjs: cmds/wsjs/wsjs.go cfg/cfg.go fsengine/fsengine.go ottoengine/ottoengine.go
	go build cmds/wsjs/wsjs.go

wskeygen: cmds/wskeygen/wskeygen.go cfg/cfg.go keygen/keygen.go
	go build cmds/wskeygen/wskeygen.go

wsinit: cmds/wsinit/wsinit.go cfg/cfg.go keygen/keygen.go
	go build cmds/wsinit/wsinit.go

slugify: cmds/slugify/slugify.go slug/slug.go
	go build cmds/slugify/slugify.go

unslugify: cmds/unslugify/unslugify.go slug/slug.go
	go build cmds/unslugify/unslugify.go


install: ws slugify unslugify test wskeygen
	go install cmds/ws/ws.go
	go install cmds/wsjs/wsjs.go
	go install cmds/slugify/slugify.go
	go install cmds/unslugify/unslugify.go
	go install cmds/wskeygen/wskeygen.go
	go install cmds/wsinit/wsinit.go

clean: 
	if [ -f ws ]; then rm ws; fi
	if [ -f wsjs ]; then rm wsjs; fi
	if [ -f slugify ]; then rm slugify; fi
	if [ -f unslugify ]; then rm unslugify; fi
	if [ -f wskeygen ]; then rm wskeygen; fi
	if [ -f wsinit ]; then rm wsinit; fi

test: slug
	cd slug && go test
