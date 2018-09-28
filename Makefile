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

release:
	./scripts/check-gofmt.sh
	go build -o golint github.com/golang/lint/golint
	./golint $(PKGS)
	go vet $(PKGS)
	go build ${BUILD_FLAGS} -o tk8 main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS} -o tk8-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS}  -o tk8-linux-amd64 main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build ${BUILD_FLAGS}  -o tk8-linux-386 main.go