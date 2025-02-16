build: client server
client:
	cd cmd/client && go build -o ../../bin/fyle-client
server:
	cd cmd/server && go build -o ../../bin/fyle-server

run-server:
	@cd cmd/server && go run main.go
run-client: 
	@cd cmd/client && go run main.go

.PHONY: client server