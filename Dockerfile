# builder
FROM golang:1.25.6-alpine3.23 AS builder
ARG VERSION

RUN apk upgrade --no-cache \
    && apk add --no-cache gcc musl-dev

WORKDIR /todolist-rest-api

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=1 go build -a -o build/todolist -ldflags="-s -w -X main.Version=${VERSION}" ./cmd/todolist

# image
FROM alpine:3.23

RUN apk upgrade --no-cache \
    && apk add --no-cache sqlite-libs ca-certificates
    
WORKDIR /todolist-rest-api
COPY --from=builder /todolist-rest-api/build/todolist  ./todolist
COPY --from=builder /todolist-rest-api/configs/config.docker.yaml ./config.yaml
RUN mkdir ./data

EXPOSE 8083

ENTRYPOINT ["./todolist"]
