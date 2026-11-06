# builder
FROM golang:1.25.6-alpine3.23 AS builder
ARG VERSION

RUN apk add --no-cache gcc musl-dev

WORKDIR /todo-api

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -a -o build/server -ldflags="-s -w -X main.Version=${VERSION}" ./cmd/api

# image
FROM alpine:3.23

RUN apk upgrade --no-cache \
    && apk add --no-cache sqlite-libs ca-certificates
    
WORKDIR /todo-api
RUN mkdir ./data

COPY --from=builder /todo-api/build/server  ./server

EXPOSE 8080

ENTRYPOINT [ "server" ]
