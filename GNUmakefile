TEST         ?= ./internal/...
WEBSITE_REPO = github.com/hashicorp/terraform-website
PKG_NAME     = dnsimple
HOSTNAME     = registry.terraform.io
NAMESPACE    = dnsimple
BINARY       = terraform-provider-${PKG_NAME}
VERSION      = $(shell git describe --tags --always | cut -c 2-)
OS_ARCH      := $(shell echo "$$(uname -s)_$$(go env GOARCH)" | tr A-Z a-z)

.PHONY: default
default: build

.PHONY: build
build: fmtcheck
	go install

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/${VERSION}/${OS_ARCH}
	mv ${GOPATH}/bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m

.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 10m $(ARGS)

.PHONY: sweep
sweep:
	go run $(CURDIR)/tools/sweep/main.go

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: fmtcheck
fmtcheck:
	@test -z "$$(gofumpt -d -e . | tee /dev/stderr)"

.PHONY: errcheck
errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: website
website:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"
