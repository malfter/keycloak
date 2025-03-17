#!/usr/bin/env bash

GOLANGCI_LINT_VERION=v1.63.4

# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin "${GOLANGCI_LINT_VERION}"

golangci-lint --version
