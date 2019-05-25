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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'bindatafs' -o /go/bin/qor-example
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/seeds config/db/seeds/main.go config/db/seeds/seeds.go

# -----------------------------------------------------------------------------
# step 2: exec
# FROM phusion/baseimage:0.11
FROM golang:1.12.5

RUN mkdir /go-app
WORKDIR /go-app
COPY --from=build-step /go/bin/qor-example /go-app/qor-example
COPY --from=build-step /go/bin/seeds /go-app/seeds
COPY app ./app
COPY public ./public
COPY config ./config

CMD ["/go-app/qor-example"]

