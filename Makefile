BINARY=golumen
VERSION=$(shell git describe --tags --dirty --always)
PKG=github.com/erikw/golumen
LDFLAGS=-ldflags "-X $(PKG)/internal/version.Version=$(VERSION)"

.PHONY: all build clean run install test release

all: build

build:
	go build $(LDFLAGS) -o $(BINARY)

clean:
	$(RM) $(BINARY)

# Usage: $ make run ARGS="-h"
run:
	go run $(LDFLAGS) ./main.go $(ARGS)

install:
	go install $(LDFLAGS)

test:
	go test ./...

# Usage: $ make release VERSION=v0.2.0
release:
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
