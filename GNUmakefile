TEST         ?= ./internal/...
WEBSITE_REPO = github.com/hashicorp/terraform-website
PKG_NAME     = dnsimple
HOSTNAME     = registry.terraform.io
NAMESPACE    = dnsimple
BINARY       = terraform-provider-${PKG_NAME}
VERSION      = $(shell git describe --tags --always | cut -c 2-)
OS_ARCH      := $(shell echo "$$(uname -s)_$$(go env GOARCH)" | tr A-Z a-z)

default: build

build: fmtcheck
	go install

install: build
# VERSION contains the Git commit, which is not a valid version for Terraform. We also use 0.0.1 so that the version never conflicts with versions from the registry, and also so it's easy to see when a local override is being used.
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/0.0.1/${OS_ARCH}
	mv ${GOPATH}/bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PKG_NAME}/0.0.1/${OS_ARCH}

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 10m $(ARGS)

fmt:
	gofmt -s -w .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

website:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

.PHONY: build test testacc fmt fmtcheck errcheck website
