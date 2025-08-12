#!/usr/bin/make -f

test: fmt
	go test -race -cover -timeout=1s -count=1 ./...

fmt:
	@go version && go fmt ./... && go mod tidy

install: test
	go install github.com/mdw-katas/coding-challenges.fyi/cmd/...
