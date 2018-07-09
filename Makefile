<<<<<<< 1d205c236040b5b4063654bac8c16809d7a32324

=======
>>>>>>> Path changed in Makefile
PROJECT := github.com/kubernauts/tk8
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -ldflags="-w -X $(PROJECT)/cmd.GITCOMMIT=$(GITCOMMIT)"

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

