include .env
export

build: deps
	go build -o bin/client cmd/client/client.go

run: build
	./bin/client

deps:
	go mod tidy