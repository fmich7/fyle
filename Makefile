run-server:
	@cd server && go run main.go
run-client: 
	@cd client && go run main.go
test:
	cd server && go test ./...
	cd client && go test ./...

.PHONY: client server