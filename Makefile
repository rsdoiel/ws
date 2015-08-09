
slugify: cmds/slugify/slugify.go slug/slug.go
	go build cmds/slugify/slugify.go

unslugify: cmds/unslugify/unslugify.go slug/slug.go
	go build cmds/unslugify/unslugify.go

ws: cmds/ws/ws.go
	go build cmds/ws/ws.go

build: ws slugify unslugify

install: ws slugify unslugify test
	go install cmds/ws/ws.go
	go install cmds/slugify/slugify.go
	go install cmds/unslugify/unslugify.go

clean: 
	if [ -f ws ]; then rm ws; fi
	if [ -f slugify ]; then rm slugify; fi
	if [ -f unslugify ]; then rm unslugify; fi

test: slug
	cd slug && go test
