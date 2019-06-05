# step 1: build
FROM golang:1.12.5 as build-step
# FROM golang:1.12.5-alpine3.9 as build-step
# RUN apk add --update --no-cache build-base ca-certificates git

RUN mkdir /go-app
WORKDIR /go-app
COPY go.mod .
COPY go.sum .
ENV GOPROXY=https://goproxy.io
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'bindatafs' -a -o /go/bin/qor-demo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/seeds config/db/seeds/main.go config/db/seeds/seeds.go

# -----------------------------------------------------------------------------
# step 2: exec
# FROM phusion/baseimage:0.11
FROM golang:1.12.5-alpine3.9
# FROM alpine:3.9.4

RUN apk add --no-cache openssl
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# RUN mkdir /go-app
# WORKDIR /go-app
COPY --from=build-step /go/bin/qor-demo /go/bin/qor-demo
COPY --from=build-step /go/bin/seeds /go/bin/seeds
COPY --from=build-step /go/pkg/mod /go/pkg/mod
EXPOSE 7000
COPY app ./app
RUN rm app/*/*.go
COPY config/locales ./config/locales
COPY config/db/seeds/data ./config/db/seeds/data

CMD dockerize -wait tcp://db:5432 -timeout 30s /go/bin/qor-demo

