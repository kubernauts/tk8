PROJECT := github.com/kubernauts/tk8
VERSION := $(shell git tag 2>/dev/null|tail -n 1)
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_FLAGS := -ldflags "-w -s -X $(PROJECT)/pkg/common.GITCOMMIT=$(GITCOMMIT) -X $(PROJECT)/pkg/common.VERSION=$(VERSION)"



default: bin

.PHONY: bin
bin:
	# go get -u ./...
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

.PHONY: test
test:
	gocov test ./... | gocov-xml > coverage.xml
	gometalinter.v1 --checkstyle > report.xml
	sonar-scanner \
  -Dsonar.projectKey=mmmac \
  -Dsonar.host.url=http://localhost:9000 \
  -Dsonar.login=616782f26ee441b650bd709eff9f8acee0a0fd75 \
	-X

.PHONY: release
release:
	go get -u ./...
	./scripts/check-gofmt.sh
	golint $(PKGS)
	go vet $(PKGS)
	go build ${BUILD_FLAGS} -o tk8 main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS} -o tk8-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS}  -o tk8-linux-amd64 main.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build ${BUILD_FLAGS}  -o tk8-linux-386 main.go