default: lint test build

tools: ## Install the tools used to test and build
	@echo "==> Installing build tools"
	GO111MODULE=off go get -u github.com/ahmetb/govvv
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

build: ## Build Sherpa for development purposes
	@echo "==> Running $@..."
	govvv build -o sherpa ./cmd -version $(shell git describe --tags --abbrev=0 $(git rev-list --tags --max-count=1) |cut -c 2- |awk '{print $1}')+dev -pkg github.com/jrasell/sherpa/pkg/build

test: ## Run the Sherpa test suite with coverage
	@echo "==> Running $@..."
	@go test ./... -cover -v -tags -race \
		"$(BUILDTAGS)" $(shell go list ./... | grep -v vendor)

acctest: ## Run the Sherpa acceptance test suite
	@echo "==> Running $@..."
	@SHERPA_ACC=1 go test ./test -count 1 -v -mod vendor

release: ## Trigger the release build script
	@echo "==> Running $@..."
	@goreleaser --rm-dist

.PHONY: lint
lint: ## Run golangci-lint
	@echo "==> Running $@..."
	golangci-lint run cmd/... pkg/...

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Sherpa make commands:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
