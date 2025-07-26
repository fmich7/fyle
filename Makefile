build: build-client build-server

build-client: 
	go build -o ./bin/client ./cmd/fyle

build-server: 
	go build -o ./bin/server ./cmd/server

server:
	@go run ./cmd/server/main.go --config=/home/fmich/projects/fyle/cmd/server/example_server.env

install client: 
	go install ./cmd/fyle

test:
	go test -v -timeout 10s -race ./pkg/... ./tests/...
coverage:
	go test -timeout 10s -race ./pkg/... ./tests/... -coverprofile=cover.out
	go tool cover -html=cover.out

.PHONY: client server
