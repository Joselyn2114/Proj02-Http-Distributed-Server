.PHONY: test coverage integration all

test:
	@go test ./... -timeout 30s

coverage:
	@go test -v ./worker/... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

integration:
	@go test -timeout 30s -run Integration

all: test integration
