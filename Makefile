.DEFAULT_GOAL:=help
.PHONY: build clean cover help lint test

build: clean ## Build binary
	@go build -o awsipranges ./cmd/awsipranges/main.go

clean: ## Clean up binary, coverage report
	@rm -f awsipranges coverage.txt

cover: test ## Display test coverage report
	@go tool cover -func=coverage.txt

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

lint: ## Lint source code
	@golangci-lint run ./...

test: ## Run unit tests
	@go test -v -coverprofile=coverage.txt ./...

tidy: ## Tidy cmd module
	@cd cmd/awsipranges && go mod tidy
