build: build-client build-server
build-client:
	cd cmd/client && go build -o ../../bin/fyle-client
build-server:
	cd cmd/server && go build -o ../../bin/fyle-server

server:
	@cd cmd/server && go run main.go
client: build-client
	@./bin/fyle-client

test:
	go test -v -timeout 10s -race ./pkg/... ./tests/...
coverage:
	go test -timeout 10s -race ./pkg/... ./tests/... -coverprofile=cover.out
	go tool cover -html=cover.out

.PHONY: client server
