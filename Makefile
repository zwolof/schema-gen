BINARY  := schema-gen
BIN_DIR := bin
GO_FILES := $(shell find . -name '*.go' -not -path './$(BIN_DIR)/*' -not -path './exported*')

.PHONY: help run build fmt fmt-check vet tidy race check clean install-tools

.DEFAULT_GOAL := help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

run: ## Regenerate schemas into ./exported
	go run .

build: $(BIN_DIR)/$(BINARY) ## Build binary to ./bin/

$(BIN_DIR)/$(BINARY): $(GO_FILES)
	@mkdir -p $(BIN_DIR)
	go build -o $@ .

fmt: ## Format Go source (gofmt + goimports if available)
	@gofmt -s -w .
	@command -v goimports >/dev/null 2>&1 \
		&& goimports -w -local go-csitems-parser . \
		|| echo "  (goimports not found — run 'make install-tools')"

fmt-check: ## Fail if any file needs formatting (for CI)
	@unformatted=$$(gofmt -s -l .); \
		if [ -n "$$unformatted" ]; then \
			echo "needs formatting:"; echo "$$unformatted"; exit 1; \
		fi

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go.mod
	go mod tidy

race: ## Run with -race detector
	go run -race .

profile: ## Run with CPU + heap profiling; writes cpu.prof and mem.prof
	go run . -cpuprofile cpu.prof -memprofile mem.prof
	@echo
	@echo "  Inspect: go tool pprof -http=:8080 cpu.prof"
	@echo "           go tool pprof -http=:8081 mem.prof"

check: fmt-check vet ## Format-check + vet (pre-commit sanity)

clean: ## Remove build artefacts and generated output
	rm -rf $(BIN_DIR)/ exported/ cpu.prof mem.prof

install-tools: ## Install optional formatting tools
	go install golang.org/x/tools/cmd/goimports@latest
