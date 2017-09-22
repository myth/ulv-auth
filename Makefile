.PHONY: all build test install install_deps

default: build

test: install_deps
	go test -cover $(go list ./... | grep -v /vendor/)

build: install_deps
	mkdir -p build
	go build -v -o build/ulv

install:
	install build/ulv /usr/bin/ulv

install_deps:
	go get github.com/golang/dep/cmd/dep
	go get golang.org/x/tools/cmd/cover
	dep ensure

