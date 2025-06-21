.PHONY: test integration all

test:
	@go test ./... -timeout 30s

integration:
	@go test -timeout 30s -run Integration

all: test integration
