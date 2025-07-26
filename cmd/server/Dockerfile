FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 3000

CMD ["go", "run", "cmd/server/main.go", "--config=/app/cmd/server/example_server.env"]