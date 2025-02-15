FROM golang:1.23.2 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

EXPOSE 8080

CMD ["./app"]

