GOCMD=go

MAIN=cmd/lwc/main.go
BIN=bin

VERSION=$(shell git describe --tags --abbrev --always)
# Same format as goreleaser
DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS="-s -w -X main.version=$(VERSION) -X main.date=$(DATE)"

all: build test

clean:
	rm -rf $(BIN)

build: build_debug build_release

test: test_unit test_integration

bin:
	mkdir -p $(BIN)

build_debug: bin
	$(GOCMD) build -o $(BIN)/lwc-debug -gcflags=all='-N -l' $(MAIN)

build_release: bin
	$(GOCMD) build -o $(BIN)/lwc -ldflags=$(LDFLAGS) $(MAIN)

test_unit:
	$(GOCMD) test -v ./...

test_integration: build_release
	test/integration.sh
