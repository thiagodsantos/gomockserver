FROM golang:alpine3.19

WORKDIR /app

CMD ["go", "run", "mock-server.go"]
