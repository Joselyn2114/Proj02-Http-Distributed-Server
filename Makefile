#

.PHONY: test coverage

test:
	@go test ./... -timeout 30s

coverage:
	@go test -v ./... -timeout 30s -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
