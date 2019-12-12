# step 1: build
FROM golang:1.12.5 as build-step
# FROM golang:1.12.5-alpine3.9 as build-step
# RUN apk add --update --no-cache build-base ca-certificates git

# fatal: unable to access 'https://go.googlesource.com/sys/'
# maybe outage
# https://github.com/golang/go/issues/32395
ENV GOPROXY=https://proxy.golang.org

RUN mkdir /go-app
WORKDIR /go-app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'bindatafs' -a -o /go/bin/qor-demo
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/seeds config/db/seeds/main/main.go

# COPY ["app", "vendor", "config", "./tmp/"]
RUN  mkdir tmp && cp -r app vendor config ./tmp && rm tmp/app/*/*.go && rm tmp/config/*/*/*.go

# -----------------------------------------------------------------------------
# step 2: exec
# FROM phusion/baseimage:0.11
# FROM golang:1.12.5-alpine3.9
FROM alpine:3.9.4

ENV TZ='Asia/Shanghai'

RUN apk update && apk add --no-cache openssl ca-certificates curl netcat-openbsd tzdata && \
  cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
  rm -rf /var/cache/apk/*

RUN mkdir /qor
WORKDIR /qor
COPY --from=build-step /go/bin/qor-demo ./qor
COPY --from=build-step /go-app/tmp .
# COPY --from=build-step /go/pkg/mod /go/pkg/mod
COPY --from=build-step /go-app/rules.grl.tmpl .

CMD ./qor
