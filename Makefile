GO		    := $(shell which go)
GOPATH      := $(shell go env GOPATH)
GOBIN       := $(GOPATH)/bin
GOLINT      := $(GOBIN)/golangci-lint
GORELEASER  := $(GOBIN)/goreleaser
GOTESTFMT   := $(GOBIN)/gotestfmt

init:
	pre-commit install
	$(GO) version
	$(GO) mod download
	$(GO) mod tidy

clean:
	rm -rf dist || true
	mkdir dist

install-gotestfmt:
	test -e $(GOTESTFMT) || $(GO) install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
test: install-gotestfmt
	go test -json -v ./... 2>&1 | tee /tmp/gotest.log | $(GOTESTFMT)

lint: install-golangci-lint
	$(GOLINT) run --verbose ./...

format: install-golangci-lint
	$(GOLINT) run --verbose --fix ./...

upgrade:
	$(GO) get -u all

run:
	@$(GO) run ./

install-golangci-lint:
	test -e $(GOLINT) || $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
install-goreleaser:
	test -e $(GORELEASER) || $(GO) install github.com/goreleaser/goreleaser@latest

build: install-goreleaser
	$(GORELEASER) check
	$(GORELEASER) release --single-target

release: install-goreleaser
	$(GORELEASER) check
	$(GORELEASER) release --snapshot --clean