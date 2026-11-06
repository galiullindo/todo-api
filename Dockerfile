# builder
FROM golang:1.25.6-alpine3.23 AS builder
ARG VERSION

RUN apk add --no-cache gcc musl-dev

WORKDIR /go/src/github.com/galiullindo/todo-api

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -a -o build/api -ldflags="-s -w -X main.Version=${VERSION}" ./cmd/api

# image
FROM alpine:3.23

RUN apk upgrade --no-cache \
    && apk add --no-cache sqlite-libs ca-certificates
    

COPY --from=builder /go/src/github.com/galiullindo/todo-api/build/api  /usr/local/bin/api
RUN ln -s /usr/local/bin/api /api

EXPOSE 8080

ENTRYPOINT [ "api" ]
