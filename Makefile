
ws: ws.go
	go build ws.go

install: ws
	go install ws.go

clean:
	rm ws

