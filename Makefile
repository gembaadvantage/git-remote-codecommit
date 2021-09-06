# Copyright (c) 2021 Gemba Advantage
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# in the Software without restriction, including without limitation the rights
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

BINDIR     := $(CURDIR)/bin
BINNAME    ?= git-remote-codecommit
BINVERSION := ''
LDFLAGS    := -w -s

GOBIN = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

# Interrogate git for build time information
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_BRANCH = $(shell git branch --show-current)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)

BINVERSION = ${GIT_TAG}

ifneq ($(GIT_BRANCH),'main')
	BINVERSION := $(BINVERSION)-${GIT_SHA}
endif

# Set build time information
LDFLAGS += -X github.com/gembaadvantage/git-remote-codecommit/internal/version.version=${BINVERSION}
LDFLAGS += -X github.com/gembaadvantage/git-remote-codecommit/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/gembaadvantage/git-remote-codecommit/internal/version.gitBranch=${GIT_BRANCH}

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(BINNAME)' ./cmd/grc

.PHONY: test
test:
	go test -race -vet=off -p 1 -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: lint
lint:
	golangci-lint run --timeout 5m0s

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'