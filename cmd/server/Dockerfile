FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY cmd/server/example_server.env .

RUN CGO_ENABLED=0 GOOS=linux go build -o /fyle-server cmd/server/main.go

EXPOSE 3000

CMD ["/fyle-server"]