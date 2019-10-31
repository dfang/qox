# step 1: build
FROM golang:1.12.5 as build-step
# FROM golang:1.12.5-alpine3.9 as build-step
# RUN apk add --update --no-cache build-base ca-certificates git

RUN mkdir /go-app
WORKDIR /go-app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'bindatafs' -a -o /go/bin/qor-demo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/seeds config/db/seeds/main.go config/db/seeds/seeds.go

# -----------------------------------------------------------------------------
# step 2: exec
# FROM phusion/baseimage:0.11
# FROM golang:1.12.5-alpine3.9
FROM alpine:3.9.4

RUN apk update && apk add --no-cache openssl ca-certificates \
  && rm -rf /var/cache/apk/*

RUN mkdir /qor
WORKDIR /qor
COPY --from=build-step /go/bin/qor-demo /go/bin/seeds ./
# COPY --from=build-step /go/pkg/mod /go/pkg/mod
EXPOSE 7000
COPY app ./app
COPY vendor ./vendor
COPY config/locales ./config/locales
COPY config/db/seeds/data ./config/db/seeds/data
RUN rm app/*/*.go

CMD ./qor-demo
