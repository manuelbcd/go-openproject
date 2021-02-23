.PHONY: test-coverage
test: ## Runs all unit, integration and example tests.
	go test -race -v ./...

.PHONY: lint
lint: ## Linter code
    golint ./..

.PHONY: vet
vet: ## Runs go vet (to detect suspicious constructs).
    go vet ./...

.PHONY: fmt
fmt: ## Runs go fmt (to check for go coding guidelines).
	gofmt -d -s .

.PHONY: staticcheck
staticcheck: ## Runs static analysis to prevend bugs, foster code simplicity, performance and editor integration.
	go get -u honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

.PHONY: all
all: test vet fmt staticcheck ## Runs all source code quality targets (like test, vet, fmt, staticcheck)
