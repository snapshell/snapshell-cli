# Makefile for snapshell-cli

.PHONY: test
test:
	go test ./... -v


.PHONY: install
install:
	go build -o snapshell main.go
	cp snapshell ~/go/bin