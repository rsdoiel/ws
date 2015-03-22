
ws: cmds/ws/ws.go
	go build cmds/ws/ws.go

install: ws
	go install cmds/ws/ws.go

clean:
	rm ws

test: ws
	./bin/ws-demo.sh
