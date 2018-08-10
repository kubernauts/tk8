PROJECT := github.com/kubernauts/tk8
VERSION := $(shell git tag 2>/dev/null|tail -n 1)
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -ldflags="-w -X $(PROJECT)/cmd.GITCOMMIT=$(GITCOMMIT) -X $(PROJECT)/cmd.VERSION=$(VERSION)"

default: bin

.PHONY: bin
bin:
	go build ${BUILD_FLAGS} -o tk8 main.go

.PHONY: install
install:
	go install ${BUILD_FLAGS}

.PHONY: validate
validate: gofmt vet lint

.PHONY: gofmt
gofmt:
	./scripts/check-gofmt.sh

.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)
