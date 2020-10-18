TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

build:
		go build -ldflags "-X main.buildVersion=$(TAG)" -o bin/qox .

run:
		bin/qox

test:
		go test ./...
