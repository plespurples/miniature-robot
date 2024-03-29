.PHONY: serverbuild

default:

up:
	CONFIG=./config.dev.yml go run cmd/websocket.go

serverbuild:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./dist/main-linux-amd64 ./cmd/websocket.go

deploy: serverbuild
	ssh -p 1014 jansvabik@s1.jetl.io 'mkdir -p /home/jansvabik/purples-websocket'
	ssh -t s1.jetl.io -p 1014 'sudo systemctl stop purples-websocket.service'
	scp -rO -P 1014 ./dist/main-linux-amd64 jansvabik@s1.jetl.io:/home/jansvabik/purples-websocket/main
	scp -rO -P 1014 ./config.prod.yml jansvabik@s1.jetl.io:/home/jansvabik/purples-websocket/config.prod.yml
	ssh -t s1.jetl.io -p 1014 'sudo systemctl restart purples-websocket.service'
