.PHONY: run build test clean

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./... -v

clean:
	rm -rf bin/
