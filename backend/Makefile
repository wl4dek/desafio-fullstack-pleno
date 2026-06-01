.PHONY: run build test clean

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./... -v -p=1

test-unit:
	go test ./internal/auth/... ./internal/children/... ./internal/statistics/... ./internal/config/... -v

test-integration:
	TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/full-stack-test?sslmode=disable" \
	  go test ./internal/database/... ./internal/children/... ./internal/statistics/... -v -p=1 -run Integration

clean:
	rm -rf bin/
