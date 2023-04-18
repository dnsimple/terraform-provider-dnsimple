TEST         ?= ./internal/...
WEBSITE_REPO = github.com/hashicorp/terraform-website
PKG_NAME     = dnsimple
HOSTNAME     = registry.terraform.io
NAMESPACE    = dnsimple
BINARY       = terraform-provider-${PKG_NAME}
VERSION      = $(shell git describe --tags --always | cut -c 2-)
OS_ARCH      = darwin_$(shell uname -m)

default: build

build: fmtcheck
	go install

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/${VERSION}/${OS_ARCH}
	mv ${GOPATH}/bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/${VERSION}/${OS_ARCH}

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 10m $(ARGS)

sweep:
	go run $(CURDIR)/tools/sweep/main.go

fmt:
	gofmt -s -w .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

website:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

.PHONY: build test testacc fmt fmtcheck errcheck website
