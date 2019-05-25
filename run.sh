#! /usr/local/bin/bash

dropdb qor_example

createdb qor_example

go run config/db/seeds/main.go config/db/seeds/seeds.go

go run main.go

