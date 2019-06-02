# step 1: build
FROM golang:1.12.5 as build-step
# FROM golang:1.12.5-alpine3.9 as build-step
# RUN apk add --update --no-cache build-base ca-certificates git

RUN mkdir /go-app
WORKDIR /go-app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'bindatafs' -o /go-app/qor-example
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/seeds config/db/seeds/main.go config/db/seeds/seeds.go
CMD ["/go-app/qor-example"]
